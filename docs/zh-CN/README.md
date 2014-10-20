## 目录

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

m.Combo("/").
	Get(func() string { return "GET" }).
	Patch(func() string { return "PATCH" }).
	Post(func() string { return "POST" }).
	Put(func() string { return "PUT" }).
	Delete(func() string { return "DELETE" }).
	Options(func() string { return "OPTIONS" }).
	Head(func() string { return "HEAD" })

m.NotFound(func() {
	// handle 404
})
```

路由匹配的顺序是按照他们被定义的顺序执行的. 最先被定义的路由将会首先被用户请求匹配并调用。

如果您想要使用子路径但让路由代码保持简洁，可以调用 `m.SetURLPrefix(suburl)`。

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
