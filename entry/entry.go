package entry

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

type EntryWalkFunc func(p string, t *thing.Thing) error

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

func (e *Entry) Walk(p string, f EntryWalkFunc) error {
	if e.IsLeaf() {
		return f(p, e.Thing)
	}
	if e.Next == nil {
		return nil
	}

	ks := maps.Keys(e.Next)
	for _, k := range slices.Sorted(ks) {
		p := filepath.Join(p, k)
		if err := e.Next[k].Walk(p, f); err != nil {
			return err
		}
	}

	return nil
}

func (e *Entry) UnmarshalYAML(n *yaml.Node) error {
	nodes := map[string]yaml.Node{}
	if err := n.Decode(nodes); err != nil {
		return err
	}
	if len(nodes) == 0 {
		return nil
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
