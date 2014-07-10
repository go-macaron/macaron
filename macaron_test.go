// Copyright 2013 Martini Authors
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

package macaron

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_New(t *testing.T) {
	m := New()
	if m == nil {
		t.Error("martini.New() cannot return nil")
	}
}

func Test_Macaron_Run(t *testing.T) {
	// just test that Run doesn't bomb
	go New().Run()
}

func Test_Macaron_ServeHTTP(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	m := New()
	m.Use(func(c *Context) {
		result += "foo"
		c.Next()
		result += "ban"
	})
	m.Use(func(c *Context) {
		result += "bar"
		c.Next()
		result += "baz"
	})
	m.Get("/", func() {})
	m.Action(func(res http.ResponseWriter, req *http.Request) {
		result += "bat"
		res.WriteHeader(http.StatusBadRequest)
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)

	expect(t, result, "foobarbatbazban")
	expect(t, response.Code, http.StatusBadRequest)
}

func Test_Macaron_Handlers(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	batman := func(c *Context) {
		result += "batman!"
	}

	m := New()
	m.Use(func(c *Context) {
		result += "foo"
		c.Next()
		result += "ban"
	})
	m.Handlers(
		batman,
		batman,
		batman,
	)
	m.Get("/", func() {})
	m.Action(func(res http.ResponseWriter, req *http.Request) {
		result += "bat"
		res.WriteHeader(http.StatusBadRequest)
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)

	expect(t, result, "batman!batman!batman!bat")
	expect(t, response.Code, http.StatusBadRequest)
}

func Test_Macaron_EarlyWrite(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	m := New()
	m.Use(func(res http.ResponseWriter) {
		result += "foobar"
		res.Write([]byte("Hello world"))
	})
	m.Use(func() {
		result += "bat"
	})
	m.Get("/", func() {})
	m.Action(func(res http.ResponseWriter) {
		result += "baz"
		res.WriteHeader(http.StatusBadRequest)
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)

	expect(t, result, "foobar")
	expect(t, response.Code, http.StatusOK)
}

func Test_Macaron_Written(t *testing.T) {
	response := httptest.NewRecorder()

	m := New()
	m.Handlers(func(res http.ResponseWriter) {
		res.WriteHeader(http.StatusOK)
	})

	ctx := m.createContext(response, &http.Request{Method: "GET"})
	expect(t, ctx.Written(), false)

	ctx.run()
	expect(t, ctx.Written(), true)
}

func Test_Macaron_Basic_NoRace(t *testing.T) {
	m := New()
	handlers := []Handler{func() {}, func() {}}
	// Ensure append will not realloc to trigger the race condition
	m.handlers = handlers[:1]
	req, _ := http.NewRequest("GET", "/", nil)
	for i := 0; i < 2; i++ {
		go func() {
			response := httptest.NewRecorder()
			m.ServeHTTP(response, req)
		}()
	}
}

func Test_SetENV(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", "development"},
		{"not_development", "not_development"},
	}

	for _, test := range tests {
		setENV(test.in)
		if Env != test.out {
			expect(t, Env, test.out)
		}
	}
}
