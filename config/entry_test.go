package config_test

import (
	"testing"

	"github.com/lesomnus/bring/config"
	"github.com/lesomnus/bring/thing"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestEntryYamlParse(t *testing.T) {
	require := require.New(t)
	data := `
foo/:
  bar/:
    baz:
      url: https://baz.com
qux:
  url: https://qux.com
`

	root := config.Entry{}
	err := yaml.Unmarshal([]byte(data), &root)
	require.NoError(err)

	require.False(root.IsLeaf())
	require.Len(root.Next, 2)
	require.Contains(root.Next, "foo/")
	require.Contains(root.Next, "qux")

	foo := root.Next["foo/"]
	require.False(foo.IsLeaf())
	require.Len(foo.Next, 1)
	require.Contains(foo.Next, "bar/")

	bar := foo.Next["bar/"]
	require.False(bar.IsLeaf())
	require.Len(bar.Next, 1)
	require.Contains(bar.Next, "baz")

	baz := bar.Next["baz"]
	require.True(baz.IsLeaf())
	require.Equal("https://baz.com", baz.Thing.Url.String())

	qux := root.Next["qux"]
	require.True(qux.IsLeaf())
	require.Equal("https://qux.com", qux.Thing.Url.String())
}

func TestEntryWalk(t *testing.T) {
	require := require.New(t)
	data := `
f:
b/:
  c/:
    d:
  z/:
  e:
a:
`

	root := config.Entry{}
	err := yaml.Unmarshal([]byte(data), &root)
	require.NoError(err)

	visited := []string{}
	root.Walk("", func(p string, t *thing.Thing) error {
		visited = append(visited, p)
		return nil
	})
	require.Equal(
		visited,
		[]string{"a", "b/c/d", "b/e", "f"},
	)
}
