// Copyright 2013 Martini Authors
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
	"encoding/xml"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type Greeting struct {
	One string `json:"one"`
	Two string `json:"two"`
}

type GreetingXML struct {
	XMLName xml.Name `xml:"greeting"`
	One     string   `xml:"one,attr"`
	Two     string   `xml:"two,attr"`
}

func Test_Render_JSON(t *testing.T) {
	Convey("Render JSON", t, func() {
		m := Classic()
		m.Use(Renderer())
		m.Get("/foobar", func(r Render) {
			r.JSON(300, Greeting{"hello", "world"})
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusMultipleChoices)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentJSON+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, `{"one":"hello","two":"world"}`)
	})

	Convey("Render JSON with prefix", t, func() {
		m := Classic()
		prefix := ")]}',\n"
		m.Use(Renderer(RenderOptions{
			PrefixJSON: []byte(prefix),
		}))
		m.Get("/foobar", func(r Render) {
			r.JSON(300, Greeting{"hello", "world"})
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusMultipleChoices)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentJSON+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, prefix+`{"one":"hello","two":"world"}`)
	})

	Convey("Render Indented JSON", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			IndentJSON: true,
		}))
		m.Get("/foobar", func(r Render) {
			r.JSON(300, Greeting{"hello", "world"})
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusMultipleChoices)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentJSON+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, `{
  "one": "hello",
  "two": "world"
}`)
	})

	Convey("Render JSON and return string", t, func() {
		m := Classic()
		m.Use(Renderer())
		m.Get("/foobar", func(r Render) {
			result, err := r.JSONString(Greeting{"hello", "world"})
			So(err, ShouldBeNil)
			So(result, ShouldEqual, `{"one":"hello","two":"world"}`)
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
	})

	Convey("Render with charset JSON", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Charset: "foobar",
		}))
		m.Get("/foobar", func(r Render) {
			r.JSON(300, Greeting{"hello", "world"})
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusMultipleChoices)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentJSON+"; charset=foobar")
		So(resp.Body.String(), ShouldEqual, `{"one":"hello","two":"world"}`)
	})
}

func Test_Render_XML(t *testing.T) {
	Convey("Render XML", t, func() {
		m := Classic()
		m.Use(Renderer())
		m.Get("/foobar", func(r Render) {
			r.XML(300, GreetingXML{One: "hello", Two: "world"})
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusMultipleChoices)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentXML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, `<greeting one="hello" two="world"></greeting>`)
	})

	Convey("Render XML with prefix", t, func() {
		m := Classic()
		prefix := ")]}',\n"
		m.Use(Renderer(RenderOptions{
			PrefixXML: []byte(prefix),
		}))
		m.Get("/foobar", func(r Render) {
			r.XML(300, GreetingXML{One: "hello", Two: "world"})
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusMultipleChoices)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentXML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, prefix+`<greeting one="hello" two="world"></greeting>`)
	})

	Convey("Render Indented XML", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			IndentXML: true,
		}))
		m.Get("/foobar", func(r Render) {
			r.XML(300, GreetingXML{One: "hello", Two: "world"})
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusMultipleChoices)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentXML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, `<greeting one="hello" two="world"></greeting>`)
	})
}

