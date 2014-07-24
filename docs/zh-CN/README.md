## 目录

- [Macaron 核心](#macaron-%E6%A0%B8%E5%BF%83)
	- [处理器](#%E5%A4%84%E7%90%86%E5%99%A8)
	- [路由设置](#%E8%B7%AF%E7%94%B1%E8%AE%BE%E7%BD%AE)
	- [自定义服务](#%E8%87%AA%E5%AE%9A%E4%B9%89%E6%9C%8D%E5%8A%A1)
	- [静态文件](#%E9%9D%99%E6%80%81%E6%96%87%E4%BB%B6)
- [中间件机制](#%E4%B8%AD%E9%97%B4%E4%BB%B6%E6%9C%BA%E5%88%B6)
	- [Next()](#next)
	- [Gzip](#gzip)
	- [Render](#render)
	- [Cookie](#cookie)
- [多站点支持](#%E5%A4%9A%E7%AB%99%E7%82%B9%E6%94%AF%E6%8C%81)
- [Macaron 环境变量](#macaron-%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F)

## Macaron 核心

为了更快速的启用 Macaron, [macaron.Classic()](https://gowalker.org/github.com/Unknwon/macaron#Classic) 提供了一些默认的方便 Web 开发的工具:

```go
  m := macaron.Classic()
  // ... middleware and routing goes here
  m.Run()
```

下面是 Macaron 核心已经包含的功能  [macaron.Classic()](https://gowalker.org/github.com/Unknwon/macaron#Classic):

- Request/Response Logging（请求/响应日志）- [macaron.Logger](https://gowalker.org/github.com/Unknwon/macaron#Logger)
- Panic Recovery（容错恢复）- [macaron.Recovery](https://gowalker.org/github.com/Unknwon/macaron#Recovery)
- Static File serving（静态文件服务）- [macaron.Static](https://gowalker.org/github.com/Unknwon/macaron#Static)

### 处理器

处理器是 Macaron 的灵魂和核心所在. 一个处理器基本上可以是任何的函数:

```go
m.Get("/", func() {
	println("hello world")
})
```

#### 返回值

当一个处理器返回结果的时候, Macaron 将会把返回值作为字符串写入到当前的 [http.ResponseWriter](http://gowalker.org/net/http#ResponseWriter) 里面:

```go
m.Get("/", func() string {
	return "hello world" // HTTP 200 : "hello world"
})
```

另外你也可以选择性的返回状态码:

```go
m.Get("/", func() (int, string) {
	return 418, "i'm a teapot" // HTTP 418 : "i'm a teapot"
})
```

#### 注入服务

处理器是通过反射来调用的. Macaron 通过*Dependency Injection* *（依赖注入）* 来为处理器注入参数列表。 **这样使得 Macaron 与 Go 语言的 `http.HandlerFunc` 接口完全兼容** 

如果你加入一个参数到你的处理器, Macaron 将会搜索它参数列表中的服务，并且通过类型判断来解决依赖关系:

```go
m.Get("/", func(rw http.ResponseWriter, req *http.Request) { 
	// rw and req are injected by Macaron
	rw.WriteHeader(200) // HTTP 200
})
```

下面的这些服务已经被包含在核心 Macaron 中: [macaron.Classic()](https://gowalker.org/github.com/Unknwon/macaron#Classic):

- [*log.Logger](http://gowalker.org/log#Logger) - Macaron 全局日志器
- [*macaron.Context](https://gowalker.org/github.com/Unknwon/macaron#Context) - HTTP 请求上下文
- [http.ResponseWriter](http://gowalker.org/net/http/#ResponseWriter) - HTTP 响应结果的流接口
- [*http.Request](http://gowalker.org/net/http/#Request) - HTTP 请求

### 路由设置

在 Macaron 中, 路由是一个 HTTP 方法配对一个 URL 匹配模型. 每一个路由可以对应一个或多个处理器方法:

```go
m.Get("/", func() {
	// show something
})

m.Patch("/", func() {
	// update something
})

m.Post("/", func() {
	// create something
})

m.Put("/", func() {
	// replace something
})

m.Delete("/", func() {
	// destroy something
})

m.Options("/", func() {
	// http options
})

m.Any("/", func() {
	// do anything
})

m.Route("/", "GET,POST", func() {
	// combine something
})

m.NotFound(func() {
	// handle 404
})
```

路由匹配的顺序是按照他们被定义的顺序执行的. 最先被定义的路由将会首先被用户请求匹配并调用.

路由模型可能包含参数列表, 可以通过  [*Context.Params](https://gowalker.org/github.com/Unknwon/macaron#Context_Params) 来获取:

```go
m.Get("/hello/:name", func(ctx *macaron.Context) string {
	return "Hello " + ctx.Params(":name")
})
```

路由匹配可以通过全局匹配的形式:

```go
m.Get("/hello/*", func(ctx *macaron.Context) string {
	return "Hello " + ctx.Params("*")
})
```

您还可以使用使用正则表达式来书写路由规则：

```go
m.Get("/user/:username([\w]+)", func(ctx *macaron.Context) string {
	return fmt.Sprintf ("Hello %s", ctx.Params(":username"))
})

m.Get("/user/:username([\w]+)", func(ctx *macaron.Context) string {
	return fmt.Sprintf ("Hello %s", ctx.Params(":username"))
})

m.Get("/cms_:id([0-9]+).html", func(ctx *macaron.Context) string {
	return fmt.Sprintf ("The ID is %s", ctx.Params(":id"))
})
```

路由处理器可以被相互叠加使用, 例如很有用的地方可以是在验证和授权的时候:

```go
m.Get("/secret", authorize, func() {
	// this will execute as long as authorize doesn't write a response
})
```

路由还可以通过路由组来进行注册：

```go
m.Group("/books", func(r *macaron.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
    
    m.Group("/chapters", func(r *macaron.Router) {
	    r.Get("/:id", GetBooks)
	    r.Post("/new", NewBook)
	    r.Put("/update/:id", UpdateBook)
	    r.Delete("/delete/:id", DeleteBook)
	})
})
```

同样的，您可以为某一组路由设置集体的中间件：

```go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
}, MyMiddleware1, MyMiddleware2)
```

### 自定义服务

服务即是被注入到处理器中的参数. 你可以映射一个服务到 *全局* 或者 *请求* 的级别.

#### 全局映射

因为 Macaron 实现了 `inject.Injector` 的接口, 那么映射成为一个服务就非常简单:

```go
db := &MyDatabase{}
m := martini.Classic()
m.Map(db) // the service will be available to all handlers as *MyDatabase
// ...
m.Run()
```

#### 请求级别的映射

映射在请求级别的服务可以用 [*macaron.Context](https://gowalker.org/github.com/Unknwon/macaron#Context) 来完成:

```go
func MyCustomLoggerHandler(ctx *macaron.Context) {
	logger := &MyCustomLogger{ctx.Req}
	ctx.Map(logger) // mapped as *MyCustomLogger
}
```

#### 映射值到接口

关于服务最强悍的地方之一就是它能够映射服务到接口. 例如说, 假设你想要覆盖 [http.ResponseWriter](http://gowalker.org/net/http#ResponseWriter) 成为一个对象, 那么你可以封装它并包含你自己的额外操作, 你可以如下这样来编写你的处理器:

```go
func WrapResponseWriter(ctx *macaron.Context) {
	rw := NewSpecialResponseWriter(ctx.Resp)
	// override ResponseWriter with our wrapper ResponseWriter
	ctx.MapTo(rw, (*http.ResponseWriter)(nil)) 
}
```

### 静态文件

[macaron.Classic()](https://gowalker.org/github.com/Unknwon/macaron#Classic) 默认会服务位于你服务器环境根目录下的 "public" 文件夹。你可以通过加入 [macaron.Static](https://gowalker.org/github.com/Unknwon/macaron#Static) 的处理器来加入更多的静态文件服务的文件夹：

```go
m.Use(macaron.Static("assets")) // serve from the "assets" directory as well
```

## 中间件机制

中间件处理器是工作于请求和路由之间的. 本质上来说和 Macaron 其他的处理器没有分别. 你可以像如下这样添加一个中间件处理器到它的堆中:

```go
m.Use(func() {
  // do some middleware stuff
})
```

你可以通过 `Handlers` 函数对中间件堆有完全的控制. 它将会替换掉之前的任何设置过的处理器:

```go
m.Handlers(
	Middleware1,
	Middleware2,
	Middleware3,
)
```

中间件处理器可以非常好处理一些功能，像 logging(日志), authorization(授权), authentication(认证), sessions(会话), error pages(错误页面), 以及任何其他的操作需要在 HTTP 请求发生之前或者之后的:

```go
// validate an api key
m.Use(func(ctx *macaron.Context) {
	if ctx.Req.Header.Get("X-API-KEY") != "secret123" {
		ctx.Resp.WriteHeader(http.StatusUnauthorized)
	}
})
```

### Next()

[Context.Next()](https://gowalker.org/github.com/Unknwon/macaron#Context_Next) 是一个可选的函数用于中间件处理器暂时放弃执行直到其他的处理器都执行完毕. 这样就可以很好的处理在 HTTP 请求完成后需要做的操作：

```go
// log before and after a request
m.Use(func(ctx *macaron.Context, log *log.Logger){
	log.Println("before a request")

	ctx.Next()

	log.Println("after a request")
})
```

### Gzip

您需要在注册任何其它有响应输出的中间件之前注册 Gzip 中间件：

```go
m.Use(macaron.Gzip())
```

### Render

中间件 [macaron.Render](https://gowalker.org/github.com/Unknwon/macaron#Render) 已经被集成到 [*macaron.Context](https://gowalker.org/github.com/Unknwon/macaron#Context) 中，注册之后便可使用它：

```go
m.Use(macaron.Renderer(macaron.RenderOptions{}))
```

该中间件有一个可选的 [macaron.RenderOptions{}](https://gowalker.org/github.com/Unknwon/macaron#RenderOptions) 参数。之后，您可以就直接通过 [*macaron.Context](https://gowalker.org/github.com/Unknwon/macaron#Context) 来调用渲染方法，然后使用 `ctx.Data`来存储您需要渲染的模板变量：

```go
func Home(ctx *macaron.Context) {
	ctx.Data["title"] = "my home page"
	ctx.HTML(200, "home", ctx.Data)
}
```

### Cookie

最基本的 Cookie 用法：

- [ctx.SetCookie](https://gowalker.org/github.com/Unknwon/macaron#Context_SetCookie)
- [ctx.GetCookie](https://gowalker.org/github.com/Unknwon/macaron#Context_GetCookie)

如果需要更加安全的 Cookie 机制，可以先使用 [macaron.SetDefaultCookieSecret](https://gowalker.org/github.com/Unknwon/macaron#Macaron_SetDefaultCookieSecret) 设定密钥，然后使用：

- [ctx.SetSecureCookie](https://gowalker.org/github.com/Unknwon/macaron#Context_SetSecureCookie)
- [ctx.GetSecureCookie](https://gowalker.org/github.com/Unknwon/macaron#Context_GetSecureCookie)

这两个方法将会自动使用您设置的默认密钥进行加密/解密 Cookie 值。

对于那些对安全性要求特别高的应用，可以为每次设置 Cookie 使用不同的密钥加密/解密：

- [ctx.SetSuperSecureCookie](https://gowalker.org/github.com/Unknwon/macaron#Context_SetSuperSecureCookie)
- [ctx.GetSuperSecureCookie](https://gowalker.org/github.com/Unknwon/macaron#Context_GetSuperSecureCookie)

## 多站点支持

如果您想要运行 2 或 2 个以上的实例在一个程序里，[HostSwitcher](https://gowalker.org/github.com/Unknwon/macaron#HostSwitcher) 就是您需要的特性：

```go
func main() {
	m1 := macaron.Classic()
	// Register m1 middlewares and routers.

	m2 := macaron.Classic()
	// Register m2 middlewares and routers.

	hs := macaron.NewHostSwitcher()
	// Set instance corresponding to host address.
	hs.Set("gowalker.org", m1)
	hs.Set("gogs.io", m2)
	hs.Run()
}
```

## Macaron 环境变量

一些 Macaron 处理器依赖 `macaron.Env` 全局变量为开发模式和部署模式表现出不同的行为，不过更建议使用环境变量 `MACARON_ENV=production` 来指示当前的模式为部署模式。
