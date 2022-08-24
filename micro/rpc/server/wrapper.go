package server

import "context"

type HandlerFunc func(ctx context.Context, req *Request, argv, rsp interface{}) error

type HandlerWrapper func(handlerFunc HandlerFunc) HandlerFunc
