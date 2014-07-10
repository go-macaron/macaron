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
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Recovery(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()

	setENV(DEV)
	m := New()
	// replace log for testing
	m.Map(log.New(buff, "[Macaron] ", 0))
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "unpredictable")
	})
	m.Use(Recovery())
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		panic("here is a panic!")
	})
	m.Get("/", func() {})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(recorder, req)

	expect(t, recorder.Code, http.StatusInternalServerError)
	expect(t, recorder.HeaderMap.Get("Content-Type"), "text/html")
	refute(t, recorder.Body.Len(), 0)
	refute(t, len(buff.String()), 0)
}

func Test_Recovery_ResponseWriter(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder2 := httptest.NewRecorder()

	setENV(DEV)
	m := New()
	m.Use(Recovery())
	m.Use(func(c *Context) {
		c.MapTo(recorder2, (*http.ResponseWriter)(nil))
		panic("here is a panic!")
	})
	m.Get("/", func() {})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(recorder, req)

	expect(t, recorder2.Code, http.StatusInternalServerError)
	expect(t, recorder2.HeaderMap.Get("Content-Type"), "text/html")
	refute(t, recorder2.Body.Len(), 0)
}
