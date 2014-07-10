Macaron [![wercker status](https://app.wercker.com/status/282aa746d272d0eaa703a86852445a67/s "wercker status")](https://app.wercker.com/project/bykey/282aa746d272d0eaa703a86852445a67)
=======================

Package macaron is a high productive and modular design web framework in Go.

##### Current version: 0.0.1

## Getting Started

To install Macaron:

	go get github.com/Unknwon/macaron
	
The very basic usage of Macaron:

```
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

- Easy to plugin/unplugin features with modular design.
- Very simple steps to turn Martini middlewares to Macaron.
- Handy dependency injection powered by [inject](https://github.com/codegangsta/inject).
- Extreamly fast radix tree-based HTTP request router powered by [HttpRouter](https://github.com/julienschmidt/httprouter).

## Credits

- Basic design of [Martini](https://github.com/go-martini/martini).

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.