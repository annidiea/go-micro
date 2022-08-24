package client

import "time"

var (
	DefaultPoolSize    = 2
	DefaultPoolTTL     = 10 * time.Minute
	DefaultConnTimeout = 30 * time.Second

	//默认重试次数
	DefaultRetries = 1
	//默认重试超时时间
	DefaultRequestTimeout = 3 * time.Second
	//默认重试验证方法
	DefaultRetry = RetryAlways
)

type Server struct {
	Openssl       bool
	CertFile      string
	TlsServerName string
	NetWork       string
	Address       string
}

type dialOptions struct {
	//需要连接的服务
	Servers map[string]*Server
	//连接池大小
	poolsize int

	//连接生命周期
	poolTTL time.Duration

	//连接超时
	connTimeout time.Duration

	//调用属性
	callOptions CallOptions
}

type CallOptions struct {
	CallWrappers []CallWrapper

	// Address of remote hosts
	address []string
	// 根据异常校验是否重试
	retry RetryFunc
	// 重试次数
	retries int
	// 请求超时
	requestTimeout time.Duration
}

//var defaultDialptions = dialOptions{
//	Servers: make(map[string]*Server),
//}

func newDialOptions() *dialOptions {
	return &dialOptions{
		Servers:     make(map[string]*Server),
		poolsize:    DefaultPoolSize,
		poolTTL:     DefaultPoolTTL,
		connTimeout: DefaultConnTimeout,
		callOptions: CallOptions{
			retry:          DefaultRetry,
			retries:        DefaultRetries,
			requestTimeout: DefaultRequestTimeout,
		},
	}
}

type DialOption interface {
	apply(*dialOptions)
}

type funcDialOption struct {
	f func(*dialOptions)
}

func newFuncDialOption(f func(*dialOptions)) *funcDialOption {
	return &funcDialOption{
		f: f,
	}
}

func (fdo *funcDialOption) apply(do *dialOptions) {
	fdo.f(do)
}

func SetServer(name string, server *Server) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		if server.CertFile != "" && server.TlsServerName != "" {
			server.Openssl = true
		}

		o.Servers[name] = server
	})
}

func SetPoolSize(size int) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.poolsize = size
	})
}

func SetPoolTTL(ttl time.Duration) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.poolTTL = ttl
	})
}

func SetConnectTimeOut(timeout time.Duration) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.connTimeout = timeout
	})
}

func RequestTimeout(timeout time.Duration) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.callOptions.requestTimeout = timeout
	})
}

func Retries(retries int) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.callOptions.retries = retries
	})
}

func Retry(fn RetryFunc) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.callOptions.retry = fn
	})
}

type CallOption func(options *CallOptions)

// 全局设置，请求超时
func WithRequestTimeout(timeout time.Duration) CallOption {
	return func(options *CallOptions) {
		options.requestTimeout = timeout
	}
}

func WithRetries(retries int) CallOption {

	return func(options *CallOptions) {
		options.retries = retries
	}
}

func WithRetry(fn RetryFunc) CallOption {
	return func(options *CallOptions) {
		options.retry = fn
	}
}

func WrapCall(cw ...CallWrapper) CallOption {
	return func(options *CallOptions) {
		options.CallWrappers = append(options.CallWrappers, cw...)
	}
}

func WithWrapCall(cw ...CallWrapper) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.callOptions.CallWrappers = append(options.callOptions.CallWrappers, cw...)
	})
}

type requestOptions struct {
}

type RequestOption func(options *requestOptions)
