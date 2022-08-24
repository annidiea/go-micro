package breaker

import (
	"fmt"
)

type State int

const (
	StateClosed State = iota + 1
	StateHalfOpen
	StateOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	default:
		return fmt.Sprintf("unknow state:%d", s)
	}
}

// 次数，请求量，失败量
type Counts struct {
	Request uint32

	TotaolSuccess uint32

	TotalFailures uint32

	//连续的成功量
	ConsecutiveSuccess uint32

	//连续的失败量
	ConsecutiveFailure uint32
}

func (c *Counts) request() {
	c.Request++
}

func (c *Counts) success() {
	c.TotaolSuccess++
	c.ConsecutiveSuccess++

	c.ConsecutiveFailure = 0
}

func (c *Counts) failure() {
	c.TotalFailures++
	c.ConsecutiveFailure++
	c.ConsecutiveSuccess = 0
}

func (c *Counts) clear() {
	c.Request = 0
	c.TotaolSuccess = 0
	c.TotalFailures = 0
	c.ConsecutiveSuccess = 0
	c.ConsecutiveFailure = 0
}
