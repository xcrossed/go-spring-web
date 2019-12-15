/*
 * Copyright 2012-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package SpringWeb

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/go-spring/go-spring-parent/spring-const"
	"github.com/go-spring/go-spring-parent/spring-logger"
)

//
// 定义 Web 处理函数
//
type Handler func(WebContext)

//
// 路由映射器
//
type Mapper struct {
	Method  string
	Path    string
	Handler Handler
	Filters []Filter
}

func NewMapper(method string, path string, fn Handler, filters []Filter) Mapper {
	return Mapper{
		Method:  method,
		Path:    path,
		Handler: fn,
		Filters: filters,
	}
}

//
// 路由表
//
type WebMapper interface {
	// 获取路由表
	GetMapper() map[string]Mapper

	// 通过路由分组注册 Web 处理函数
	Route(path string, filters ...Filter) *Route

	// 通过路由分组注册 Web 处理函数
	Group(path string, fn GroupHandler, filters ...Filter)

	// 注册 GET 方法处理函数
	GET(path string, fn Handler, filters ...Filter)

	// 注册 POST 方法处理函数
	POST(path string, fn Handler, filters ...Filter)

	// 注册 PATCH 方法处理函数
	PATCH(path string, fn Handler, filters ...Filter)

	// 注册 PUT 方法处理函数
	PUT(path string, fn Handler, filters ...Filter)

	// 注册 DELETE 方法处理函数
	DELETE(path string, fn Handler, filters ...Filter)

	// 注册 HEAD 方法处理函数
	HEAD(path string, fn Handler, filters ...Filter)

	// 注册 OPTIONS 方法处理函数
	OPTIONS(path string, fn Handler, filters ...Filter)
}

//
// 定义 Web 路由分组。分组的限制：分组内路由只能共享相同的 filters。
//
type Route struct {
	basePath string
	filters  []Filter
	mapper   WebMapper
}

//
// 工厂函数
//
func NewRoute(mapper WebMapper, path string, filters []Filter) *Route {
	return &Route{
		mapper:   mapper,
		basePath: path,
		filters:  filters,
	}
}

//
// 定义分组处理函数
//
type GroupHandler func(*Route)

func (g *Route) GET(path string, fn Handler) *Route {
	g.mapper.GET(g.basePath+path, fn, g.filters...)
	return g
}

func (g *Route) POST(path string, fn Handler) *Route {
	g.mapper.POST(g.basePath+path, fn, g.filters...)
	return g
}

func (g *Route) PATCH(path string, fn Handler) *Route {
	g.mapper.PATCH(g.basePath+path, fn, g.filters...)
	return g
}

func (g *Route) PUT(path string, fn Handler) *Route {
	g.mapper.PUT(g.basePath+path, fn, g.filters...)
	return g
}

func (g *Route) DELETE(path string, fn Handler) *Route {
	g.mapper.DELETE(g.basePath+path, fn, g.filters...)
	return g
}

func (g *Route) HEAD(path string, fn Handler) *Route {
	g.mapper.HEAD(g.basePath+path, fn, g.filters...)
	return g
}

func (g *Route) OPTIONS(path string, fn Handler) *Route {
	g.mapper.OPTIONS(g.basePath+path, fn, g.filters...)
	return g
}

//
// 定义 Web 容器接口
//
type WebContainer interface {
	// 监听的 IP
	GetIP() string
	SetIP(ip string)

	// 监听的 Port
	GetPort() []int
	SetPort(port ...int)

	// 是否启用 SSL
	EnableSSL() bool
	SetEnableSSL(enable bool)

	// SSL 证书
	GetKeyFile() string
	SetKeyFile(keyFile string)
	GetCertFile() string
	SetCertFile(certFile string)

	// 启动 Web 容器，非阻塞
	Start()

	// 停止 Web 容器
	Stop(ctx context.Context)

	// 继承路由表的方法
	WebMapper

	// 获取 IoC 容器里面注册的 Filter 对象
	Filters(s ...string) []Filter
}

//
// WebContainer 基本实现
//
type BaseWebContainer struct {
	ip        string
	port      []int
	enableSSL bool
	keyFile   string
	certFile  string
	mapper    map[string]Mapper
}

func (c *BaseWebContainer) Init() {
	c.mapper = make(map[string]Mapper)
}

func (c *BaseWebContainer) GetIP() string {
	return c.ip
}

func (c *BaseWebContainer) SetIP(ip string) {
	c.ip = ip
}

func (c *BaseWebContainer) GetPort() []int {
	return c.port
}

func (c *BaseWebContainer) SetPort(port ...int) {
	c.port = port
}

func (c *BaseWebContainer) EnableSSL() bool {
	return c.enableSSL
}

func (c *BaseWebContainer) SetEnableSSL(enable bool) {
	c.enableSSL = enable
}

func (c *BaseWebContainer) GetKeyFile() string {
	return c.keyFile
}

func (c *BaseWebContainer) SetKeyFile(keyFile string) {
	c.keyFile = keyFile
}

func (c *BaseWebContainer) GetCertFile() string {
	return c.certFile
}

func (c *BaseWebContainer) SetCertFile(certFile string) {
	c.certFile = certFile
}

func (c *BaseWebContainer) GetMapper() map[string]Mapper {
	return c.mapper
}

func (c *BaseWebContainer) Route(path string, filters ...Filter) *Route {
	return NewRoute(c, path, filters)
}

func (c *BaseWebContainer) Group(path string, fn GroupHandler, filters ...Filter) {
	fn(NewRoute(c, path, filters))
}

func (c *BaseWebContainer) GET(path string, fn Handler, filters ...Filter) {
	c.mapper[path] = NewMapper("GET", path, fn, filters)
}

func (c *BaseWebContainer) PATCH(path string, fn Handler, filters ...Filter) {
	c.mapper[path] = NewMapper("PATCH", path, fn, filters)
}

func (c *BaseWebContainer) PUT(path string, fn Handler, filters ...Filter) {
	c.mapper[path] = NewMapper("PUT", path, fn, filters)
}

func (c *BaseWebContainer) POST(path string, fn Handler, filters ...Filter) {
	c.mapper[path] = NewMapper("POST", path, fn, filters)
}

func (c *BaseWebContainer) DELETE(path string, fn Handler, filters ...Filter) {
	c.mapper[path] = NewMapper("DELETE", path, fn, filters)
}

func (c *BaseWebContainer) HEAD(path string, fn Handler, filters ...Filter) {
	c.mapper[path] = NewMapper("HEAD", path, fn, filters)
}

func (c *BaseWebContainer) OPTIONS(path string, fn Handler, filters ...Filter) {
	c.mapper[path] = NewMapper("OPTIONS", path, fn, filters)
}

func (c *BaseWebContainer) Filters(s ...string) []Filter {
	panic(SpringConst.UNIMPLEMENTED_METHOD)
}

//
// 定义 Web 上下文接口，设计理念：为社区中优秀的 Web 服务器提供一个抽象层，使得
// 底层可以灵活切换，因此在功能上取这些 Web 服务器功能的交集，同时提供获取底层对
// 象的接口，以便在不能满足用户要求的时候使用底层实现的能力，当然这种功能要慎用。
//
type WebContext interface {
	/////////////////////////////////////////
	// 通用能力部分

	SpringLogger.LoggerContext

	// 获取封装的底层上下文对象
	NativeContext() interface{}

	// Get retrieves data from the context.
	Get(key string) interface{}

	// Set saves data in the context.
	Set(key string, val interface{})

	/////////////////////////////////////////
	// Request Part

	// Request returns `*http.Request`.
	Request() *http.Request

	// IsTLS returns true if HTTP connection is TLS otherwise false.
	IsTLS() bool

	// IsWebSocket returns true if HTTP connection is WebSocket otherwise false.
	IsWebSocket() bool

	// Scheme returns the HTTP protocol scheme, `http` or `https`.
	Scheme() string

	// ClientIP implements a best effort algorithm to return the real client IP,
	// it parses X-Real-IP and X-Forwarded-For in order to work properly with
	// reverse-proxies such us: nginx or haproxy. Use X-Forwarded-For before
	// X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
	ClientIP() string

	// Path returns the registered path for the handler.
	Path() string

	// Handler returns the matched handler by router.
	Handler() Handler

	// ContentType returns the Content-Type header of the request.
	ContentType() string

	// GetHeader returns value from request headers.
	GetHeader(key string) string

	// GetRawData return stream data.
	GetRawData() ([]byte, error)

	// Param returns path parameter by name.
	PathParam(name string) string

	// ParamNames returns path parameter names.
	PathParamNames() []string

	// ParamValues returns path parameter values.
	PathParamValues() []string

	// QueryParam returns the query param for the provided name.
	QueryParam(name string) string

	// QueryParams returns the query parameters as `url.Values`.
	QueryParams() url.Values

	// QueryString returns the URL query string.
	QueryString() string

	// FormValue returns the form field value for the provided name.
	FormValue(name string) string

	// FormParams returns the form parameters as `url.Values`.
	FormParams() (url.Values, error)

	// FormFile returns the multipart form file for the provided name.
	FormFile(name string) (*multipart.FileHeader, error)

	// SaveUploadedFile uploads the form file to specific dst.
	SaveUploadedFile(file *multipart.FileHeader, dst string) error

	// MultipartForm returns the multipart form.
	MultipartForm() (*multipart.Form, error)

	// Cookie returns the named cookie provided in the request.
	Cookie(name string) (*http.Cookie, error)

	// Cookies returns the HTTP cookies sent with the request.
	Cookies() []*http.Cookie

	// Bind binds the request body into provided type `i`. The default binder
	// does it based on Content-Type header.
	Bind(i interface{}) error

	/////////////////////////////////////////
	// Response Part

	// ResponseWriter returns `http.ResponseWriter`.
	ResponseWriter() http.ResponseWriter

	// Status sets the HTTP response code.
	Status(code int)

	// Header is a intelligent shortcut for c.Writer.Header().Set(key, value).
	// It writes a header in the response.
	// If value == "", this method removes the header `c.Writer.Header().Del(key)`
	Header(key, value string)

	// SetCookie adds a `Set-Cookie` header in HTTP response.
	SetCookie(cookie *http.Cookie)

	// NoContent sends a response with no body and a status code.
	NoContent(code int)

	// String writes the given string into the response body.
	String(code int, format string, values ...interface{})

	// HTML sends an HTTP response with status code.
	HTML(code int, html string)

	// HTMLBlob sends an HTTP blob response with status code.
	HTMLBlob(code int, b []byte)

	// JSON sends a JSON response with status code.
	JSON(code int, i interface{})

	// JSONPretty sends a pretty-print JSON with status code.
	JSONPretty(code int, i interface{}, indent string)

	// JSONBlob sends a JSON blob response with status code.
	JSONBlob(code int, b []byte)

	// JSONP sends a JSONP response with status code. It uses `callback`
	// to construct the JSONP payload.
	JSONP(code int, callback string, i interface{})

	// JSONPBlob sends a JSONP blob response with status code. It uses
	// `callback` to construct the JSONP payload.
	JSONPBlob(code int, callback string, b []byte)

	// XML sends an XML response with status code.
	XML(code int, i interface{})

	// XMLPretty sends a pretty-print XML with status code.
	XMLPretty(code int, i interface{}, indent string)

	// XMLBlob sends an XML blob response with status code.
	XMLBlob(code int, b []byte)

	// Blob sends a blob response with status code and content type.
	Blob(code int, contentType string, b []byte)

	// Stream sends a streaming response with status code and content type.
	Stream(code int, contentType string, r io.Reader)

	// File sends a response with the content of the file.
	File(file string)

	// Attachment sends a response as attachment, prompting client to save the
	// file.
	Attachment(file string, name string)

	// Inline sends a response as inline, opening the file in the browser.
	Inline(file string, name string)

	// Redirect redirects the request to a provided URL with status code.
	Redirect(code int, url string)

	// SSEvent writes a Server-Sent Event into the body stream.
	SSEvent(name string, message interface{})

	/////////////////////////////////////////
	// 错误处理部分

	// Error invokes the registered HTTP error handler.
	Error(err error)
}

//
// 定义 Web 过滤器
//
type Filter interface {
	// 函数内部通过 chain.Next() 驱动链条向后执行
	Invoke(ctx WebContext, chain *FilterChain)
}

//
// 包装 Web 处理函数的过滤器
//
type handlerFilter struct {
	fn Handler
}

func (h *handlerFilter) Invoke(ctx WebContext, _ *FilterChain) {
	h.fn(ctx)
}

//
// 把 Web 处理函数转换成 Web 过滤器
//
func HandlerFilter(fn Handler) Filter {
	return &handlerFilter{
		fn: fn,
	}
}

//
// 定义 Web 过滤器链条
//
type FilterChain struct {
	filters []Filter
	next    int
}

//
// 工厂函数
//
func NewFilterChain(filters []Filter) *FilterChain {
	return &FilterChain{
		filters: filters,
	}
}

//
// 执行下一个 Web 过滤器
//
func (chain *FilterChain) Next(ctx WebContext) {
	if chain.next >= len(chain.filters) {
		return
	}
	f := chain.filters[chain.next]
	chain.next++
	f.Invoke(ctx, chain)
}

//
// 执行 Web 处理函数
//
func InvokeHandler(ctx WebContext, fn Handler, filters []Filter) {
	if len(filters) > 0 {
		filters = append(filters, HandlerFilter(fn))
		chain := NewFilterChain(filters)
		chain.Next(ctx)
	} else {
		fn(ctx)
	}
}

//
// 定义 WebContainer 的工厂函数
//
type Factory func() WebContainer

//
// 保存 WebContainer 的工厂函数
//
var WebContainerFactory Factory

//
// 注册 WebContainer 的工厂函数
//
func RegisterWebContainerFactory(fn Factory) {
	WebContainerFactory = fn
}
