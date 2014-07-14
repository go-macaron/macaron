// Copyright 2014 Unknown
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package macaron is a high productive and modular design web framework in Go.
package macaron

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"github.com/julienschmidt/httprouter"

	"github.com/Unknwon/macaron/inject"
)

func Version() string {
	return "0.0.4.0714"
}

// Handler can be any callable function.
// Macaron attempts to inject services into the handler's argument list,
// and panics if an argument could not be fullfilled via dependency injection.
type Handler interface{}

// validateHandler makes sure a handler is a callable function,
// and panics if it is not.
func validateHandler(handler Handler) {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic("mocaron handler must be a callable function")
	}
}

// Macaron represents the top level web application.
// inject.Injector methods can be invoked to map services on a global level.
type Macaron struct {
	inject.Injector
	handlers []Handler
	action   Handler
	*Router
	logger *log.Logger
}

// New creates a bare bones Macaron instance.
// Use this method if you want to have full control over the middleware that is used.
func New() *Macaron {
	m := &Macaron{
		Injector: inject.New(),
		action:   func() {},
		Router: &Router{
			router: httprouter.New(),
		},
		logger: log.New(os.Stdout, "[Macaron] ", 0),
	}
	m.Router.m = m
	m.Map(m.logger)
	m.Map(defaultReturnHandler())
	m.router.NotFound = func(resp http.ResponseWriter, req *http.Request) {
		c := m.createContext(resp, req)
		c.handlers = append(m.handlers, func(resp http.ResponseWriter) (int, string) {
			return 404, "404 Not Found"
		})
		c.run()
	}
	return m
}

// Classic creates a classic Macaron with some basic default middleware:
// mocaron.Logger, mocaron.Recovery and mocaron.Static.
func Classic() *Macaron {
	m := New()
	m.Use(Logger())
	m.Use(Recovery())
	m.Use(Static("public"))
	return m
}

// Handlers sets the entire middleware stack with the given Handlers.
// This will clear any current middleware handlers,
// and panics if any of the handlers is not a callable function
func (m *Macaron) Handlers(handlers ...Handler) {
	m.handlers = make([]Handler, 0)
	for _, handler := range handlers {
		m.Use(handler)
	}
}

// Action sets the handler that will be called after all the middleware has been invoked.
// This is set to macaron.Router in a macaron.Classic().
func (m *Macaron) Action(handler Handler) {
	validateHandler(handler)
	m.action = handler
}

// Use adds a middleware Handler to the stack,
// and panics if the handler is not a callable func.
// Middleware Handlers are invoked in the order that they are added.
func (m *Macaron) Use(handler Handler) {
	validateHandler(handler)
	m.handlers = append(m.handlers, handler)
}

func (m *Macaron) createContext(resp http.ResponseWriter, req *http.Request) *Context {
	c := &Context{
		Injector: inject.New(),
		handlers: m.handlers,
		action:   m.action,
		rw:       NewResponseWriter(resp),
		index:    0,
		Req:      req,
		Resp:     resp,
	}
	c.SetParent(m)
	c.Map(c)
	c.MapTo(c.rw, (*http.ResponseWriter)(nil))
	c.Map(req)
	return c
}

// ServeHTTP is the HTTP Entry point for a Macaron instance.
// Useful if you want to control your own HTTP server.
// Be aware that none of middleware will run without registering any router.
func (m *Macaron) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	m.router.ServeHTTP(resp, req)
}

// getDefaultListenAddr returns default server listen address of Macaron.
func getDefaultListenAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	host := os.Getenv("HOST")
	return host + ":" + port
}

// Run the http server. Listening on os.GetEnv("PORT") or 4000 by default.
func (m *Macaron) Run() {
	addr := getDefaultListenAddr()

	logger := m.Injector.GetVal(reflect.TypeOf(m.logger)).Interface().(*log.Logger)
	logger.Printf("listening on %s (%s)\n", addr, Env)
	logger.Fatalln(http.ListenAndServe(addr, m))
}

// __________               __
// \______   \ ____  __ ___/  |_  ___________
//  |       _//  _ \|  |  \   __\/ __ \_  __ \
//  |    |   (  <_> )  |  /|  | \  ___/|  | \/
//  |____|_  /\____/|____/ |__|  \___  >__|
//         \/                        \/

// Router represents a Macaron router layer.
type Router struct {
	m      *Macaron
	router *httprouter.Router
	prefx  string
	groups []group
}

type group struct {
	pattern  string
	handlers []Handler
}

// Handle registers a new request handle with the given pattern, method and handlers.
func (r *Router) Handle(method string, pattern string, handlers []Handler) {
	if len(r.groups) > 0 {
		groupPattern := ""
		h := make([]Handler, 0)
		for _, g := range r.groups {
			groupPattern += g.pattern
			h = append(h, g.handlers...)
		}

		pattern = groupPattern + pattern
		h = append(h, handlers...)
		handlers = h
	}

	r.router.Handle(method, pattern, func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		c := r.m.createContext(resp, req)
		c.params = params
		c.handlers = append(r.m.handlers, handlers...)
		c.run()
	})
}

func (r *Router) Group(pattern string, fn func(*Router), h ...Handler) {
	r.groups = append(r.groups, group{pattern, h})
	fn(r)
	r.groups = r.groups[:len(r.groups)-1]
}

// Get is a shortcut for r.Handle("GET", pattern, handlers)
func (r *Router) Get(pattern string, h ...Handler) {
	r.Handle("GET", pattern, h)
}

// Patch is a shortcut for r.Handle("PATCH", pattern, handlers)
func (r *Router) Patch(pattern string, h ...Handler) {
	r.Handle("PATCH", pattern, h)
}

// Post is a shortcut for r.Handle("POST", pattern, handlers)
func (r *Router) Post(pattern string, h ...Handler) {
	r.Handle("POST", pattern, h)
}

// Put is a shortcut for r.Handle("PUT", pattern, handlers)
func (r *Router) Put(pattern string, h ...Handler) {
	r.Handle("PUT", pattern, h)
}

// Delete is a shortcut for r.Handle("DELETE", pattern, handlers)
func (r *Router) Delete(pattern string, h ...Handler) {
	r.Handle("DELETE", pattern, h)
}

// Options is a shortcut for r.Handle("OPTIONS", pattern, handlers)
func (r *Router) Options(pattern string, h ...Handler) {
	r.Handle("OPTIONS", pattern, h)
}

// Head is a shortcut for r.Handle("HEAD", pattern, handlers)
func (r *Router) Head(pattern string, h ...Handler) {
	r.Handle("HEAD", pattern, h)
}

// Configurable http.HandlerFunc which is called when no matching route is
// found. If it is not set, http.NotFound is used.
// Be sure to set 404 response code in your handler.
func (r *Router) NotFound(handlers ...Handler) {
	r.router.NotFound = func(resp http.ResponseWriter, req *http.Request) {
		c := r.m.createContext(resp, req)
		c.handlers = append(r.m.handlers, handlers...)
		c.run()
	}
}

// \_   _____/ _______  __
//  |    __)_ /    \  \/ /
//  |        \   |  \   /
// /_______  /___|  /\_/
//         \/     \/

const (
	DEV  string = "development"
	PROD string = "production"
	TEST string = "test"
)

// Env is the environment that Macaron is executing in.
// The MACARON_ENV is read on initialization to set this variable.
var Env = DEV
var Root string

func setENV(e string) {
	if len(e) > 0 {
		Env = e
	}
}

func init() {
	setENV(os.Getenv("MACARON_ENV"))
	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}
	Root = filepath.Dir(path)
}
