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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Context(t *testing.T) {
	Convey("Do advanced encapsulation operations", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory: "fixtures/basic",
		}))

		Convey("Render HTML", func() {
			m.Get("/html", func(ctx *Context) {
				ctx.HTML(304, "hello", "Unknwon") // 304 for logger test.
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/html", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "<h1>Hello Unknwon</h1>\n")

			m.Get("/html2", func(ctx *Context) {
				ctx.Data["Name"] = "Unknwon"
				ctx.HTML(200, "hello2")
			})

			resp = httptest.NewRecorder()
			req, err = http.NewRequest("GET", "/html2", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "<h1>Hello Unknwon</h1>\n")
		})

		Convey("Parse from and query", func() {
			m.Get("/query", func(ctx *Context) string {
				return ctx.Query("name")
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/query?name=Unknwon", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "Unknwon")
		})

		Convey("URL parameter", func() {
			m.Get("/:name", func(ctx *Context) string {
				return ctx.Params(":name")
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/user", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "user")
		})

		Convey("Set and get cookie", func() {
			m.Get("/set", func(ctx *Context) {
				ctx.SetCookie("user", "Unknwon", 1)
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/set", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Header().Get("Set-Cookie"), ShouldEqual, "user=Unknwon; Path=/; Max-Age=1")

			m.Get("/get", func(ctx *Context) string {
				ctx.GetCookie("404")
				return ctx.GetCookie("user")
			})

			resp = httptest.NewRecorder()
			req, err = http.NewRequest("GET", "/get", nil)
			So(err, ShouldBeNil)
			req.Header.Set("Cookie", "user=Unknwon; Path=/; Max-Age=1")
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "Unknwon")
		})

		Convey("Set and get secure cookie", func() {
			m.SetDefaultCookieSecret("macaron")
			m.Get("/set", func(ctx *Context) {
				ctx.SetSecureCookie("user", "Unknwon", 1)
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/set", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(strings.HasPrefix(resp.Header().Get("Set-Cookie"), "user=VW5rbndvbg==|"), ShouldBeTrue)

			m.Get("/get", func(ctx *Context) string {
				name, ok := ctx.GetSecureCookie("user")
				So(ok, ShouldBeTrue)
				return name
			})

			resp = httptest.NewRecorder()
			req, err = http.NewRequest("GET", "/get", nil)
			So(err, ShouldBeNil)
			req.Header.Set("Cookie", "user=VW5rbndvbg==|1409244667158399419|6097781707f68d9940ba1ef0e78cc84aaeebc48f; Path=/; Max-Age=1")
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "Unknwon")
		})
	})
}
