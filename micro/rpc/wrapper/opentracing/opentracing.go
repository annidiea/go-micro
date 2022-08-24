package opentracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"go-micro/rpc/client"
	"go-micro/rpc/server"
	"net/http"
)

// client => server : context统一，同过一个方法对context中的SpanContext进行解析即可
// 考虑当前的SpanContext是第一次还是第N
// 记录新的SpanContext到数据
type headerKey struct {
}

// 针对客户端
func NewCallWrapper(ot opentracing.Tracer) client.CallWrapper {
	if ot == nil {
		ot = opentracing.GlobalTracer()
	}

	return func(call client.CallFunc) client.CallFunc {
		return func(ctx context.Context, req client.Request, rsp interface{}, opts client.CallOptions) error {
			name := fmt.Sprintf("%s.%s", req.Service(), req.Method())
			ctx, span, err := StartSpanFromContext(ctx, ot, name)
			if err != nil {
				return err
			}
			defer span.Finish()
			// 发送给下一个节点
			//req.SetHeader("opentracing", HeaderFromContext(ctx))
			req.Header()["Uber-Trace-Id"] = HeaderFromContext(ctx)["Uber-Trace-Id"]

			if err = call(ctx, req, rsp, opts); err != nil {

				span.LogFields(opentracinglog.String("error", err.Error()))
				span.SetTag("error", true)
			}
			return err
		}
	}
}

func NewHandlerWrapper(ot opentracing.Tracer) server.HandlerWrapper {
	if ot == nil {
		ot = opentracing.GlobalTracer()
	}
	return func(call server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req *server.Request, argv, rsp interface{}) error {
			name := fmt.Sprintf("%s", req.ServiceMethod)
			ctx = context.WithValue(ctx, headerKey{}, opentracing.HTTPHeadersCarrier(req.Header))

			ctx, span, err := StartSpanFromContext(ctx, ot, name)
			if err != nil {
				return err
			}
			defer span.Finish()
			// 发送给下一个节点
			if err = call(ctx, req, argv, rsp); err != nil {
				span.LogFields(opentracinglog.String("test", fmt.Sprintf("%v", err)))
				span.LogFields(opentracinglog.String("error", err.Error()))
				span.SetTag("error", true)
			}
			return err
		}
	}
}

func StartSpanFromContext(ctx context.Context, tracer opentracing.Tracer, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	// 先判断header是不是存在信息
	carrier := HeaderFromContext(ctx)
	fmt.Println("carrier", carrier)
	// 解析；解析是不是存在父节点，第二个是不是存在下级
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	} else if spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, carrier); err == nil {
		opts = append(opts, opentracing.ChildOf(spanCtx))
	}

	// 创建新的span
	sp := tracer.StartSpan(name, opts...)

	// 还需去获取新的span
	if err := sp.Tracer().Inject(sp.Context(), opentracing.HTTPHeaders, carrier); err != nil {
		return nil, nil, err
	}

	ctx = opentracing.ContextWithSpan(ctx, sp)
	ctx = context.WithValue(ctx, headerKey{}, carrier)
	return ctx, sp, nil
}

func HeaderFromContext(ctx context.Context) opentracing.HTTPHeadersCarrier {
	v := ctx.Value(headerKey{})
	val, ok := v.(opentracing.HTTPHeadersCarrier)
	if !ok {
		h := http.Header{}
		return opentracing.HTTPHeadersCarrier(h)
	}

	return val
}