func Test_Render_HTML(t *testing.T) {
	Convey("Render HTML", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory: "fixtures/basic",
		}))
		m.Use(Renderer(RenderOptions{
			Name:      "basic2",
			Directory: "fixtures/basic2",
		}))
		m.Get("/foobar", func(r Render) {
			r.HTML(200, "hello", "jeremy")
			r.SetTemplatePath("", "fixtures/basic2")
		})
		m.Get("/foobar2", func(r Render) {
			r.HTMLSet(200, "basic2", "hello", "jeremy")
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentHTML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, "<h1>Hello jeremy</h1>\n")

		resp = httptest.NewRecorder()
		req, err = http.NewRequest("GET", "/foobar2", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentHTML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, "<h1>What's up, jeremy</h1>\n")

		Convey("Change render templates path", func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/foobar", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)

			So(resp.Code, ShouldEqual, http.StatusOK)
			So(resp.Header().Get(ContentType), ShouldEqual, ContentHTML+"; charset=UTF-8")
			So(resp.Body.String(), ShouldEqual, "<h1>What's up, jeremy</h1>\n")
		})
	})

	Convey("Render HTML and return string", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory: "fixtures/basic",
		}))
		m.Use(Renderer(RenderOptions{
			Name:      "basic2",
			Directory: "fixtures/basic2",
		}))
		m.Get("/foobar", func(r Render) {
			result, err := r.HTMLString("hello", "jeremy")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "<h1>Hello jeremy</h1>\n")
		})
		m.Get("/foobar2", func(r Render) {
			result, err := r.HTMLSetString("basic2", "hello", "jeremy")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "<h1>What's up, jeremy</h1>\n")
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		resp = httptest.NewRecorder()
		req, err = http.NewRequest("GET", "/foobar2", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
	})

	Convey("Render with nested HTML", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory: "fixtures/basic",
		}))
		m.Get("/foobar", func(r Render) {
			r.HTML(200, "admin/index", "jeremy")
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentHTML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, "<h1>Admin jeremy</h1>\n")
	})

	Convey("Render bad HTML", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory: "fixtures/basic",
		}))
		m.Get("/foobar", func(r Render) {
			r.HTML(200, "nope", nil)
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		So(resp.Body.String(), ShouldEqual, "html/template: \"nope\" is undefined\n")
	})
}

func Test_Render_XHTML(t *testing.T) {
	Convey("Render XHTML", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory:       "fixtures/basic",
			HTMLContentType: ContentXHTML,
		}))
		m.Get("/foobar", func(r Render) {
			r.HTML(200, "hello", "jeremy")
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentXHTML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, "<h1>Hello jeremy</h1>\n")
	})
}

func Test_Render_Extensions(t *testing.T) {
	Convey("Render with extensions", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory:  "fixtures/basic",
			Extensions: []string{".tmpl", ".html"},
		}))
		m.Get("/foobar", func(r Render) {
			r.HTML(200, "hypertext", nil)
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentHTML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, "Hypertext!\n")
	})
}

func Test_Render_Funcs(t *testing.T) {
	Convey("Render with functions", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory: "fixtures/custom_funcs",
			Funcs: []template.FuncMap{
				{
					"myCustomFunc": func() string {
						return "My custom function"
					},
				},
			},
		}))
		m.Get("/foobar", func(r Render) {
			r.HTML(200, "index", "jeremy")
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Body.String(), ShouldEqual, "My custom function\n")
	})
}

func Test_Render_Layout(t *testing.T) {
	Convey("Render with layout", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory: "fixtures/basic",
			Layout:    "layout",
		}))
		m.Get("/foobar", func(r Render) {
			r.HTML(200, "content", "jeremy")
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Body.String(), ShouldEqual, "head\n<h1>jeremy</h1>\n\nfoot\n")
	})

	Convey("Render with current layout", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory: "fixtures/basic",
			Layout:    "current_layout",
		}))
		m.Get("/foobar", func(r Render) {
			r.HTML(200, "content", "jeremy")
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Body.String(), ShouldEqual, "content head\n<h1>jeremy</h1>\n\ncontent foot\n")
	})

	Convey("Render with override layout", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory: "fixtures/basic",
			Layout:    "layout",
		}))
		m.Get("/foobar", func(r Render) {
			r.HTML(200, "content", "jeremy", HTMLOptions{
				Layout: "another_layout",
			})
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentHTML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, "another head\n<h1>jeremy</h1>\n\nanother foot\n")
	})
}

func Test_Render_Delimiters(t *testing.T) {
	Convey("Render with delimiters", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Delims:    Delims{"{[{", "}]}"},
			Directory: "fixtures/basic",
		}))
		m.Get("/foobar", func(r Render) {
			r.HTML(200, "delims", "jeremy")
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentHTML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, "<h1>Hello jeremy</h1>")
	})
}

