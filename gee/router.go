package gee

import (
	"log"
	"net/http"
	"strings"
)

/**
自定义的路由处理函数
*/
type HandlerFunc func(ctx *Context)

type router struct {
	handlers map[string]HandlerFunc
	root     map[string]*node
}

/**
初始化路由
*/
func NewRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc), root: make(map[string]*node)}
}

/**
解析全路径
*/
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

/**
添加路由参数
*/
func (r *router) addRoute(method, path string, handleFunc HandlerFunc) {
	key := method + "-" + path
	if _, ok := r.handlers[key]; !ok {
		log.Printf("成功添加路由 %s === %s", method, path)
		r.handlers[key] = handleFunc
	}

	parts := parsePattern(path)
	if item, ok := r.root[method]; !ok { // node不存在
		item = &node{}
		r.root[method] = item
	}

	r.root[method].insert(path, parts, 0) // 插入路由前缀树
}

/**
获取路由信息
*/
func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchPart := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.root[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchPart, 0)
	if n != nil { // 有值
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' { // 以：开头需要动态绑定参数
				params[part[1:]] = searchPart[index]
			}

			if part[0] == '*' && len(part) > 1 { // 静态文件路径
				params[part[1:]] = strings.Join(searchPart[index:], "/")
				break
			}
		}

		return n, params
	}

	return nil, nil
}

/**
路由匹配
*/
func (r *router) handle(ctx *Context) {
	node, params := r.getRoute(ctx.Method, ctx.Path)
	if node != nil {
		ctx.Params = params
		key := ctx.Method + "-" + node.pattern
		ctx.handlers = append(ctx.handlers, r.handlers[key])
	} else {
		ctx.handlers = append(ctx.handlers, func(ctx *Context) {
			_, _ = ctx.String(http.StatusNotFound, "当前路由未找到，method:%s, path:%s", ctx.Method, ctx.Path)
		})
	}

	ctx.Next()
}
