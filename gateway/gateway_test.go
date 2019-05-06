package gateway

import "testing"

const testApiDoc = `version: 1.1.0
baseUrl: /v1/
description: 朋来互动用户API接口
types:
  - name: UserInfo
    description: 用户基本信息
    members:
      name:
        description: 用户姓名
        type: string
        required: false
apis:
  - url: /auth/login
    method: post
    description: 用户登录，使用朋来账户登录
    params:
      body:
        userName:
          description: 用户名（朋来用户登录）
          type: string
          required: true
          minLength: 6
          maxLength: 20
        password:
          description: ''
          type: string
          required: true
          minLength: 6
          maxLength: 16
    returns:
      '200':
        status: '200'
        description: 登录成功
        type: json
        data:
          token:
            description: 授权TOKEN
            type: string
          nickName:
            description: 用户昵称
            type: string
          id:
            description: 用户ID
            type: number
          avatar:
            description: 用户头像地址
            type: string
      '404':
        status: '404'
        description: 没有找到用户登录，或者密码错误
        type: json
        data:
          errCode:
            description: ''
            type: number
          errMsg:
            description: ''
            type: string
  - url: /auth/logout
    method: post
    description: 用户登出
    params:
      header:
        pl-token:
          description: 朋来TOKEN
          type: string
          required: true
    returns:
      '200':
        status: '200'
        description: 登出成功
        type: json
        data: {}
  - url: /auth/register
    method: post
    description: 注册账户，注册为朋来账户
    params:
      body:
        phone:
          description: 用户名称，使用手机号(中国大陆)作为用户名
          type: string
          required: true
          length: 11
        password:
          description: 用户密码
          type: string
          required: true
          minLength: 6
          maxLength: 16
        verifyCode:
          description: ''
          type: string
          required: true
          length: 4
    returns:
      '200':
        status: '200'
        description: 注册并登录成功
        type: json
        data:
          token:
            description: 授权TOKEN
            type: string
          nickName:
            description: 用户昵称
            type: string
          id:
            description: 用户ID
            type: number
          avatar:
            description: 用户头像地址
            type: string
  - url: /auth/verifycode
    method: post
    description: 发送验证码
    params:
      body:
        content:
          description: |
            发送的内容：
            {"loginname":"13258335315","type":1}
            json字符串整体加密，loginname:手机号码,type
            1:注册
            2:忘记密码
            3:更换手机号
            4:验证支付密码
            5:绑定手机号
            6:验资实名信息
            7:众筹(验证个人信息)
            8:重置密码。
          type: string
          required: true
    returns:
      '200':
        status: '200'
        description: 短信发送成功
        type: json
        data: {}
  - url: /auth/resetpassword
    method: post
    description: 重置用户密码
    params:
      body:
        phone:
          description: ''
          type: string
          required: true
          length: 11
        newpassword:
          description: ''
          type: string
          required: true
          minLength: 6
          maxLength: 16
        verifyCode:
          description: ''
          type: string
          required: true
          length: 4
    returns:
      '200':
        status: '200'
        description: ''
        type: json
        data: {}
  - url: /user/info
    method: get
    description: 获取当前用户信息
    params:
      header:
        pl-token:
          description: ''
          type: string
          required: true
    returns:
      '200':
        status: '200'
        description: ''
        type: json
        data:
          nickName:
            description: ''
            type: string
          phone:
            description: 用户手机号，朋来登录账号，手机号需要打码：比如：132****3333
            type: string
          avatar:
            description: 用户头像
            type: string
      '401':
        status: '401'
        description: 未登录或者登录超时，需要重新登录
        type: json
        data: {}
  - url: /user/info
    method: post
    description: 修改用户基本信息
    params:
      body:
        nickName:
          description: ''
          type: string
          required: false
        avatar:
          description: ''
          type: string
          required: false
      header:
        pl-token:
          description: ''
          type: string
          required: true
    returns:
      '200':
        status: '200'
        description: ''
        type: json
        data:
          nickName:
            description: ''
            type: string
          phone:
            description: 用户手机号，朋来登录账号，手机号需要打码：比如：132****3333
            type: string
          avatar:
            description: 用户头像
            type: string
      '401':
        status: '401'
        description: 未登录或登录超时，需要重新登录
        type: json
        data: {}
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
