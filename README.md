Macaron [![wercker status](https://app.wercker.com/status/282aa746d272d0eaa703a86852445a67/s "wercker status")](https://app.wercker.com/project/bykey/282aa746d272d0eaa703a86852445a67)
=======================

Package macaron is a high productive and modular design web framework in Go.

##### Current version: 0.1.2

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

## Features

- Powerful routing.
- Serve multiple sites in one program.
- Unlimited nested group routers.
- Easy to plugin/unplugin features with modular design.
- Integrated most frequently used middlewares with less reflection.
- Very simple steps to turn Martini middlewares to Macaron.
- Handy dependency injection powered by [inject](https://github.com/codegangsta/inject).

## Getting Help

- Visit [Go Walker](https://gowalker.org/github.com/Unknwon/macaron) for API documentation.
- Documentation
	- [简体中文](docs/zh-CN)
	- [English](docs/en-US)

## FAQs

### Where do I find middleware X?

Start by looking in the [macaron-contrib](https://github.com/macaron-contrib) projects. If it is not there feel free to contact a [macaron-contrib](https://github.com/macaron-contrib) team member about adding a new repo to the organization.

- [i18n](https://github.com/macaron-contrib/i18n) - Internationalization and Localization
- [cache](https://github.com/macaron-contrib/cache) - Cache manager
- [session](https://github.com/macaron-contrib/session) - Session manager
- [csrf](https://github.com/macaron-contrib/csrf) - Generates and validates csrf tokens
- [captcha](https://github.com/macaron-contrib/captcha) - Captcha service
- [pongo2](https://github.com/macaron-contrib/pongo2) - Pongo2 template engine support

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
- Router layer of [beego](https://github.com/astaxie/beego).

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.