package breaker

import "time"

var (
	DefaultTimeout           = time.Duration(60) * time.Second
	DefaultMaxRequest uint32 = 1
)

func defaultReadyToTrip(counts Counts) bool {
	return counts.ConsecutiveFailure > 5
}

func defaultAcceptable(err error) bool {
	return err == nil
}

type option struct {
	name string

	//最大的请求次数
	maxRequests uint32

	//设置熔断超时时限
	timeout time.Duration

	//自定义熔断验证方法，判断是否开启熔断
	readyToTrip func(counts Counts) bool

	//在熔断状态发生改变时
	onStateChange func(name string, from State, to State)
}

type Option func(opt *option)

func NewOption() *option {
	return &option{
		name:          "",
		maxRequests:   DefaultMaxRequest,
		timeout:       DefaultTimeout,
		readyToTrip:   defaultReadyToTrip,
		onStateChange: nil,
	}
}

func WithName(name string) Option {
	return func(opt *option) {
		opt.name = name
	}
}

func WithMaxRequest(maxRequest uint32) Option {
	return func(opt *option) {
		opt.maxRequests = maxRequest
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(opt *option) {
		opt.timeout = timeout
	}
}

func WithReadyToTrip(readyToTrip func(counts Counts) bool) Option {
	return func(opt *option) {
		opt.readyToTrip = readyToTrip
	}
}

func WithOnStateChange(onStateChange func(name string, from State, to State)) Option {
	return func(opt *option) {
		opt.onStateChange = onStateChange
	}
}
