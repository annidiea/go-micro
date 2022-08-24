package server

import (
	"encoding/json"
	"go-micro/core/errors"

	"io"
	"sync"
)

var errMissingParams = errors.New("go-micro/rpc/server", "jsonrpc: request body missing params", 500)

type serverCodec struct {
	dec     *json.Decoder // for reading JSON values
	enc     *json.Encoder // for writing JSON values
	c       io.Closer
	req     serverRequest
	mutex   sync.Mutex // protects seq, pending
	seq     uint64
	pending map[uint64]*json.RawMessage
}

func NewServerCodec(conn io.ReadWriteCloser) ServerCodec {
	return &serverCodec{
		dec:     json.NewDecoder(conn),
		enc:     json.NewEncoder(conn),
		c:       conn,
		pending: make(map[uint64]*json.RawMessage),
	}
}

type serverRequest struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
	Id     *json.RawMessage `json:"id"`
}

func (r *serverRequest) reset() {
	r.Method = ""
	r.Params = nil
	r.Id = nil
}

type serverResponse struct {
	Id     *json.RawMessage `json:"id"`
	Result interface{}      `json:"result"`
	Error  interface{}      `json:"error"`
}

func (c *serverCodec) ReadRequestHeader(r *Request) error {
	c.req.reset()
	if err := c.dec.Decode(&c.req); err != nil {
		return err
	}
	r.ServiceMethod = c.req.Method
	c.mutex.Lock()
	c.seq++
	c.pending[c.seq] = c.req.Id
	c.req.Id = nil
	r.Seq = c.seq
	c.mutex.Unlock()

	return nil
}

func (c *serverCodec) ReadRequestBody(r *Request, x interface{}) error {
	if x == nil {
		return nil
	}
	if c.req.Params == nil {
		return errMissingParams
	}
	var params [1]Message

	if e := json.Unmarshal(*c.req.Params, &params); e != nil {
		return e
	}
	r.Header = params[0].Header

	return json.Unmarshal(params[0].Body, x)
}

var null = json.RawMessage([]byte("null"))

func (c *serverCodec) WriteResponse(r *Response, x interface{}) error {
	c.mutex.Lock()
	b, ok := c.pending[r.Seq]

	if !ok {
		c.mutex.Unlock()

		return errors.New("go-micro/rpc/server/serverCodec.WriteResponse", "invalid sequence number in response", 500)
	}
	delete(c.pending, r.Seq)
	c.mutex.Unlock()

	if b == nil {
		b = &null
	}
	resp := serverResponse{Id: b}
	if r.Error == "" {
		resp.Result = x
	} else {
		resp.Error = r.Error
	}
	return c.enc.Encode(resp)
}

func (c *serverCodec) Close() error {
	return c.c.Close()
}

//func ServeConn(conn io.ReadWriteCloser) {
//	ServeCodec(NewServerCodec(conn))
//}
