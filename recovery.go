// Copyright 2013 Martini Authors
// Copyright 2014 The Macaron Authors
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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"

	"github.com/go-macaron/inject"
)

const (
	panicHtml = `<html>
<head><title>PANIC: %s</title>
<meta charset="utf-8" />
<style type="text/css">
html, body {
	font-family: "Roboto", sans-serif;
	color: #333333;
	background-color: #ea5343;
	margin: 0px;
}
h1 {
	color: #d04526;
	background-color: #ffffff;
	padding: 20px;
	border-bottom: 1px dashed #2b3848;
}
pre {
	margin: 20px;
	padding: 20px;
	border: 2px solid #2b3848;
	background-color: #ffffff;
	white-space: pre-wrap;       /* css-3 */
	white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
	white-space: -pre-wrap;      /* Opera 4-6 */
	white-space: -o-pre-wrap;    /* Opera 7 */
	word-wrap: break-word;       /* Internet Explorer 5.5+ */
}
</style>
</head><body>
<h1>PANIC</h1>
<pre style="font-weight: bold;">%s</pre>
<pre>%s</pre>
</body>
</html>`
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
// While Martini is in development mode, Recovery will also output the panic as HTML.
func Recovery() Handler {
	return func(c *Context, log *log.Logger) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				log.Printf("PANIC: %s\n%s", err, stack)

				// Lookup the current responsewriter
				val := c.GetVal(inject.InterfaceOf((*http.ResponseWriter)(nil)))
				res := val.Interface().(http.ResponseWriter)

				// respond with panic message while in development mode
				var body []byte
				if Env == DEV {
					res.Header().Set("Content-Type", "text/html")
					body = []byte(fmt.Sprintf(panicHtml, err, err, stack))
				}

				res.WriteHeader(http.StatusInternalServerError)
				if nil != body {
					_, _ = res.Write(body)
				}
			}
		}()

		c.Next()
	}
}
