package main

import "github.com/gin-gonic/gin"

// Context 是一个自定义的请求上下文结构体，封装了 gin.Context
type Context struct {
	C *gin.Context
}

// New 创建一个新的自定义 Context 实例
func New(c *gin.Context) *Context {
	return &Context{C: c}
}

// HandlerFunc 定义 Handler 函数签名
type HandlerFunc func(*Context)

// Wrap 适配器函数
// 它接收我们自定义的 HandlerFunc，并返回一个标准的 gin.HandlerFunc
func Wrap(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在这里，我们执行了之前在每个 Handler 开头都要做的工作
		// 创建自定义上下文
		ctx := New(c)

		// 调用我们真正的业务逻辑 Handler
		handler(ctx)
	}
}

// wrapHandlers 是一个辅助函数，用于将我们自定义的 HandlerFunc 列表转换为 gin.HandlerFunc 列表
func wrapHandlers(handlers []HandlerFunc) []gin.HandlerFunc {
	wrappedHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		wrappedHandlers[i] = Wrap(handler)
	}
	return wrappedHandlers
}

// CustomRouterGroup 是我们自定义的路由组，它封装了 gin.RouterGroup
type CustomRouterGroup struct {
	*gin.RouterGroup
}

// Group 重写了原生的 Group 方法，确保返回的是我们自己的 CustomRouterGroup
// 这样可以支持无限层级的路由分组，且每一层都支持自定义 Handler
func (g *CustomRouterGroup) Group(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	// 调用 wrapHandlers 转换中间件
	ginHandlers := wrapHandlers(handlers)
	// 调用原生的 Group 方法创建新的 gin.RouterGroup
	newGinGroup := g.RouterGroup.Group(relativePath, ginHandlers...)
	// 将新的 gin.RouterGroup 包装成我们自己的 CustomRouterGroup 并返回
	return &CustomRouterGroup{RouterGroup: newGinGroup}
}

// 以下是所有 HTTP 方法的重写
// 它们都遵循相同的模式：
// 1. 接收自定义的 HandlerFunc
// 2. 调用 wrapHandlers 进行转换
// 3. 调用 gin.RouterGroup 中对应的原生方法
// 4. 返回 *CustomRouterGroup 以支持链式调用 (e.g., r.GET(...).Use(...))
//    注意：这里的 Use 仍然是 gin.Use，接收 gin.HandlerFunc

func (g *CustomRouterGroup) POST(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	g.RouterGroup.POST(relativePath, wrapHandlers(handlers)...)
	return g
}

func (g *CustomRouterGroup) GET(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	g.RouterGroup.GET(relativePath, wrapHandlers(handlers)...)
	return g
}

func (g *CustomRouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	g.RouterGroup.DELETE(relativePath, wrapHandlers(handlers)...)
	return g
}

func (g *CustomRouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	g.RouterGroup.PATCH(relativePath, wrapHandlers(handlers)...)
	return g
}

func (g *CustomRouterGroup) PUT(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	g.RouterGroup.PUT(relativePath, wrapHandlers(handlers)...)
	return g
}

func (g *CustomRouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	g.RouterGroup.OPTIONS(relativePath, wrapHandlers(handlers)...)
	return g
}

func (g *CustomRouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	g.RouterGroup.HEAD(relativePath, wrapHandlers(handlers)...)
	return g
}

func (g *CustomRouterGroup) Any(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	g.RouterGroup.Any(relativePath, wrapHandlers(handlers)...)
	return g
}

// Handle 方法也需要重写
func (g *CustomRouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	g.RouterGroup.Handle(httpMethod, relativePath, wrapHandlers(handlers)...)
	return g
}

// 注意: Static, StaticFile, StaticFS 等方法不需要重写，因为它们不接收 HandlerFunc
// 通过内嵌 *gin.RouterGroup，这些方法被自动继承，可以直接使用。

// Engine 是我们的自定义引擎，封装了 gin.Engine
type Engine struct {
	*gin.Engine
	RouterGroup *CustomRouterGroup
}

// NewEngine 创建并返回一个我们自定义的 Engine
func NewEngine() *Engine {
	e := gin.Default() // 或者 gin.New()，并按需添加中间件
	return &Engine{
		Engine: e,
		RouterGroup: &CustomRouterGroup{
			RouterGroup: &e.RouterGroup,
		},
	}
}

// 为了让 r.POST(...) 这种顶层调用生效，我们也需要在 Engine 上实现这些方法
// 这些方法直接代理到其内部的 RouterGroup

func (e *Engine) Group(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	return e.RouterGroup.Group(relativePath, handlers...)
}

func (e *Engine) POST(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	return e.RouterGroup.POST(relativePath, handlers...)
}

func (e *Engine) GET(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	return e.RouterGroup.GET(relativePath, handlers...)
}

func (e *Engine) DELETE(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	return e.RouterGroup.DELETE(relativePath, handlers...)
}

func (e *Engine) PATCH(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	return e.RouterGroup.PATCH(relativePath, handlers...)
}

func (e *Engine) PUT(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	return e.RouterGroup.PUT(relativePath, handlers...)
}

func (e *Engine) OPTIONS(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	return e.RouterGroup.OPTIONS(relativePath, handlers...)
}

func (e *Engine) HEAD(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	return e.RouterGroup.HEAD(relativePath, handlers...)
}

func (e *Engine) Any(relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	return e.RouterGroup.Any(relativePath, handlers...)
}

func (e *Engine) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) *CustomRouterGroup {
	return e.RouterGroup.Handle(httpMethod, relativePath, handlers...)
}
