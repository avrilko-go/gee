package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Method     string
	Path       string
	Writer     http.ResponseWriter
	Req        *http.Request
	Params     map[string]string
	StatusCode int
	handlers   []HandlerFunc // 中间件处理
	index      int           // 中间件处理游标 默认为-1
	engine     *Engine
}

type H map[string]interface{}

/**
初始化一个上下文对象
*/
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Method: r.Method,
		Path:   r.URL.Path,
		Writer: w,
		Req:    r,
		Params: make(map[string]string),
		index:  -1,
	}
}

/**
获取url上的查询参数
*/
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

/**
错误处理
*/
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	// 直接跳过其他中间件的执行了
	_ = c.JSON(code, H{"message": err})
}

/**
查询表单发送过来的数据
*/
func (c *Context) PostForm(key string) string {
	return c.Req.PostFormValue(key)
}

/**
设置http请求返回的头部信息
*/
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

/**
设置http的状态码
*/
func (c *Context) SetHttpStatus(code int) {
	c.Writer.WriteHeader(code)
	c.StatusCode = code
}

/**
返回json数据
*/
func (c *Context) JSON(code int, data interface{}) error {
	c.SetHttpStatus(code)
	c.SetHeader("Content-Type", "application/json")
	encode := json.NewEncoder(c.Writer)
	return encode.Encode(data)
}

/**
返回HTML信息
*/
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHttpStatus(code)
	c.SetHeader("Content-Type", "text/html")
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.SetHttpStatus(http.StatusInternalServerError)
	}
}

/**
格式化字符串显示
*/
func (c *Context) String(code int, format string, value ...interface{}) (int, error) {
	c.SetHttpStatus(code)
	c.SetHeader("Content-Type", "text/plain")
	return c.Writer.Write([]byte(fmt.Sprintf(format, value...)))
}

/**
返回数据
*/
func (c *Context) Data(code int, data []byte) (int, error) {
	c.SetHttpStatus(code)
	return c.Writer.Write(data)
}

/**
获取动态路由URL上的参数
*/
func (c *Context) Param(key string) string {
	return c.Params[key]
}

/**
中间件处理函数
*/
func (c *Context) Next() {
	c.index++
	if c.index > len(c.handlers) {
		return
	}
	c.handlers[c.index](c)
}
