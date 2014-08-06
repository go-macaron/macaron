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
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

var currentRoot, _ = os.Getwd()

func Test_Static(t *testing.T) {
	response := httptest.NewRecorder()
	response.Body = new(bytes.Buffer)

	m := New()
	r := NewRouter()

	m.Use(Static(currentRoot))
	m.Action(r.Handle)

	req, err := http.NewRequest("GET", "http://localhost:4000/macaron.go", nil)
	if err != nil {
		t.Error(err)
	}
	m.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
	expect(t, response.Header().Get("Expires"), "")
	if response.Body.Len() == 0 {
		t.Errorf("Got empty body for GET request")
	}
}

func Test_Static_Local_Path(t *testing.T) {
	Root = os.TempDir()
	response := httptest.NewRecorder()
	response.Body = new(bytes.Buffer)

	m := New()
	r := NewRouter()

	m.Use(Static("."))
	f, err := ioutil.TempFile(Root, "static_content")
	if err != nil {
		t.Error(err)
	}
	f.WriteString("Expected Content")
	f.Close()
	m.Action(r.Handle)

	req, err := http.NewRequest("GET", "http://localhost:4000/"+path.Base(f.Name()), nil)
	if err != nil {
		t.Error(err)
	}
	m.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
	expect(t, response.Header().Get("Expires"), "")
	expect(t, response.Body.String(), "Expected Content")
}

func Test_Static_Head(t *testing.T) {
	response := httptest.NewRecorder()
	response.Body = new(bytes.Buffer)

	m := New()
	r := NewRouter()

	m.Use(Static(currentRoot))
	m.Action(r.Handle)

	req, err := http.NewRequest("HEAD", "http://localhost:4000/macaron.go", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
	if response.Body.Len() != 0 {
		t.Errorf("Got non-empty body for HEAD request")
	}
}

func Test_Static_As_Post(t *testing.T) {
	response := httptest.NewRecorder()

	m := New()
	r := NewRouter()

	m.Use(Static(currentRoot))
	m.Action(r.Handle)

	req, err := http.NewRequest("POST", "http://localhost:4000/macaron.go", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusNotFound)
}

func Test_Static_BadDir(t *testing.T) {
	response := httptest.NewRecorder()

	m := Classic()

	req, err := http.NewRequest("GET", "http://localhost:4000/macaron.go", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	refute(t, response.Code, http.StatusOK)
}

func Test_Static_Options_Logging(t *testing.T) {
	response := httptest.NewRecorder()

	var buffer bytes.Buffer
	m := NewWithLogger(&buffer)
	opt := StaticOptions{}
	m.Use(Static(currentRoot, opt))

	req, err := http.NewRequest("GET", "http://localhost:4000/macaron.go", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
	expect(t, buffer.String(), "[Macaron] [Static] Serving /macaron.go\n")

	// Now without logging
	m.Handlers()
	buffer.Reset()

	// This should disable logging
	opt.SkipLogging = true
	m.Use(Static(currentRoot, opt))

	m.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
	expect(t, buffer.String(), "")
}

func Test_Static_Options_ServeIndex(t *testing.T) {
	response := httptest.NewRecorder()

	var buffer bytes.Buffer
	m := NewWithLogger(&buffer)
	opt := StaticOptions{IndexFile: "macaron.go"} // Define macaron.go as index file
	m.Use(Static(currentRoot, opt))

	req, err := http.NewRequest("GET", "http://localhost:4000/", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
	expect(t, buffer.String(), "[Macaron] [Static] Serving /macaron.go\n")
}

func Test_Static_Options_Prefix(t *testing.T) {
	response := httptest.NewRecorder()

	var buffer bytes.Buffer
	m := NewWithLogger(&buffer)

	// Serve current directory under /public
	m.Use(Static(currentRoot, StaticOptions{Prefix: "/public"}))

	// Check file content behaviour
	req, err := http.NewRequest("GET", "http://localhost:4000/public/macaron.go", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
	expect(t, buffer.String(), "[Macaron] [Static] Serving /macaron.go\n")
}

func Test_Static_Options_Expires(t *testing.T) {
	response := httptest.NewRecorder()

	var buffer bytes.Buffer
	m := NewWithLogger(&buffer)

	// Serve current directory under /public
	m.Use(Static(currentRoot, StaticOptions{Expires: func() string { return "46" }}))

	// Check file content behaviour
	req, err := http.NewRequest("GET", "http://localhost:4000/macaron.go", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	expect(t, response.Header().Get("Expires"), "46")
}

func Test_Static_Redirect(t *testing.T) {
	response := httptest.NewRecorder()

	m := New()
	m.Use(Static(currentRoot, StaticOptions{Prefix: "/public"}))

	req, err := http.NewRequest("GET", "http://localhost:4000/public", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusFound)
	expect(t, response.Header().Get("Location"), "/public/")
}
