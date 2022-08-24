package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"go-micro/core/debug"
	"go-micro/core/errors"
	"io/ioutil"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

type rpcClient struct {
	opts *dialOptions
	mp   *managePool
	id   int64
}

func NewClient(opt ...DialOption) (client *rpcClient) {
	opts := newDialOptions()
	for _, o := range opt {
		o.apply(opts)
	}

	client = &rpcClient{
		opts: opts,
		mp:   newManagePool(),
	}

	poolOpts := PoolOptions{
		Size: opts.poolsize,
		TTL:  opts.poolTTL,
	}

	for serverName, server := range opts.Servers {
		debug.PrintDirExePos(dir+"NewClient", "创建 %v 连接池", serverName)
		poolOpts.CreateConnectHandle = client.newConnect(serverName, server)
		pool, err := initPool(poolOpts)
		if err == ErrCreateConnHandleNotExit {
			debug.PrintErrDirExePos(dir+"NewClient", err, "创建%v连接池出现异常", serverName)
			continue
		}

		client.mp.Add(serverName, pool)
	}

	return
}

// 回收连接
func (c *rpcClient) ConnRelease(serverName string, conn Conn) {
	pool, ok := c.mp.Get(serverName)
	if !ok {
		return
	}

	pool.Release(conn)
}

// 根据服务名创建连接
func (c *rpcClient) NewConnect(serverName string) (Conn, error) {
	pool, ok := c.mp.Get(serverName)
	if !ok {
		debug.PrintErrDirExePos(dir+":NewConnect", ErrNotServer, "获取%v连接池错误", serverName)
		return nil, ErrNotServer
	}

	ctx, _ := context.WithTimeout(context.TODO(), c.opts.connTimeout)
	conn, err := pool.Get(ctx)
	if err != nil {
		debug.PrintErrDirExePos(dir+":NewConnect", ErrNotServer, "从连接池中获取%v连接错误", serverName)
		return nil, errors.NotFound("go-micro/rpc/client/rpcClient.NewConnect", "server %s: not found", serverName)
	}

	return conn, nil

}

func (c *rpcClient) call(ctx context.Context, req Request, resp interface{}, callOption CallOptions) error {
	conn, err := c.NewConnect(req.Service())
	defer func() {
		c.ConnRelease(req.Service(), conn)
	}()

	if err != nil {
		debug.PrintErrDirExePos(dir, err, "获取服务连接 %v 异常", req.Service())
		return err
	}

	return conn.Call(ctx, req, resp, callOption)

}

func (c *rpcClient) Call(ctx context.Context, req Request, resp interface{}, callOption ...CallOption) error {
	//覆盖执行操作信息
	callOpts := c.opts.callOptions

	for _, opt := range callOption {
		opt(&callOpts)
	}

	//fmt.Println(callOpts)

	//是否设置超时
	d, ok := ctx.Deadline()
	if !ok {
		// no deadline so we create a new one
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, callOpts.requestTimeout)
		defer cancel()
	} else {
		opt := WithRequestTimeout(d.Sub(time.Now()))
		opt(&callOpts)
	}

	//复制call方法
	rcall := c.call

	//执行中间件方法
	for i := len(callOpts.CallWrappers); i > 0; i-- {
		rcall = callOpts.CallWrappers[i-1](rcall)
	}

	//执行失败重试
	retries := callOpts.retries
	ch := make(chan error, retries+1)
	var gerr error
	for i := 0; i < retries; i++ {
		go func(i int) {
			ch <- rcall(ctx, req, resp, callOpts)
		}(i)

		select {
		case <-ctx.Done():
			return errors.Timeout("go-micro/rpc/client/rpcClient.Call", "server %s.%s timeout", req.Service(), req.Method())
		case err := <-ch:
			// if the call succeeded lets bail early
			if err == nil {
				return nil
			}

			retry, rerr := callOpts.retry(ctx, req, i, err)
			if rerr != nil {
				return rerr
			}

			if !retry {
				return err
			}

			gerr = err

		}

	}

	return gerr

}

func (c *rpcClient) NewRequest(serverName string, serverMethod string, req interface{}, opts ...RequestOption) Request {
	return newRequest(serverName, serverMethod, req, opts...)
}

func (c *rpcClient) newConnect(serverName string, s *Server) CreateConnectHandle {
	return func() (Conn, error) {
		c.id++
		//建立连接
		client, err := c.getClient(s)
		if err != nil {
			debug.PrintErrDirExePos(dir+":newConnect", err, "创建%v服务出现异常", serverName)
			return &connect{
				id:  c.id,
				err: errors.InternalServerError("go-micro/rpc/client/rpcClient.NewConnect", "server %s: create rpc client err %s ", serverName, err),
			}, err
		}

		return &connect{
			client:  client,
			id:      c.id,
			addr:    s.Address,
			err:     nil,
			created: time.Now(),
		}, nil
	}
}

func (c *rpcClient) getClient(s *Server) (client *rpc.Client, err error) {

	debug.DD("openssl %v; s = %v", s.Openssl, s)
	if !s.Openssl {
		return jsonrpc.Dial(s.NetWork, s.Address)
	}

	certPool := x509.NewCertPool()
	certBytes, err := ioutil.ReadFile(s.CertFile)
	if err != nil {
		return nil, err
	}
	certPool.AppendCertsFromPEM(certBytes)

	config := &tls.Config{
		RootCAs:    certPool,
		ServerName: s.TlsServerName,
	}

	conn, err := tls.Dial(s.NetWork, s.Address, config)
	//conn, err := rpc.DialHTTP("tcp", ":8000")
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}

	//client, err := rpc.Dial("tcp", ":9503")
	return jsonrpc.NewClient(conn), err
}

//调度失败-》重试
