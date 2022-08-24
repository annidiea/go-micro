package client

import "net/http"

type rpcRequest struct {
	//服务名
	service string
	//方法
	method      string
	endpoint    string
	contentType string
	//参数
	body interface{}
	opts requestOptions

	//header map[string]interface{}
	header http.Header
}

func newRequest(service, endpoint string, request interface{}, reqOpts ...RequestOption) Request {
	var opts requestOptions

	for _, o := range reqOpts {
		o(&opts)
	}

	return &rpcRequest{
		service:  service,
		method:   endpoint,
		endpoint: endpoint,
		body:     request,
		opts:     opts,

		//header: make(map[string]interface{}),
		header: http.Header{},
	}
}

func (r *rpcRequest) ContentType() string {
	return r.contentType
}

func (r *rpcRequest) Service() string {
	return r.service
}

func (r *rpcRequest) Method() string {
	return r.method
}

func (r *rpcRequest) Endpoint() string {
	return r.endpoint
}

func (r *rpcRequest) Body() interface{} {
	return r.body
}

//func (r *rpcRequest) SetHeader(key string, value interface{}) {
//	r.header[key] = value
//}
//
//func (r *rpcRequest) GetHeader() map[string]interface{} {
//	return r.header
//}

func (r *rpcRequest) Header() http.Header {
	return r.header
}
