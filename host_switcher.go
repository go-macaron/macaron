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
	"log"
	"net/http"
)

// HostSwitcher represents a global multi-site support layer.
type HostSwitcher struct {
	switches map[string]*Macaron
}

// NewHostSwitcher initalizes and returns a new host switcher.
// You have to use this function to get a new host switcher.
func NewHostSwitcher() *HostSwitcher {
	return &HostSwitcher{
		switches: make(map[string]*Macaron),
	}
}

// Set adds a new switch to host switcher.
func (hs *HostSwitcher) Set(host string, m *Macaron) {
	hs.switches[host] = m
}

// Remove removes a switch from host switcher.
func (hs *HostSwitcher) Remove(host string) {
	delete(hs.switches, host)
}

// ServeHTTP is the HTTP Entry point for a Host Switcher instance.
func (hs *HostSwitcher) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if h := hs.switches[req.Host]; h != nil {
		h.ServeHTTP(resp, req)
	} else {
		http.Error(resp, "Not Found", http.StatusNotFound)
	}
}

// RunOnAddr runs server in given address and port.
func (hs *HostSwitcher) RunOnAddr(addr string) {
	log.Fatalln(http.ListenAndServe(addr, hs))
}

// Run the http server. Listening on os.GetEnv("PORT") or 4000 by default.
func (hs *HostSwitcher) Run() {
	hs.RunOnAddr(getDefaultListenAddr())
}
