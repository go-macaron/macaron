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
	"errors"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// StaticOptions is a struct for specifying configuration options for the macaron.Static middleware.
type StaticOptions struct {
	// Prefix is the optional prefix used to serve the static directory content
	Prefix string
	// SkipLogging will disable [Static] log messages when a static file is served.
	SkipLogging bool
	// IndexFile defines which file to serve as index if it exists.
	IndexFile string
	// Expires defines which user-defined function to use for producing a HTTP Expires Header
	// https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
	Expires func() string
	// FileSystem is the interface for supporting any implmentation of file system.
	FileSystem http.FileSystem
}

type staticMap struct {
	lock sync.RWMutex
	data map[string]*http.Dir
}

func (sm *staticMap) Set(dir *http.Dir) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	sm.data[string(*dir)] = dir
}

func (sm *staticMap) Get(name string) *http.Dir {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	return sm.data[name]
}

func (sm *staticMap) Delete(name string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	delete(sm.data, name)
}

var statics = staticMap{sync.RWMutex{}, map[string]*http.Dir{}}

var ErrMissingTrailingSlash = errors.New("Request missing trailing slash")

// staticFileSystem implements http.FileSystem interface.
type staticFileSystem struct {
	dir *http.Dir
}

func newStaticFileSystem(directory string) staticFileSystem {
	if !filepath.IsAbs(directory) {
		directory = filepath.Join(Root, directory)
	}
	dir := http.Dir(directory)
	statics.Set(&dir)
	return staticFileSystem{&dir}
}

func (fs staticFileSystem) Open(name string) (http.File, error) {
	return fs.dir.Open(name)
}

func prepareStaticOptions(dir string, options []StaticOptions) StaticOptions {
	var opt StaticOptions
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.IndexFile) == 0 {
		opt.IndexFile = "index.html"
	}
	// Normalize the prefix if provided
	if opt.Prefix != "" {
		// Ensure we have a leading '/'
		if opt.Prefix[0] != '/' {
			opt.Prefix = "/" + opt.Prefix
		}
		// Remove any trailing '/'
		opt.Prefix = strings.TrimRight(opt.Prefix, "/")
	}
	if opt.FileSystem == nil {
		opt.FileSystem = newStaticFileSystem(dir)
	}
	return opt
}

// Static returns a middleware handler that serves static files in the given directory.
func Static(directory string, staticOpt ...StaticOptions) Handler {
	// if !filepath.IsAbs(directory) {
	// 	directory = filepath.Join(Root, directory)
	// }
	// dir := http.Dir(directory)
	opt := prepareStaticOptions(directory, staticOpt)

	return func(ctx *Context, log *log.Logger) {
		// FIXME: BUG BUG BUG
		// ctx.statics[string(dir)] = &dir
		if ctx.Req.Method != "GET" && ctx.Req.Method != "HEAD" {
			return
		}
		file := ctx.Req.URL.Path
		// if we have a prefix, filter requests by stripping the prefix
		if opt.Prefix != "" {
			if !strings.HasPrefix(file, opt.Prefix) {
				return
			}
			file = file[len(opt.Prefix):]
			if file != "" && file[0] != '/' {
				return
			}
		}

		f, err := opt.FileSystem.Open(file)
		if err != nil {
			// FIXME: discard the error?
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return
		}

		// Try to serve index file
		if fi.IsDir() {
			// Redirect if missing trailing slash.
			if !strings.HasSuffix(ctx.Req.URL.Path, "/") {
				http.Redirect(ctx.Resp, ctx.Req.Request, ctx.Req.URL.Path+"/", http.StatusFound)
				return
			}

			file = path.Join(file, opt.IndexFile)
			f, err = opt.FileSystem.Open(file)
			if err != nil {
				// FIXME: discard the error?
				return
			}
			defer f.Close()

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				// FIXME: discard the error?
				return
			}
		}

		if !opt.SkipLogging {
			log.Println("[Static] Serving " + file)
		}

		// Add an Expires header to the static content
		if opt.Expires != nil {
			ctx.Resp.Header().Set("Expires", opt.Expires())
		}

		http.ServeContent(ctx.Resp, ctx.Req.Request, file, fi.ModTime(), f)
	}
}
