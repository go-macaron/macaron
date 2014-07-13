Macaron [![wercker status](https://app.wercker.com/status/282aa746d272d0eaa703a86852445a67/s "wercker status")](https://app.wercker.com/project/bykey/282aa746d272d0eaa703a86852445a67)
=======================

Package macaron is a high productive and modular design web framework in Go.

##### Current version: 0.0.2

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
	- Routing
	- Services
	- Serving Static Files
- Middleware Handlers
	- Next()
- Macaron Env
- FAQ

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

## Credits

- Basic design of [Martini](https://github.com/go-martini/martini).

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.