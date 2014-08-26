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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_GzipAll(t *testing.T) {
	// Set up
	recorder := httptest.NewRecorder()
	before := false

	m := New()
	m.Use(Gzip())
	m.Use(func(r http.ResponseWriter) {
		r.(ResponseWriter).Before(func(rw ResponseWriter) {
			before = true
		})
	})
	m.Get("/", func() {})

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(recorder, r)

	// Make our assertions
	_, ok := recorder.HeaderMap[HeaderContentEncoding]
	if ok {
		t.Error(HeaderContentEncoding + " present")
	}

	ce := recorder.Header().Get(HeaderContentEncoding)
	if strings.EqualFold(ce, "gzip") {
		t.Error(HeaderContentEncoding + " is 'gzip'")
	}

	recorder = httptest.NewRecorder()
	r.Header.Set(HeaderAcceptEncoding, "gzip")
	m.ServeHTTP(recorder, r)

	// Make our assertions
	_, ok = recorder.HeaderMap[HeaderContentEncoding]
	if !ok {
		t.Error(HeaderContentEncoding + " not present")
	}

	ce = recorder.Header().Get(HeaderContentEncoding)
	if !strings.EqualFold(ce, "gzip") {
		t.Error(HeaderContentEncoding + " is not 'gzip'")
	}

	if before == false {
		t.Error("Before hook was not called")
	}
}
