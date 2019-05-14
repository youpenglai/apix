package main

import (
	"os"
	"github.com/youpenglai/mfwgo/registry"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
	"github.com/youpenglai/apix/proxy"
	"os/signal"
	"syscall"
)

const (
	ENV_REDIS_ADDR_NAME = "REDIS_ADDR"
	ENV_REDIS_PWD_NAME = "REDIS_PWD"
	ENV_REDIS_SERVICE_NAME = "REDIS_SERVICE_NAME"

	DEFAULT_REDIS_SERVICE_NAME = "redis"
)

func getRedisServiceName() string {
	serviceName := os.Getenv(ENV_REDIS_SERVICE_NAME)
	if serviceName == "" {
		serviceName = DEFAULT_REDIS_SERVICE_NAME
	}
	return serviceName
}

func getRedisConfFromConsul() (addr, pwd string, db int, err error) {
	serviceName := getRedisServiceName()
	var serviceInfo *registry.ServiceInfo
	serviceInfo, err = registry.DiscoverService(serviceName)
	if err != nil {
		return
	}

	// TODO: 从consul中读取Redis的配置信息

	addr = fmt.Sprintf("%s:%d", serviceInfo.Address, serviceInfo.Port)
	return
}

func getRedisConfFromEnv() (addr, pwd string, db int, err error) {
	addr = os.Getenv(ENV_REDIS_ADDR_NAME)
	if addr == "" {
		err = errors.New("get redis from env failed")
	}

	pwd = os.Getenv(ENV_REDIS_PWD_NAME)
	return
}

func getRedisConf() (addr, pwd string, db int, err error) {
	addr, pwd, db, err = getRedisConfFromConsul()
	if err == nil {
		return
	}
	addr, pwd, db, err = getRedisConfFromEnv()
	return
}

func wait()(c chan os.Signal) {
	c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	return
}

var redisPool *redis.Pool

func main() {
	addr, pwd, db, err := getRedisConf()
	if err != nil {
		panic("Read redis config failure!")
	}

	redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			var conn redis.Conn
			conn, err = redis.Dial("tcp", addr, redis.DialPassword(pwd))
			if err != nil {
				return nil, err
			}

			if _, err = conn.Do("SELECT", db); err != nil {
				conn.Close()
				return nil, err
			}

			return conn, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:15,
	}

	proxyInst := proxy.InitServiceProxy()
	err = proxy.RegisterService(proxyInst, getRedisServiceName())
	if err != nil {
		panic("Register service error:" + err.Error())
	}

	proxy.HandleServiceCall(proxyInst, func(call *proxy.ProxyServiceCall) (data []byte, err error) {
		
		return
	})

	<-wait()
}
