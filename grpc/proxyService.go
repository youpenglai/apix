package grpc

import (
	"io"
	"sync/atomic"
	"encoding/binary"
	"encoding/json"
	"sync"
	"os"
	"bytes"
)

var globalId uint64

const (
	ipcMsgTypeCall = iota
	ipcMsgTypeReply
)

type IPCMessage struct {
	id uint64
	msgType int8
	body []byte
}

func (msg *IPCMessage) SetId(id uint64) {
	if id == 0 {
		msg.id = getGId()
	} else {
		msg.id = id
	}
}

func (msg *IPCMessage) GetId() uint64 {
	return msg.id
}

func (msg *IPCMessage) SetData(data []byte) {
	msg.body = data
}

func (msg *IPCMessage) GetData() ([]byte) {
	return msg.body
}

func getGId() uint64 {
	return atomic.AddUint64(&globalId, 1)
}

type ProxyServiceRegMsg struct {
	ServiceNames []string `json:"serviceNames"`
}

type ProxyServiceCall struct {
	ServiceName string `json:"serviceName"`
	Method string `json:"method"`
	Params []byte `json:"params"`
}

func (psc *ProxyServiceCall) Marshal() (data []byte, err error) {
	buff := bytes.NewBuffer(nil)
	if _, err = buff.WriteString(psc.ServiceName); err != nil {
		return
	}
	if err = buff.WriteByte(0); err != nil {
		return
	}

	if _, err = buff.WriteString(psc.Method); err != nil {
		return
	}
	if err = buff.WriteByte(0); err != nil {
		return
	}

	if _, err = buff.Write(psc.Params); err != nil {
		return
	}

	data = buff.Bytes()

	return
}

func (psc *ProxyServiceCall) UnMarshal(rawData []byte) (err error) {
	end := bytes.IndexByte(rawData, 0)
	psc.ServiceName = string(rawData[:end])
	s := end + 1
	end = bytes.IndexByte(rawData[s:], 0)

	psc.Method = string(rawData[s: s + end])
	s = s + end + 1
	psc.Params = rawData[s:]
	return
}

type ProxyService struct {
	reader io.Reader
	writer io.Writer
	ioReady chan int
	messageBuff chan *IPCMessage

	callWaiter map[uint64]chan[]byte
	callWaiterMu sync.Mutex

	callHandler ProxyCallHandler
}

func (sp *ProxyService) Attach(reader io.Reader, writer io.Writer) {
	sp.reader = reader
	sp.writer = writer
	sp.ioReady <- 1
}

func (sp *ProxyService) writeMessage(msg *IPCMessage) (err error) {
	sp.messageBuff <- msg
	return nil
}

func (sp *ProxyService) readMessage(msg *IPCMessage) (err error) {
	if err = binary.Read(sp.reader, binary.LittleEndian, &msg.id); err != nil {
		return
	}
	if err = binary.Read(sp.reader, binary.LittleEndian, &msg.msgType); err != nil {
		return
	}
	var size int32
	if err = binary.Read(sp.reader, binary.LittleEndian, &size); err != nil {
		return
	}
	if size == 0 {
		return
	}

	msg.body = make([]byte, size)
	err = binary.Read(sp.reader, binary.LittleEndian, msg.body)
	return
}

func (sp *ProxyService) CallAsync(param interface{}) (retCh chan[]byte, err error) {
	var paramData []byte
	switch param.(type) {
	case []byte:
		paramData = param.([]byte)
	default:
		if paramData, err = json.Marshal(param); err != nil {
			return
		}
	}

	var msg IPCMessage
	msg.SetId(0)
	msg.msgType = ipcMsgTypeCall
	msg.SetData(paramData)
	sp.writeMessage(&msg)
	retCh = make(chan[]byte, 1)
	sp.callWaiterMu.Lock()
	sp.callWaiter[msg.GetId()] = retCh
	sp.callWaiterMu.Unlock()

	return
}

func (sp *ProxyService) CallSync(param interface{}) (retData []byte, err error) {
	var retCh chan[]byte
	if retCh, err = sp.CallAsync(param); err != nil {
		return
	}

	retData = <-retCh
	return
}

type ProxyCallHandler func(param []byte) (retData []byte, err error)

func (sp *ProxyService) OnCall(handler ProxyCallHandler) {
	sp.callHandler = handler
}

// 暂时使用
func (sp *ProxyService) handleCall(callId uint64, data []byte) error {
	if sp.callHandler == nil {
		return nil
	}

	retData, err := sp.callHandler(data)
	if err != nil {
		return err
	}

	var msg IPCMessage
	msg.id = callId
	msg.msgType = ipcMsgTypeReply
	msg.body = retData
	return sp.writeMessage(&msg)
}

func IPCCallHandler(proxy *ProxyService) {
	<- proxy.ioReady
	for {
		var msg IPCMessage
		if err := proxy.readMessage(&msg); err != nil {
			// TODO:  错误处理，EOF?
			return
		}

		if msg.msgType == ipcMsgTypeCall {
			err := proxy.handleCall(msg.id, msg.body)
			if err != nil {
				// TODO: 错误处理
			}
		} else {
			proxy.callWaiterMu.Lock()
			waiter, exist := proxy.callWaiter[msg.id]
			if exist {
				waiter <- msg.GetData()
				delete(proxy.callWaiter, msg.id)
			}
			proxy.callWaiterMu.Unlock()
		}
	}
}

func writeMessageHandler(proxy *ProxyService) {
	for {
		var err error
		var size int32

		msg := <- proxy.messageBuff
		size = int32(len(msg.body))
		if err = binary.Write(proxy.writer, binary.LittleEndian, msg.id); err != nil {
			return
		}
		if err = binary.Write(proxy.writer, binary.LittleEndian, msg.msgType); err != nil {
			return
		}

		if err = binary.Write(proxy.writer, binary.LittleEndian, size); err != nil {
			return
		}
		if err = binary.Write(proxy.writer, binary.LittleEndian, msg.body); err != nil {
			return
		}
	}
}

func NewServiceProxy() *ProxyService {
	proxy := &ProxyService{
		callWaiter: make(map[uint64]chan[]byte),
		ioReady: make(chan int, 1),
		messageBuff: make(chan *IPCMessage, 1),
	}

	go IPCCallHandler(proxy)
	go writeMessageHandler(proxy)
	return proxy
}

type ServiceCallHandler func(call *ProxyServiceCall) (data []byte, err error)

func HandleServiceCall(proxyService *ProxyService, handler ServiceCallHandler) {
	proxyService.OnCall(func(param []byte) (retData []byte, err error) {
		var call ProxyServiceCall

		if err = call.UnMarshal(param); err != nil {
			return
		}
		return handler(&call)
	})
}

func RegisterService(proxyService *ProxyService, serviceNames ...string) error {
	_, err := proxyService.CallSync(ProxyServiceRegMsg{ServiceNames: serviceNames})
	return err
}

func InitServiceProxy() *ProxyService {
	proxy := NewServiceProxy()
	proxy.Attach(os.Stdin, os.Stdout)
	return proxy
}
