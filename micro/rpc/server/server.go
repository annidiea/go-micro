package server

import (
	"crypto/tls"
	"fmt"
	"go-micro/core/debug"
	"net"
	"os"
)

var dir = "core/rpc/server/"

type RpcServer struct {
	opts  serverOptions
	count int
	svr   *Server
}

func NewRpcServer(opt ...ServerOption) *RpcServer {
	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return &RpcServer{
		opts: opts,
		svr:  NewServer(opts),
	}
}

// 注册服务
func (s *RpcServer) Register(server interface{}) error {
	return s.svr.Register(server)
}

func (s *RpcServer) RegisterName(name string, rcvr interface{}) error {
	return s.svr.RegisterName(name, rcvr)
}

// 启动服务
func (s *RpcServer) Run(addr ...string) (err error) {
	defer func() {
		debug.DE(err)
	}()
	address := resolveAddress(addr)

	debug.DD("listening and serving TCP on %s \n", address)

	lis, err := s.listen(address)

	if err != nil {
		return
	}

	for {
		conn, err := lis.Accept()

		if err != nil {
			continue
		}

		s.count++
		debug.PrintDirExePos(dir+"server.go", "连接数 %d", s.count)
		go func(conn net.Conn) {
			debug.PrintDirExePos(dir+"server.go", "连接数 %d, %s", s.count, "进入请求")
			s.svr.ServeCodec(NewServerCodec(conn))
			debug.PrintDirExePos(dir+"server.go", "连接数 %d, %s", s.count, "完成请求")
			s.count--
		}(conn)
	}

	return nil
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			debug.DD("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debug.DD("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

func (s *RpcServer) listen(address string) (lis net.Listener, err error) {
	//未开启认证
	if !s.opts.openssl {
		debug.DD("未开启tls认证")
		return net.Listen("tcp", address)
	}

	//开启认证
	debug.DD("开启tls认证")
	fmt.Println(s.opts.certFile)
	cert, err := tls.LoadX509KeyPair(s.opts.certFile, s.opts.keyFile)

	if err != nil {
		return
	}
	return tls.Listen("tcp", address, &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
}
