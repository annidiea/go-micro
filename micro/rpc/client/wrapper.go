package client

import "context"

type CallFunc func(ctx context.Context, req Request, resp interface{}, callOption CallOptions) error

type CallWrapper func(callFunc CallFunc) CallFunc
