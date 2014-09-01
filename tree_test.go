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
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_splitSegment(t *testing.T) {
	type result struct {
		Ok    bool
		Parts []string
		Regex string
	}
	cases := map[string]result{
		"admin":              result{false, nil, ""},
		":id":                result{true, []string{":id"}, ""},
		"?:id":               result{true, []string{":", ":id"}, ""},
		":id:int":            result{true, []string{":id"}, "([0-9]+)"},
		":name:string":       result{true, []string{":name"}, `([\w]+)`},
		":id([0-9]+)":        result{true, []string{":id"}, "([0-9]+)"},
		":id([0-9]+)_:name":  result{true, []string{":id", ":name"}, "([0-9]+)_(.+)"},
		"cms_:id_:page.html": result{true, []string{":id", ":page"}, "cms_(.+)_(.+).html"},
		"*":                  result{true, []string{":splat"}, ""},
		"*.*":                result{true, []string{".", ":path", ":ext"}, ""},
	}
	Convey("Splits segment into parts", t, func() {
		for key, result := range cases {
			ok, parts, regex := splitSegment(key)
			So(ok, ShouldEqual, result.Ok)
			if result.Parts == nil {
				So(parts, ShouldBeNil)
			} else {
				So(parts, ShouldNotBeNil)
				So(strings.Join(parts, " "), ShouldEqual, strings.Join(result.Parts, " "))
			}
			So(regex, ShouldEqual, result.Regex)
		}
	})
}
