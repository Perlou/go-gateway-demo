package middleware

import (
	"context"
	"net"

	"github.com/perlou/go-gateway-demo/proxy/tcp_proxy"
)

//目标定位是 tcp、http通用的中间件
//知其然也知其所以然

type TcpHandlerFunc func(*TcpSliceRouterContext)

// router 结构体
type TcpSliceRouter struct {
	groups []*TcpSliceGroup
}

// group 结构体
type TcpSliceGroup struct {
	*TcpSliceRouter
	path     string
	handlers []TcpHandlerFunc
}

// router上下文
type TcpSliceRouterContext struct {
	conn net.Conn
	Ctx  context.Context
	*TcpSliceGroup
	index int8
}

func newTcpSliceRouterContext(conn net.Conn, r *TcpSliceRouter, ctx context.Context) *TcpSliceRouterContext {
	newTcpSliceGroup := &TcpSliceGroup{}
	*newTcpSliceGroup = *r.groups[0] //浅拷贝数组指针
	c := &TcpSliceRouterContext{conn: conn, TcpSliceGroup: newTcpSliceGroup, Ctx: ctx}
	c.Reset()
	return c
}

func (c *TcpSliceRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

func (c *TcpSliceRouterContext) Set(key, val interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, val)
}

type TcpSliceRouterHandler struct {
	coreFunc func(*TcpSliceRouterContext) tcp_proxy.TCPHandler
	router   *TcpSliceRouter
}

func (w *TcpSliceRouterHandler) ServeTCP(ctx context.Context, conn net.Conn) {
	c := newTcpSliceRouterContext(conn, w.router, ctx)
	c.handlers = append(c.handlers, func(c *TcpSliceRouterContext) {
		w.coreFunc(c).ServeTCP(ctx, conn)
	})
	c.Reset()
	c.Next()
}

func NewTcpSliceRouterHandler(coreFunc func(*TcpSliceRouterContext) tcp_proxy.TCPHandler, router *TcpSliceRouter) *TcpSliceRouterHandler {
	return &TcpSliceRouterHandler{
		coreFunc: coreFunc,
		router:   router,
	}
}

// 构造 router
func NewTcpSliceRouter() *TcpSliceRouter {
	return &TcpSliceRouter{}
}

// 创建 Group
func (g *TcpSliceRouter) Group(path string) *TcpSliceGroup {
	return &TcpSliceGroup{
		TcpSliceRouter: g,
		path:           path,
	}
}

// 构造回调方法
func (g *TcpSliceGroup) Use(middlewares ...TcpHandlerFunc) *TcpSliceGroup {
	g.handlers = append(g.handlers, middlewares...)
	existsFlag := false
	for _, oldGroup := range g.TcpSliceRouter.groups {
		if oldGroup == g {
			existsFlag = true
		}
	}
	if !existsFlag {
		g.TcpSliceRouter.groups = append(g.TcpSliceRouter.groups, g)
	}
	return g
}

// 从最先加入中间件开始回调
func (c *TcpSliceRouterContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		//fmt.Println("c.index")
		//fmt.Println(c.index)
		c.handlers[c.index](c)
		c.index++
	}
}

// 跳出中间件方法
func (c *TcpSliceRouterContext) Abort() {
	c.index = abortIndex
}

// 是否跳过了回调
func (c *TcpSliceRouterContext) IsAborted() bool {
	return c.index >= abortIndex
}

// 重置回调
func (c *TcpSliceRouterContext) Reset() {
	c.index = -1
}
