// Copyright 2014 Unknwon
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

package macaron

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Unknwon/macaron/inject"
)

// Context represents the runtime context of current request of Macaron instance.
// It is the integration of most frequently used middlewares and helper methods.
type Context struct {
	inject.Injector
	handlers []Handler
	action   Handler
	index    int

	*Router
	Req    *http.Request
	Resp   ResponseWriter
	params Params
	Render // Not nil only if you use macaran.Render middleware.
	Data   map[string]interface{}
}

func (c *Context) handler() Handler {
	if c.index < len(c.handlers) {
		return c.handlers[c.index]
	}
	if c.index == len(c.handlers) {
		return c.action
	}
	panic("invalid index for context handler")
}

func (c *Context) Next() {
	c.index += 1
	c.run()
}

func (c *Context) Written() bool {
	return c.Resp.Written()
}

func (c *Context) run() {
	for c.index <= len(c.handlers) {
		vals, err := c.Invoke(c.handler())
		if err != nil {
			panic(err)
		}
		c.index += 1

		// if the handler returned something, write it to the http response
		if len(vals) > 0 {
			ev := c.GetVal(reflect.TypeOf(ReturnHandler(nil)))
			handleReturn := ev.Interface().(ReturnHandler)
			handleReturn(c, vals)
		}

		if c.Written() {
			return
		}
	}
}

// RemoteAddr returns more real IP address.
func (ctx *Context) RemoteAddr() string {
	addr := ctx.Req.Header.Get("X-Real-IP")
	if len(addr) == 0 {
		addr = ctx.Req.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = ctx.Req.RemoteAddr
			if i := strings.LastIndex(addr, ":"); i > -1 {
				addr = addr[:i]
			}
		}
	}
	return addr
}

// HTML calls Render.HTML but allows less arguments.
func (ctx *Context) HTML(status int, name string, binding ...interface{}) {
	if len(binding) == 0 {
		ctx.Render.HTML(status, name, ctx.Data)
	} else {
		ctx.Render.HTML(status, name, binding[0])
		if len(binding) > 1 {
			ctx.Render.HTML(status, name, binding[1].(HTMLOptions))
		}
	}
}

// Query querys form parameter.
func (ctx *Context) Query(name string) string {
	ctx.Req.ParseForm()
	return ctx.Req.Form.Get(name)
}

// Params return value of given param name.
func (ctx *Context) Params(name string) string {
	return ctx.params[name]
}

// SetCookie sets given cookie value to response header.
func (ctx *Context) SetCookie(name string, value string, others ...interface{}) {
	cookie := http.Cookie{}
	cookie.Name = name
	cookie.Value = value

	if len(others) > 0 {
		switch v := others[0].(type) {
		case int:
			cookie.MaxAge = v
		case int64:
			cookie.MaxAge = int(v)
		case int32:
			cookie.MaxAge = int(v)
		}
	}

	// default "/"
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			cookie.Path = v
		}
	} else {
		cookie.Path = "/"
	}

	// default empty
	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			cookie.Domain = v
		}
	}

	// default empty
	if len(others) > 3 {
		switch v := others[3].(type) {
		case bool:
			cookie.Secure = v
		default:
			if others[3] != nil {
				cookie.Secure = true
			}
		}
	}

	// default false. for session cookie default true
	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			cookie.HttpOnly = true
		}
	}

	ctx.Resp.Header().Add("Set-Cookie", cookie.String())
}

// GetCookie returns given cookie value from request header.
func (ctx *Context) GetCookie(name string) string {
	cookie, err := ctx.Req.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

var defaultCookieSecret string

// SetDefaultCookieSecret sets global default secure cookie secret.
func (m *Macaron) SetDefaultCookieSecret(secret string) {
	defaultCookieSecret = secret
}

// SetSecureCookie sets given cookie value to response header with default secret string.
func (ctx *Context) SetSecureCookie(name, value string, others ...interface{}) {
	ctx.SetSuperSecureCookie(defaultCookieSecret, name, value, others...)
}

// GetSecureCookie returns given cookie value from request header with default secret string.
func (ctx *Context) GetSecureCookie(key string) (string, bool) {
	return ctx.GetSuperSecureCookie(defaultCookieSecret, key)
}

// SetSuperSecureCookie sets given cookie value to response header with secret string.
func (ctx *Context) SetSuperSecureCookie(Secret, name, value string, others ...interface{}) {
	vs := base64.URLEncoding.EncodeToString([]byte(value))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", h.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")
	ctx.SetCookie(name, cookie, others...)
}

// GetSuperSecureCookie returns given cookie value from request header with secret string.
func (ctx *Context) GetSuperSecureCookie(Secret, key string) (string, bool) {
	val := ctx.GetCookie(key)
	if val == "" {
		return "", false
	}

	parts := strings.SplitN(val, "|", 3)

	if len(parts) != 3 {
		return "", false
	}

	vs := parts[0]
	timestamp := parts[1]
	sig := parts[2]

	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)

	if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res), true
}

// ServeFile serves given file to response.
func (ctx *Context) ServeFile(file string, names ...string) {
	var name string
	if len(names) > 0 {
		name = names[0]
	} else {
		name = path.Base(file)
	}
	ctx.Resp.Header().Set("Content-Description", "File Transfer")
	ctx.Resp.Header().Set("Content-Type", "application/octet-stream")
	ctx.Resp.Header().Set("Content-Disposition", "attachment; filename="+name)
	ctx.Resp.Header().Set("Content-Transfer-Encoding", "binary")
	ctx.Resp.Header().Set("Expires", "0")
	ctx.Resp.Header().Set("Cache-Control", "must-revalidate")
	ctx.Resp.Header().Set("Pragma", "public")
	http.ServeFile(ctx.Resp, ctx.Req, file)
}
