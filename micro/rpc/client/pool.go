package client

import (
	"context"
	"fmt"
	"go-micro/core/debug"
	"time"
)

var dir = "core/rpc/client"

type Pool interface {
	Get(ctx context.Context) (Conn, error)

	Release(conn Conn)

	Close()
}

// 定一个创建连接的方法
type CreateConnectHandle func() (Conn, error)

type PoolOptions struct {
	Size int
	TTL  time.Duration
	CreateConnectHandle
}

// 管理连接池
type managePool struct {
	pools map[string]Pool
}

func newManagePool() *managePool {
	return &managePool{
		pools: make(map[string]Pool),
	}
}

func (mp *managePool) Add(tab string, pool Pool) {
	mp.pools[tab] = pool
}

func (mp *managePool) Get(tab string) (Pool, bool) {
	pool, ok := mp.pools[tab]
	return pool, ok
}

type pool struct {
	count int //用于创建的连接
	conns chan Conn
	size  int
	ttl   time.Duration
	CreateConnectHandle
}

func initPool(options PoolOptions) (*pool, error) {
	if options.Size <= 0 {
		return nil, ErrPoolSize
	}

	p := &pool{
		size:                options.Size,
		ttl:                 options.TTL,
		conns:               make(chan Conn, options.Size),
		CreateConnectHandle: options.CreateConnectHandle,
	}

	return p, p.init()
}

// 初识连接
func (p *pool) init() error {
	if p.CreateConnectHandle == nil {
		return ErrCreateConnHandleNotExit
	}

	debug.PrintDirExePos(dir+"init", "连接数 %d", p.size)
	// 创建连接
	for i := 0; i < p.size; i++ {
		//time.Sleep(time.Second)
		conn, err := p.CreateConnectHandle()
		if err != nil {
			continue
			//return err
		}

		p.count++
		p.conns <- conn
	}
	return nil
}

// 获取连接
func (p *pool) Get(ctx context.Context) (Conn, error) {
	//判断正在使用的连接是否少于总连接数
	//判断正在使用的连接是否少于总连接数
	if p.count < p.size {
		p.createConn()
	}
	for {
		select {
		case conn := <-p.conns:
			if d := time.Since(conn.Created()); d > p.ttl {
				conn.Close()
				p.count--
				// 创建新的连接
				p.createConn()
				continue
			}
			return conn, nil
		case <-ctx.Done():
			return nil, ErrPoolGetTimeout
		}
	}
}

// 释放连接
func (p *pool) Release(conn Conn) {
	// 可能连接为nil
	if conn == nil {
		return
	}

	if conn.Error() == nil {
		p.conns <- conn
		return
	}
	conn.Close()
	p.createConn()
}

func (p *pool) createConn() {
	go func() {
		retry := 0
		// 创建的时候会出现异常所以用for
		for {
			if p.CreateConnectHandle == nil {
				return
			}
			conn, err := p.CreateConnectHandle()
			if retry > 10 {
				return
			}
			if err != nil {
				retry++
				time.Sleep(1 * time.Second)
				continue
			}
			p.count++
			p.conns <- conn
			return
		}
	}()
}
func (p *pool) Close() {
	fmt.Println("关闭连接")
}
