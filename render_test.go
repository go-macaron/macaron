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

func Test_Render_Bad_HTML(t *testing.T) {
	Convey("Render HTML", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory: "fixtures/basic",
		}))
		m.Get("/foobar", func(r Render) {
			r.HTML(200, "hello", "jeremy")
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentHTML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, "<h1>Hello jeremy</h1>\n")
	})

	Convey("Render HTML and return string", t, func() {
		m := Classic()
		m.Use(Renderer(RenderOptions{
			Directory: "fixtures/basic",
		}))
		m.Get("/foobar", func(r Render) {
			result, err := r.HTMLString("hello", "jeremy")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "<h1>Hello jeremy</h1>\n")
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
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

// TODO
func Test_Render_Extensions(t *testing.T) {
	m := Classic()
	m.Use(Renderer(RenderOptions{
		Directory:  "fixtures/basic",
		Extensions: []string{".tmpl", ".html"},
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "hypertext", nil)
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "Hypertext!\n")
}

func Test_Render_Funcs(t *testing.T) {

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

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "index", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Body.String(), "My custom function\n")
}

func Test_Render_Layout(t *testing.T) {
	m := Classic()
	m.Use(Renderer(RenderOptions{
		Directory: "fixtures/basic",
		Layout:    "layout",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "content", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Body.String(), "head\n<h1>jeremy</h1>\n\nfoot\n")
}

func Test_Render_Layout_Current(t *testing.T) {
	m := Classic()
	m.Use(Renderer(RenderOptions{
		Directory: "fixtures/basic",
		Layout:    "current_layout",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "content", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Body.String(), "content head\n<h1>jeremy</h1>\n\ncontent foot\n")
}

func Test_Render_Nested_HTML(t *testing.T) {
	m := Classic()
	m.Use(Renderer(RenderOptions{
		Directory: "fixtures/basic",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "admin/index", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "<h1>Admin jeremy</h1>\n")
}

func Test_Render_Delimiters(t *testing.T) {
	m := Classic()
	m.Use(Renderer(RenderOptions{
		Delims:    Delims{"{[{", "}]}"},
		Directory: "fixtures/basic",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "delims", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "<h1>Hello jeremy</h1>")
}

func Test_Render_BinaryData(t *testing.T) {
	m := Classic()
	m.Use(Renderer(RenderOptions{
	// nothing here to configure
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.RawData(200, []byte("hello there"))
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentBinary)
	expect(t, res.Body.String(), "hello there")
}

func Test_Render_BinaryData_CustomMimeType(t *testing.T) {
	m := Classic()
	m.Use(Renderer(RenderOptions{
	// nothing here to configure
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.Header().Set(ContentType, "image/jpeg")
		r.RawData(200, []byte("..jpeg data.."))
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), "image/jpeg")
	expect(t, res.Body.String(), "..jpeg data..")
}

func Test_Render_Status204(t *testing.T) {
	res := httptest.NewRecorder()
	r := TplRender{res, nil, nil, RenderOptions{}, "", time.Now()}
	r.Status(204)
	expect(t, res.Code, 204)
}

func Test_Render_Error404(t *testing.T) {
	res := httptest.NewRecorder()
	r := TplRender{res, nil, nil, RenderOptions{}, "", time.Now()}
	r.Error(404)
	expect(t, res.Code, 404)
}

func Test_Render_Error500(t *testing.T) {
	res := httptest.NewRecorder()
	r := TplRender{res, nil, nil, RenderOptions{}, "", time.Now()}
	r.Error(500)
	expect(t, res.Code, 500)
}

func Test_Render_Redirect_Default(t *testing.T) {
	url, _ := url.Parse("http://localhost/path/one")
	req := http.Request{
		Method: "GET",
		URL:    url,
	}
	res := httptest.NewRecorder()

	r := TplRender{res, &req, nil, RenderOptions{}, "", time.Now()}
	r.Redirect("two")

	expect(t, res.Code, 302)
	expect(t, res.HeaderMap["Location"][0], "/path/two")
}

func Test_Render_Redirect_Code(t *testing.T) {
	url, _ := url.Parse("http://localhost/path/one")
	req := http.Request{
		Method: "GET",
		URL:    url,
	}
	res := httptest.NewRecorder()

	r := TplRender{res, &req, nil, RenderOptions{}, "", time.Now()}
	r.Redirect("two", 307)

	expect(t, res.Code, 307)
	expect(t, res.HeaderMap["Location"][0], "/path/two")
}

func Test_Render_Charset_JSON(t *testing.T) {
	m := Classic()
	m.Use(Renderer(RenderOptions{
		Charset: "foobar",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.JSON(300, Greeting{"hello", "world"})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 300)
	expect(t, res.Header().Get(ContentType), ContentJSON+"; charset=foobar")
	expect(t, res.Body.String(), `{"one":"hello","two":"world"}`)
}

func Test_Render_Override_Layout(t *testing.T) {
	m := Classic()
	m.Use(Renderer(RenderOptions{
		Directory: "fixtures/basic",
		Layout:    "layout",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "content", "jeremy", HTMLOptions{
			Layout: "another_layout",
		})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "another head\n<h1>jeremy</h1>\n\nanother foot\n")
}

func Test_Render_NoRace(t *testing.T) {
	// This test used to fail if run with -race
	m := Classic()
	m.Use(Renderer(RenderOptions{
		Directory: "fixtures/basic",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "hello", "world")
	})

	done := make(chan bool)
	doreq := func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/foobar", nil)

		m.ServeHTTP(res, req)

		expect(t, res.Code, 200)
		expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
		// ContentLength should be deferred to the ResponseWriter and not Render
		expect(t, res.Header().Get(ContentLength), "")
		expect(t, res.Body.String(), "<h1>Hello world</h1>\n")
		done <- true
	}
	// Run two requests to check there is no race condition
	go doreq()
	go doreq()
	<-done
	<-done
}

func Test_GetExt(t *testing.T) {
	expect(t, getExt("test"), "")
	expect(t, getExt("test.tmpl"), ".tmpl")
	expect(t, getExt("test.go.html"), ".go.html")
}
