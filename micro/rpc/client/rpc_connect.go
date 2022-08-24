package client

import (
	"context"
	"go-micro/core/errors"
	"net/rpc"
	"time"
)

type connect struct {
	client  *rpc.Client
	id      int64
	addr    string
	err     error
	created time.Time
}

// 调用
func (c *connect) Call(ctx context.Context, req Request, resp interface{}, callOption CallOptions) error {
	ch := make(chan error, 1)
	go func() {
		//ch <- c.client.Call(req.Method(), req.Body(), resp)

		ch <- c.client.Call(req.Method(), &Message{
			//Header: req.GetHeader(),
			Header: req.Header(),
			Body:   req.Body(),
		}, resp)
	}()

	select {
	case err := <-ch:
		if err != nil {
			return errors.Parse(err.Error())
			return errors.InternalServerError("go-micro/rpc/client/rpcConnect.Call", "server %s.%s.%v", req.Service(), req.Method(), err.Error())
		}
		return nil
	case <-ctx.Done():
		return errors.Timeout("go-micro/rpc/client/rpcConnect.Call", "server %s.%s", req.Service(), req.Method())

	}
}

func (c *connect) Close() error {
	return c.client.Close()
}

func (c *connect) Created() time.Time {
	return c.created
}

func (c *connect) Remote() string {
	return c.addr
}

func (c *connect) Error() error {
	return c.err
}

func (c *connect) Id() int64 {
	return c.id
}
