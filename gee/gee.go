package gee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

type Engine struct {
	*RouterGroup  // 匿名实现继承
	router        *router
	groups        []*RouterGroup
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

/**
路由组
*/
type RouterGroup struct {
	prefix     string        // 路由前缀
	middleware []HandlerFunc // 中间件处理函数
	engine     *Engine       // 全局共享一个engine
	parent     *RouterGroup  // 父级路由组
}

/**
初始化核心
*/
func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Recovery())
	return engine
}

/**
添加路由分组
*/
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		engine: engine,
		parent: g,
	}
	engine.groups = append(engine.groups, newGroup)

	return newGroup
}

func (g *RouterGroup) AddRoute(method, path string, handleFunc HandlerFunc) {
	pattern := g.prefix + path
	g.engine.router.addRoute(method, pattern, handleFunc)
}

func (g *RouterGroup) GET(path string, handleFunc HandlerFunc) {
	g.AddRoute("GET", path, handleFunc)
}

func (g *RouterGroup) POST(path string, handleFunc HandlerFunc) {
	g.AddRoute("POST", path, handleFunc)
}

func (g *RouterGroup) PUT(path string, handleFunc HandlerFunc) {
	g.AddRoute("PUT", path, handleFunc)
}

func (g *RouterGroup) DELETE(path string, handleFunc HandlerFunc) {
	g.AddRoute("DELETE", path, handleFunc)
}

func (g *RouterGroup) PATCH(path string, handleFunc HandlerFunc) {
	g.AddRoute("PATCH", path, handleFunc)
}

func (g *RouterGroup) OPTION(path string, handleFunc HandlerFunc) {
	g.AddRoute("OPTION", path, handleFunc)
}

/**
添加中间件函数
*/
func (g *RouterGroup) Use(middleware ...HandlerFunc) {
	g.middleware = append(g.middleware, middleware...)
}

/**
创建文件处理的函数
*/
func (g *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(g.prefix, relativePath) // 获取完整的路径     /v1/assets/php.js
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(ctx *Context) {
		// 判断文件路径是否存在
		file := ctx.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			ctx.SetHttpStatus(http.StatusNotFound) // 未找到
			return
		}

		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

/**
静态文件路径
*/
func (g *RouterGroup) Static(relativePath string, root string) {
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	g.GET(urlPattern, handler)
}

/**
实现http自定义路由的接口
*/
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middleware []HandlerFunc

	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middleware = append(middleware, group.middleware...)
		}
	}

	context := NewContext(w, r)
	context.handlers = middleware // 将中间件加到context上
	context.engine = e
	e.router.handle(context)
}

/**
设置HTML渲染map
*/
func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

/**
解析html
*/
func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

/**
开始运行
*/
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
