Macaron [![wercker status](https://app.wercker.com/status/282aa746d272d0eaa703a86852445a67/s "wercker status")](https://app.wercker.com/project/bykey/282aa746d272d0eaa703a86852445a67)
=======================

Package macaron is a high productive and modular design web framework in Go.

##### Current version: 0.0.4

Anyone who uses [Martini](https://github.com/go-martini/martini) and familiar with dependency injection like me, should be very comfortable about how to use Macaron.

## Getting Started

To install Macaron:

	go get github.com/Unknwon/macaron
	
The very basic usage of Macaron:

```go
package main

import "github.com/Unknwon/macaron"

func main() {
  m := macaron.Classic()
  m.Get("/", func() string {
    return "Hello world!"
  })
  m.Run()
}
```

## Getting Help

- Visit [Go Walker](https://gowalker.org/github.com/Unknwon/macaron) for API documentation.

## Features

- Serve multiple sites in one program.
- Unlimited nested group routers.
- Easy to plugin/unplugin features with modular design.
- Integrated most frequently used middlewares with less reflection.
- Very simple steps to turn Martini middlewares to Macaron.
- Handy dependency injection powered by [inject](https://github.com/codegangsta/inject).
- Extreamly fast radix tree-based HTTP request router powered by [HttpRouter](https://github.com/julienschmidt/httprouter).

----------

## Table of Contents

- [Classic Macaron](#classic-macaron)
	- [Handlers](#handlers)
	- [Routing](#routing)
	- [Services](#services)
	- [Serving Static Files](#serving-static-files)
- [Middleware Handlers](#middleware-handlers)
	- [Next()](#next)
	- [Gzip](#gzip)
	- [Render](#render)
	- [Cookie](#cookie)
- [Macaron Env](#macaron-env)
- [FAQ](#faq)

## Classic Macaron

To get up and running quickly, [macaron.Classic()](https://gowalker.org/github.com/Unknwon/macaron#Classic) provides some reasonable defaults that work well for most web applications:

```go
  m := macaron.Classic()
  // ... middleware and routing goes here
  m.Run()
```

Below is some of the functionality [macaron.Classic()](https://gowalker.org/github.com/Unknwon/macaron#Classic) pulls in automatically:

- Request/Response Logging - [macaron.Logger](https://gowalker.org/github.com/Unknwon/macaron#Logger)
- Panic Recovery - [macaron.Recovery](https://gowalker.org/github.com/Unknwon/macaron#Recovery)
- Static File serving - [macaron.Static](https://gowalker.org/github.com/Unknwon/macaron#Static)

### Handlers

Handlers are the heart and soul of Macaron. A handler is basically any kind of callable function:

```go
m.Get("/", func() {
  println("hello world")
})
```

#### Return Values

If a handler returns something, Macaron will write the result to the current [http.ResponseWriter](http://gowalker.org/net/http#ResponseWriter) as a string:

```go
m.Get("/", func() string {
  return "hello world" // HTTP 200 : "hello world"
})
```

You can also optionally return a status code:

```go
m.Get("/", func() (int, string) {
  return 418, "i'm a teapot" // HTTP 418 : "i'm a teapot"
})
```

#### Service Injection

Handlers are invoked via reflection. Macaron makes use of *Dependency Injection* to resolve dependencies in a Handlers argument list. **This makes Macaron completely  compatible with golang's `http.HandlerFunc` interface.**

If you add an argument to your Handler, Martini will search its list of services and attempt to resolve the dependency via type assertion:

```go
m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Macaron
  res.WriteHeader(200) // HTTP 200
})
```

The following services are included with [macaron.Classic()](https://gowalker.org/github.com/Unknwon/macaron#Classic):

- [*log.Logger](http://gowalker.org/log#Logger) - Global logger for Macaron.
- [*macaron.Context](https://gowalker.org/github.com/Unknwon/macaron#Context) - HTTP request context.
- [http.ResponseWriter](http://gowalker.org/net/http/#ResponseWriter) - HTTP Response writer interface.
- [*http.Request](http://gowalker.org/net/http/#Request) - HTTP Request.

### Routing

In Macaron, a route is an HTTP method paired with a URL-matching pattern.
Each route can take one or more handler methods:

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

m.NotFound(func() {
  // handle 404
})
```

Routes are matched in the order they are defined. The first route that
matches the request is invoked.

Route patterns may include named parameters, accessible via the [*Context.Params](https://gowalker.org/github.com/Unknwon/macaron#Context_Params):

```go
m.Get("/hello/:name", func(ctx *macaron.Context) string {
  return "Hello " + ctx.Params("name")
})
```

Routes can be matched with globs:

```go
m.Get("/hello/*name", func(ctx *macaron.Context) string {
  return "Hello " + ctx.Params("name")
})
```

Regular expressions are currently **NOT SUPPORTED**.

Route handlers can be stacked on top of each other, which is useful for things like authentication and authorization:

```go
m.Get("/secret", authorize, func() {
  // this will execute as long as authorize doesn't write a response
})
```

Route groups can be added too using the Group method.

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

Just like you can pass middlewares to a handler you can pass middlewares to groups.

```go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
}, MyMiddleware1, MyMiddleware2)
```

### Services

Services are objects that are available to be injected into a Handler's argument list. You can map a service on a *Global* or *Request* level.

#### Global Mapping

A Macaron instance implements the inject.Injector interface, so mapping a service is easy:

```go
db := &MyDatabase{}
m := martini.Classic()
m.Map(db) // the service will be available to all handlers as *MyDatabase
// ...
m.Run()
```

#### Request-Level Mapping

Mapping on the request level can be done in a handler via [*macaron.Context](https://gowalker.org/github.com/Unknwon/macaron#Context):

```go
func MyCustomLoggerHandler(ctx *macaron.Context, req *http.Request) {
  logger := &MyCustomLogger{req}
  ctx.Map(logger) // mapped as *MyCustomLogger
}
```

#### Mapping values to Interfaces

One of the most powerful parts about services is the ability to map a service to an interface. For instance, if you wanted to override the [http.ResponseWriter](http://gowalker.org/net/http#ResponseWriter) with an object that wrapped it and performed extra operations, you can write the following handler:

```go
func WrapResponseWriter(res http.ResponseWriter, ctx *macaron.Context) {
  rw := NewSpecialResponseWriter(res)
  ctx.MapTo(rw, (*http.ResponseWriter)(nil)) // override ResponseWriter with our wrapper ResponseWriter
}
```

### Serving Static Files

A [macaron.Classic()](https://gowalker.org/github.com/Unknwon/macaron#Classic) instance automatically serves static files from the "public" directory in the root of your server.
You can serve from more directories by adding more [macaron.Static](https://gowalker.org/github.com/Unknwon/macaron#Static) handlers.

```go
m.Use(macaron.Static("assets")) // serve from the "assets" directory as well
```

## Middleware Handlers

Middleware Handlers sit between the incoming http request and the router. In essence they are no different than any other Handler in Macaron. You can add a middleware handler to the stack like so:

```go
m.Use(func() {
  // do some middleware stuff
})
```

You can have full control over the middleware stack with the `Handlers` function. This will replace any handlers that have been previously set:

```go
m.Handlers(
  Middleware1,
  Middleware2,
  Middleware3,
)
```

Middleware Handlers work really well for things like logging, authorization, authentication, sessions, gzipping, error pages and any other operations that must happen before or after an http request:

```go
// validate an api key
m.Use(func(res http.ResponseWriter, req *http.Request) {
  if req.Header.Get("X-API-KEY") != "secret123" {
    res.WriteHeader(http.StatusUnauthorized)
  }
})
```

### Next()

[Context.Next()](https://gowalker.org/github.com/Unknwon/macaron#Context_Next) is an optional function that Middleware Handlers can call to yield the until after the other Handlers have been executed. This works really well for any operations that must happen after an http request:

```go
// log before and after a request
m.Use(func(ctx *macaron.Context, log *log.Logger){
  log.Println("before a request")

  ctx.Next()

  log.Println("after a request")
})
```

### Gzip

Register middleware Gzip before all the other middlewares that have response.

```go
m.Use(macaron.Gzip())
```

### Render

The [*macaron.Render](https://gowalker.org/github.com/Unknwon/macaron#Render) has been integrated into [*macaron.Context](https://gowalker.org/github.com/Unknwon/macaron#Context). To use it, you have to register the render middleware first.

```go
m.Use(macaron.Renderer(macaron.RenderOptions{}))
```

Note that [macaron.RenderOptions{}](https://gowalker.org/github.com/Unknwon/macaron#RenderOptions) is optional. After that, you can directly call render methods by [*macaron.Context](https://gowalker.org/github.com/Unknwon/macaron#Context), and use `ctx.Data` to store template variables.

```go
func Home(ctx *macaron.Context) {
	ctx.Data["title"] = "my home page"
	ctx.HTML(200, "home", ctx.Data)
}
```

### Cookie

The very basic usage of cookie is just [ctx.SetCookie](https://gowalker.org/github.com/Unknwon/macaron#Context_SetCookie) and [ctx.GetCookie](https://gowalker.org/github.com/Unknwon/macaron#Context_GetCookie).

And there are more secure cookie support. First, you need to call [macaron.SetDefaultCookieSecret](https://gowalker.org/github.com/Unknwon/macaron#Macaron_SetDefaultCookieSecret), then use it by calling [ctx.SetSecureCookie](https://gowalker.org/github.com/Unknwon/macaron#Context_SetSecureCookie) and [ctx.GetSecureCookie](https://gowalker.org/github.com/Unknwon/macaron#Context_GetSecureCookie). These two methods uses default secret string you set globally to encode and decode values.

For people who wants even more secure cookies that change secret string every time, just use [ctx.SetSuperSecureCookie](https://gowalker.org/github.com/Unknwon/macaron#Context_SetSuperSecureCookie) and [ctx.GetSuperSecureCookie](https://gowalker.org/github.com/Unknwon/macaron#Context_GetSuperSecureCookie).

## Macaron Env

Some Macaron handlers make use of the `macaron.Env` global variable to provide special functionality for development environments vs production environments. It is recommended that the `MACARON_ENV=production` environment variable to be set when deploying a Martini server into a production environment.

## FAQ

### Where do I find middleware X?

Start by looking in the [macaron-contrib](https://github.com/macaron-contrib) projects. If it is not there feel free to contact a [macaron-contrib](https://github.com/macaron-contrib) team member about adding a new repo to the organization.

- [i18n](https://github.com/macaron-contrib/i18n) - Internationalization and Localization
- [cache](https://github.com/macaron-contrib/cache) - Cache manager
- [session](https://github.com/macaron-contrib/session) - Session manager

### How do I integrate with existing servers?

A Martini instance implements `http.Handler`, so it can easily be used to serve subtrees
on existing Go servers. For example this is a working Martini app for Google App Engine:

```go
package hello

import (
  "net/http"
  "github.com/Unknwon/macaron"
)

func init() {
  m := macaron.Classic()
  m.Get("/", func() string {
    return "Hello world!"
  })
  http.Handle("/", m)
}
```

### How do I change the port/host?

Macaron's `Run` function looks for the PORT and HOST environment variables and uses those. Otherwise Macaro will default to `localhost:4000`.
To have more flexibility over port and host, use the `http.ListenAndServe` function instead.

```go
  m := macaro.Classic()
  // ...
  log.Fatal(http.ListenAndServe(":8080", m))
```

### What's the idea behind this other than Martini?

- Integrate frequently used middlewares and helper methods with less reflection.
- Make a deep source study against Martini.

### Live code reload?

[Bra](https://github.com/Unknwon/bra) is the prefect fit for live reloading Macaron apps.

## Credits

- Basic design of [Martini](https://github.com/go-martini/martini).

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.