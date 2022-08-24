package client

import (
	"context"
	"net/http"
	"time"
)

type RpcClient interface {
	ConnRelease(serverName string, conn Conn)
	NewConnect(serverName string) (Conn, error)
	Call(ctx context.Context, req Request, resp interface{}, callOption ...CallOption) error
	NewRequest(serverName string, serverMethod string, req interface{}, opts ...RequestOption) Request
}

type Conn interface {
	//rpc调度服务方法
	Call(ctx context.Context, req Request, resp interface{}, CallOption CallOptions) error
	Close() error
	//获取连接的创建时间
	Created() time.Time
	//获取连接的地址
	Remote() string

	Error() error

	//连接id
	Id() int64
}

// Request is the interface for a synchronous request used by Call or Stream
type Request interface {
	// 服务名
	Service() string
	// 请求方法
	Method() string
	// 请求主题，也就是参数
	Body() interface{}

	//SetHeader(key string, value interface{})
	//
	//GetHeader() map[string]interface{}

	Header() http.Header
}
