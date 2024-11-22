package config

import (
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"github.com/lesomnus/bring/thing"
	"gopkg.in/yaml.v3"
)

type Entry struct {
	Next  map[string]*Entry
	Thing *thing.Thing
}

type EntryWalkFunc func(p string, t *thing.Thing)

func (e *Entry) IsLeaf() bool {
	return e != nil && e.Thing != nil
}

func (e *Entry) Len() int {
	l := 0
	for _, v := range e.Next {
		if !v.IsLeaf() {
			l += v.Len()
			continue
		}

		l += 1
	}

	return l
}

func (e *Entry) Walk(p string, f EntryWalkFunc) {
	if e.IsLeaf() {
		f(p, e.Thing)
		return
	}

	ks := maps.Keys(e.Next)
	for _, k := range slices.Sorted(ks) {
		p := filepath.Join(p, k)
		e.Next[k].Walk(p, f)
	}
}

func (e *Entry) UnmarshalYAML(n *yaml.Node) error {
	nodes := map[string]yaml.Node{}
	if err := n.Decode(nodes); err != nil {
		return err
	}

	e.Next = map[string]*Entry{}
	for name, node := range nodes {
		next := &Entry{}
		if strings.HasSuffix(name, "/") {
			if err := node.Decode(next); err != nil {
				return err
			}
		} else {
			t := &thing.Thing{}
			if err := node.Decode(t); err != nil {
				return err
			}

			next.Thing = t
		}

		e.Next[name] = next
	}

	return nil
}
