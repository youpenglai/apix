package gateway

import (
	"io"
	"sync/atomic"
	"encoding/binary"
	"encoding/json"
	"sync"
	"os"
)

var globalId uint64

type IPCMessage struct {
	id uint64
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

type ProxyService struct {
	reader io.Reader
	writer io.Writer

	callWaiter map[uint64]chan[]byte
	callWaiterMu sync.Mutex

	callHandler ProxyCallHandler
}

func (sp *ProxyService) Attach(reader io.Reader, writer io.Writer) {
	sp.reader = reader
	sp.writer = writer
}

func (sp *ProxyService) writeMessage(msg *IPCMessage) (err error) {
	var size int32
	size = int32(len(msg.body))
	if err = binary.Write(sp.writer, binary.LittleEndian, msg.id); err != nil {
		return
	}
	if err = binary.Write(sp.writer, binary.LittleEndian, size); err != nil {
		return
	}
	err = binary.Write(sp.writer, binary.LittleEndian, msg.body)
	return
}

func (sp *ProxyService) readMessage(msg *IPCMessage) (err error) {
	if err = binary.Read(sp.reader, binary.LittleEndian, &msg.id); err != nil {
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
	if paramData, err = json.Marshal(param); err != nil {
		return
	}
	var msg IPCMessage
	msg.SetId(0)
	msg.SetData(paramData)
	sp.writeMessage(&msg)
	retCh = make(chan[]byte)
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
	msg.body = retData
	return sp.writeMessage(&msg)
}

func IPCCallHandler(proxy *ProxyService) {
	for {
		var msg IPCMessage
		if err := proxy.readMessage(&msg); err != nil {
			// TODO:  错误处理，EOF?
			return
		}
		proxy.callWaiterMu.Lock()
		waiter, exist := proxy.callWaiter[msg.id]
		if exist {
			waiter <- msg.GetData()
			delete(proxy.callWaiter, msg.id)
		} else {
			// dispatch
			err := proxy.handleCall(msg.id, msg.body)
			if err != nil {
				// TODO: 错误处理
			}
		}
		proxy.callWaiterMu.Unlock()
	}
}

func NewServiceProxy() *ProxyService {
	proxy := &ProxyService{
		callWaiter: make(map[uint64]chan[]byte),
	}

	go IPCCallHandler(proxy)
	return proxy
}

type ServiceCallHandler func(call *ProxyServiceCall) (data []byte, err error)

func HandleServiceCall(proxyService *ProxyService, handler ServiceCallHandler) {
	proxyService.OnCall(func(param []byte) (retData []byte, err error) {
		var call ProxyServiceCall
		if err = json.Unmarshal(param, &call); err != nil {
			return
		}
		return handler(&call)
	})
}

func InitServiceProxy() *ProxyService {
	proxy := NewServiceProxy()
	proxy.Attach(os.Stdin, os.Stdout)
	return proxy
}
