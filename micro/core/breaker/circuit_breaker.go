package breaker

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrOpenState      = errors.New("熔断器开启状态")
	ErrTooManyRequest = errors.New("太多请求量")
)

type CircuitBreaker struct {
	name string

	opt *option

	////最大的请求次数
	//maxRequests uint32
	//
	////设置熔断超时时限
	//timeout time.Duration
	//
	////自定义熔断验证方法，判断是否开启熔断
	//readyToTrip func(counts Counts) bool
	//
	////在熔断状态发生改变时
	//onStateChange func(name string, from State, to State)

	mutex sync.Mutex

	//状态
	state State

	//同级熔断器次数
	counts Counts

	//记录熔断时限
	expiry time.Time
}

func NewBreaker(opts ...Option) *CircuitBreaker {
	cb := new(CircuitBreaker)
	opt := NewOption()
	for _, o := range opts {
		o(opt)
	}

	cb.name = opt.name
	cb.state = StateClosed
	cb.opt = opt

	return cb
}

func (cb *CircuitBreaker) Name() string {
	return cb.name
}

func (cb *CircuitBreaker) Counts() Counts {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	return cb.counts
}

func (cb *CircuitBreaker) State() State {
	cb.mutex.Lock()
	cb.mutex.Unlock()
	return cb.currentState()
}

/*
DoWithAcceptable(req func() error, acceptable Acceptable) error

		DoWithFallback(req func() error, fallback func(err error) error) error

		DoWithFallbackAcceptable(req func() error, fallback func(err error) error, acceptable Acceptable) error
*/

func (cb *CircuitBreaker) Do(req func() error) error {
	return cb.do(req, nil, defaultAcceptable)
}

func (cb *CircuitBreaker) DoWithAcceptable(req func() error, acceptable Acceptable) error {
	return cb.do(req, nil, acceptable)
}

func (cb *CircuitBreaker) DoWithFallback(req func() error, fallback func(err error) error) error {
	return cb.do(req, fallback, defaultAcceptable)
}

func (cb *CircuitBreaker) DoWithFallbackAcceptable(req func() error, fallback func(err error) error, acceptable Acceptable) error {
	return cb.do(req, fallback, acceptable)
}

// 熔断器执行请求
// req		：正常执行的方法
// fallback	：熔断后执行的方法
// acceptable：则是判断正常执行方法之后的异常是否可以通过，否则失败
func (cb *CircuitBreaker) do(req func() error, fallback func(err error) error, acceptable func(err error) bool) error {

	//读取当前的熔断状态，并判断是否处理方法
	if err := cb.accept(); err != nil {

		//当处于熔断状态就执行降级方法
		if fallback != nil {
			return fallback(err)
		}

		return err
	}

	defer func() {
		if e := recover(); e != nil {
			cb.failure()
			panic(e)
		}
	}()

	//执行实际方法
	err := req()
	//判断异常是否为调度失败
	if acceptable(err) {
		cb.success()
	} else {
		cb.failure()
	}

	return err

}

func (cb *CircuitBreaker) success() {
	switch cb.state {
	case StateClosed:
		//统计成功的次数
		cb.counts.success()

	case StateHalfOpen:
		//半开放的状态下失败，直接转化为开放状态
		cb.counts.success()
		if cb.counts.ConsecutiveSuccess >= cb.opt.maxRequests {
			cb.setState(StateClosed)
		}
	}
}

func (cb *CircuitBreaker) failure() {
	switch cb.state {
	case StateClosed:
		//统计失败的次数
		cb.counts.failure()

		if cb.opt.readyToTrip(cb.counts) {
			cb.setState(StateOpen)
		}
	case StateHalfOpen:
		//半开放的状态下失败，直接转化为开放状态
		cb.setState(StateOpen)
	}
}

// 判断当前的熔断器是否接受处理方法
func (cb *CircuitBreaker) accept() error {

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	state := cb.currentState()

	//根据熔断状态判断
	if state == StateOpen {
		return ErrOpenState
	} else if state == StateHalfOpen && cb.counts.Request >= cb.opt.maxRequests {
		return ErrTooManyRequest
	}

	cb.counts.request()
	return nil
}

// 获取当前是否熔断
func (cb *CircuitBreaker) currentState() State {
	now := time.Now()
	switch cb.state {
	case StateClosed:
	case StateOpen:
		//如果熔断超过时限，开始半开放状态
		if cb.expiry.Before(now) {
			cb.setState(StateHalfOpen)
		}
	}
	return cb.state
}

// 设置熔断
func (cb *CircuitBreaker) setState(state State) {
	if cb.state == state {
		return
	}

	prev := cb.state
	cb.state = state
	cb.toNewExpiry(time.Now())

	if cb.opt.onStateChange != nil {
		cb.opt.onStateChange(cb.name, prev, cb.state)
	}
}

func (cb *CircuitBreaker) toNewExpiry(now time.Time) {
	//重置，计数
	cb.counts.clear()

	var zero time.Time

	switch cb.state {
	case StateClosed:
	case StateOpen:
		cb.expiry = now.Add(cb.opt.timeout)
	default:
		cb.expiry = zero
	}
}
