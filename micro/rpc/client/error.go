package client

import "errors"

var (
	ErrNotServer = errors.New("查找不到服务配置")

	ErrPoolSize = errors.New("连接池大小不能小于0个")

	ErrAddrNotExist   = errors.New("不存在服务地址")
	ErrPoolGetTimeout = errors.New("连接获取超时")

	ErrCreateConnHandleNotExit = errors.New("创建连接的处理方法不存在")

	ErrRequestTimeout = errors.New("请求超时")
)
