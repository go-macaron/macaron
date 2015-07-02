// Copyright 2015 Unknwon
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
	"regexp"
	"strings"

	"github.com/Unknwon/com"
)

type patternType int8

const (
	_PATTERN_STATIC patternType = iota
	_PATTERN_REGEXP
	_PATTERN_PATH_EXT
	_PATTERN_HOLDER
	_PATTERN_MATCH_ALL
)

type Leaf struct {
	parent *Tree

	typ       patternType
	pattern   string
	wildcards []string
	reg       *regexp.Regexp
	optional  bool

	name   string
	handle Handle
}

var wildcardPattern = regexp.MustCompile(`:[a-zA-Z0-9]+`)

// getNextWildcard tries to find next wildcard and update pattern with corresponding regexp.
func getNextWildcard(pattern string) (wildcard string, _ string) {
	pos := wildcardPattern.FindStringIndex(pattern)
	if pos == nil {
		return "", pattern
	}
	wildcard = pattern[pos[0]:pos[1]]

	// Reach last character or no regexp is given.
	if len(pattern) == pos[1] {
		return wildcard, strings.Replace(pattern, wildcard, `(.+)`, 1)
	} else if pattern[pos[1]] != '(' {
		if len(pattern) >= pos[1]+4 && pattern[pos[1]:pos[1]+4] == ":int" {
			pattern = strings.Replace(pattern, ":int", "([0-9]+)", -1)
		} else {
			return wildcard, strings.Replace(pattern, wildcard, `(.+)`, 1)
		}
	}

	// Cut out placeholder directly.
	return wildcard, pattern[:pos[0]] + pattern[pos[1]:]
}

func getWildcards(pattern string) (string, []string) {
	wildcards := make([]string, 0, 2)

	// Keep getting next wildcard until nothing is left.
	var wildcard string
	for {
		wildcard, pattern = getNextWildcard(pattern)
		if len(wildcard) > 0 {
			wildcards = append(wildcards, wildcard)
		} else {
			break
		}
	}

	return pattern, wildcards
}

func checkPattern(pattern string) (typ patternType, wildcards []string, reg *regexp.Regexp) {
	pattern = strings.TrimLeft(pattern, "?")
	if pattern == "*" {
		typ = _PATTERN_MATCH_ALL
	} else if pattern == "*.*" {
		typ = _PATTERN_PATH_EXT
	} else if strings.Contains(pattern, ":") {
		typ = _PATTERN_REGEXP
		pattern, wildcards = getWildcards(pattern)
		if pattern == "(.+)" {
			typ = _PATTERN_HOLDER
		} else {
			reg = regexp.MustCompile(pattern)
		}
	}
	return typ, wildcards, reg
}

func NewLeaf(parent *Tree, pattern, name string, handle Handle) *Leaf {
	typ, wildcards, reg := checkPattern(pattern)
	optional := false
	if len(pattern) > 0 && pattern[0] == '?' {
		optional = true
	}
	return &Leaf{parent, typ, pattern, wildcards, reg, optional, name, handle}
}

// Tree represents a router tree in Macaron.
type Tree struct {
	parent *Tree

	typ       patternType
	pattern   string
	wildcards []string
	reg       *regexp.Regexp

	subtrees []*Tree
	leaves   []*Leaf
}

func NewSubtree(parent *Tree, pattern string) *Tree {
	typ, wildcards, reg := checkPattern(pattern)
	return &Tree{parent, typ, pattern, wildcards, reg, make([]*Tree, 0, 5), make([]*Leaf, 0, 5)}
}

func NewTree() *Tree {
	return NewSubtree(nil, "")
}

func (t *Tree) addLeaf(pattern, name string, handle Handle) bool {
	for i := 0; i < len(t.leaves); i++ {
		if t.leaves[i].pattern == pattern {
			return true
		}
	}

	leaf := NewLeaf(t, pattern, name, handle)

	// Add exact same leaf to grandparent/parent level without optional.
	if leaf.optional {
		parent := leaf.parent
		if parent.parent != nil {
			parent.parent.addLeaf(parent.pattern, name, handle)
		} else {
			parent.addLeaf("", name, handle) // Root tree can add as empty pattern.
		}
	}

	i := 0
	for ; i < len(t.leaves); i++ {
		if leaf.typ < t.leaves[i].typ {
			break
		}
	}

	if i == len(t.leaves) {
		t.leaves = append(t.leaves, leaf)
	} else {
		t.leaves = append(t.leaves[:i], append([]*Leaf{leaf}, t.leaves[i:]...)...)
	}
	return false
}

func (t *Tree) addSubtree(segment, pattern, name string, handle Handle) bool {
	for i := 0; i < len(t.subtrees); i++ {
		if t.subtrees[i].pattern == segment {
			return t.subtrees[i].addNextSegment(pattern, name, handle)
		}
	}

	subtree := NewSubtree(t, segment)
	i := 0
	for ; i < len(t.subtrees); i++ {
		if subtree.typ < t.subtrees[i].typ {
			break
		}
	}

	if i == len(t.subtrees) {
		t.subtrees = append(t.subtrees, subtree)
	} else {
		t.subtrees = append(t.subtrees[:i], append([]*Tree{subtree}, t.subtrees[i:]...)...)
	}
	return subtree.addNextSegment(pattern, name, handle)
}

