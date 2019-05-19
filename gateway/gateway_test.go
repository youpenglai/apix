package gateway

import "testing"

const testApiDoc = `version: 1.0.0
baseUrl: /api/
description: 这是一个测试API文档
types:
  - name: TestType
    description: 测试数据类型
    members:
      name:
        type: string
        required: true
        minLength: 6
      age:
        type: integer
        required: true
      tags:
        type: [string]
apis:
  - url: /grpc/test
    description: GRPC测试
    method: post
    params:
      body:
        name:
          type: string
          required: true
          minLength: 6
        age:
          type: integer
          required: true
        tags:
          type: [string]
    forwards:
      - name: grpc
        service: my-service
        grpc:
          method: hello
          paramMapper:
            name: name
            age: age
            tags: tags
        test:
          name: wang
          $gt:
            age: 10
        onfail: reject
    returns:
      - "200":
        description: 测试结果
        type: json
        data:
          type: TestType
  - url: /redis/test
    description: Redis测试
    method: post
    params:
      header:
        token:
          type: string
    forwards:
      - name: redis
        service: redis
        redis:
          key: token|pl.%s
          # type: string, hash etc...
          type: string
    returns:
      - "200":
        description: Redis测试
        data:
          token: string

`

func TestNewApiGateWay(t *testing.T) {
	gateway := NewApiGateWay()
	err := gateway.AddApiDoc("user.yaml", []byte(testApiDoc))
	if err != nil {
		t.Error(err)
		return
	}

	gateway.Serve()
}
