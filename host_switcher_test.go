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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_HostSwitcher(t *testing.T) {
	Convey("Hosting multiple instances", t, func() {
		hs := NewHostSwitcher()

		m1 := Classic()
		m1.Get("/", func() string {
			return "welcome to gowalker.org"
		})
		hs.Set("gowalker.org", m1)

		m2 := Classic()
		m2.Get("/", func() string {
			return "welcome to gogs.io"
		})
		hs.Set("gogs.io", m2)

		Convey("Remove a instance", func() {
			hs.Remove("gopm.io")
		})

		Convey("Request instance 1", func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			req.Host = "gowalker.org"
			hs.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "welcome to gowalker.org")
		})

		Convey("Request instance 2", func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			req.Host = "gogs.io"
			hs.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "welcome to gogs.io")
		})

		Convey("Request a instance that not exist", func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			req.Host = "gopm.io"
			hs.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, 404)
		})

		Convey("Just test that Run() doesn't bomb", func() {
			go hs.RunOnAddr(":4002")
		})
	})
}