func (t *Tree) addNextSegment(pattern, name string, handle Handle) bool {
	pattern = strings.TrimPrefix(pattern, "/")

	i := strings.Index(pattern, "/")
	if i == -1 {
		return t.addLeaf(pattern, name, handle)
	}
	return t.addSubtree(pattern[:i], pattern[i+1:], name, handle)
}

func (t *Tree) Add(pattern, name string, handle Handle) bool {
	pattern = strings.TrimSuffix(pattern, "/")
	return t.addNextSegment(pattern, name, handle)
}

func (t *Tree) matchLeaf(globLevel int, url string, params Params) (Handle, bool) {
	for i := 0; i < len(t.leaves); i++ {
		switch t.leaves[i].typ {
		case _PATTERN_STATIC:
			if t.leaves[i].pattern == url {
				return t.leaves[i].handle, true
			}
		case _PATTERN_REGEXP:
			results := t.leaves[i].reg.FindStringSubmatch(url)
			// Number of results and wildcasrd should be exact same.
			if len(results)-1 != len(t.leaves[i].wildcards) {
				break
			}

			for j := 0; j < len(t.leaves[i].wildcards); j++ {
				params[t.leaves[i].wildcards[j]] = results[j+1]
			}
			return t.leaves[i].handle, true
		case _PATTERN_PATH_EXT:
			j := strings.LastIndex(url, ".")
			if j > -1 {
				params[":path"] = url[:j]
				params[":ext"] = url[j+1:]
			} else {
				params[":path"] = url
			}
			return t.leaves[i].handle, true
		case _PATTERN_HOLDER:
			params[t.leaves[i].wildcards[0]] = url
			return t.leaves[i].handle, true
		case _PATTERN_MATCH_ALL:
			params["*"+com.ToStr(globLevel)] = url
			return t.leaves[i].handle, true
		}
	}
	return nil, false
}

func (t *Tree) matchSubtree(globLevel int, segment, url string, params Params) (Handle, bool) {
	for i := 0; i < len(t.subtrees); i++ {
		switch t.subtrees[i].typ {
		case _PATTERN_STATIC:
			if t.subtrees[i].pattern == segment {
				if handle, ok := t.subtrees[i].matchNextSegment(globLevel, url, params); ok {
					return handle, true
				}
			}
		case _PATTERN_REGEXP:
			results := t.subtrees[i].reg.FindStringSubmatch(segment)
			if len(results)-1 != len(t.subtrees[i].wildcards) {
				break
			}

			for j := 0; j < len(t.subtrees[i].wildcards); j++ {
				params[t.subtrees[i].wildcards[j]] = results[j+1]
			}
			if handle, ok := t.subtrees[i].matchNextSegment(globLevel, url, params); ok {
				return handle, true
			}
		case _PATTERN_HOLDER:
			if handle, ok := t.subtrees[i].matchNextSegment(globLevel+1, url, params); ok {
				params[t.subtrees[i].wildcards[0]] = segment
				return handle, true
			}
		case _PATTERN_MATCH_ALL:
			if handle, ok := t.subtrees[i].matchNextSegment(globLevel+1, url, params); ok {
				params["*"+com.ToStr(globLevel)] = segment
				return handle, true
			}
		}
	}

	if len(t.leaves) > 0 {
		leaf := t.leaves[len(t.leaves)-1]
		if leaf.typ == _PATTERN_PATH_EXT {
			url = segment + "/" + url
			j := strings.LastIndex(url, ".")
			if j > -1 {
				params[":path"] = url[:j]
				params[":ext"] = url[j+1:]
			} else {
				params[":path"] = url
			}
			return leaf.handle, true
		} else if leaf.typ == _PATTERN_MATCH_ALL {
			params["*"+com.ToStr(globLevel)] = segment + "/" + url
			return leaf.handle, true
		}
	}
	return nil, false
}

func (t *Tree) matchNextSegment(globLevel int, url string, params Params) (Handle, bool) {
	i := strings.Index(url, "/")
	if i == -1 {
		return t.matchLeaf(globLevel, url, params)
	}
	return t.matchSubtree(globLevel, url[:i], url[i+1:], params)
}

func (t *Tree) Match(url string) (Handle, Params, bool) {
	url = strings.TrimPrefix(url, "/")
	url = strings.TrimSuffix(url, "/")
	params := make(Params)
	handle, ok := t.matchNextSegment(0, url, params)
	return handle, params, ok
}

// MatchTest returns true if given URL is matched by given pattern.
func MatchTest(pattern, url string) bool {
	t := NewTree()
	t.Add(pattern, "", nil)
	_, _, ok := t.Match(url)
	return ok
}