func Test_Render_BinaryData(t *testing.T) {
	Convey("Render binary data", t, func() {
		m := Classic()
		m.Use(Renderer())
		m.Get("/foobar", func(r Render) {
			r.RawData(200, []byte("hello there"))
		})
		m.Get("/foobar2", func(r Render) {
			r.RenderData(200, []byte("hello there"))
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentBinary)
		So(resp.Body.String(), ShouldEqual, "hello there")

		resp = httptest.NewRecorder()
		req, err = http.NewRequest("GET", "/foobar2", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentHTML)
		So(resp.Body.String(), ShouldEqual, "hello there")
	})

	Convey("Render binary data with mime type", t, func() {
		m := Classic()
		m.Use(Renderer())
		m.Get("/foobar", func(r Render) {
			r.RW().Header().Set(ContentType, "image/jpeg")
			r.RawData(200, []byte("..jpeg data.."))
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, "image/jpeg")
		So(resp.Body.String(), ShouldEqual, "..jpeg data..")
	})
}

func Test_Render_Status(t *testing.T) {
	Convey("Render with status 204", t, func() {
		resp := httptest.NewRecorder()
		r := TplRender{resp, nil, &RenderOptions{}, "", time.Now()}
		r.Status(204)
		So(resp.Code, ShouldEqual, http.StatusNoContent)
	})

	Convey("Render with status 404", t, func() {
		resp := httptest.NewRecorder()
		r := TplRender{resp, nil, &RenderOptions{}, "", time.Now()}
		r.Error(404)
		So(resp.Code, ShouldEqual, http.StatusNotFound)
	})

	Convey("Render with status 500", t, func() {
		resp := httptest.NewRecorder()
		r := TplRender{resp, nil, &RenderOptions{}, "", time.Now()}
		r.Error(500)
		So(resp.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

func Test_Render_Redirect_Default(t *testing.T) {
	Convey("Render with default redirect", t, func() {
		url, err := url.Parse("http://localhost/path/one")
		So(err, ShouldBeNil)
		resp := httptest.NewRecorder()
		req := http.Request{
			Method: "GET",
			URL:    url,
		}
		r := TplRender{resp, &req, &RenderOptions{}, "", time.Now()}
		r.Redirect("two")

		So(resp.Code, ShouldEqual, http.StatusFound)
		So(resp.HeaderMap["Location"][0], ShouldEqual, "/path/two")
	})

	Convey("Render with custom redirect", t, func() {
		url, err := url.Parse("http://localhost/path/one")
		So(err, ShouldBeNil)
		resp := httptest.NewRecorder()
		req := http.Request{
			Method: "GET",
			URL:    url,
		}
		r := TplRender{resp, &req, &RenderOptions{}, "", time.Now()}
		r.Redirect("two", 307)

		So(resp.Code, ShouldEqual, http.StatusTemporaryRedirect)
		So(resp.HeaderMap["Location"][0], ShouldEqual, "/path/two")
	})
}

func Test_Render_NoRace(t *testing.T) {
	Convey("Make sure render has no race", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory: "fixtures/basic",
		}))
		m.Get("/foobar", func(r Render) {
			r.HTML(200, "hello", "world")
		})

		done := make(chan bool)
		doreq := func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/foobar", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)

			So(resp.Code, ShouldEqual, http.StatusOK)
			So(resp.Header().Get(ContentType), ShouldEqual, ContentHTML+"; charset=UTF-8")
			// ContentLength should be deferred to the ResponseWriter and not Render
			So(resp.Header().Get(ContentLength), ShouldBeBlank)
			So(resp.Body.String(), ShouldEqual, "<h1>Hello world</h1>\n")
			done <- true
		}
		// Run two requests to check there is no race condition
		go doreq()
		go doreq()
		<-done
		<-done
	})
}

func Test_GetExt(t *testing.T) {
	Convey("Get extension", t, func() {
		So(getExt("test"), ShouldBeBlank)
		So(getExt("test.tmpl"), ShouldEqual, ".tmpl")
		So(getExt("test.go.tmpl"), ShouldEqual, ".go.tmpl")
	})
}
