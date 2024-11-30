package thing_test

import (
	"testing"

	"github.com/lesomnus/bring/thing"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestThingYamlParse(t *testing.T) {
	test := func(data string, f func(require *require.Assertions, v thing.Thing, err error)) func(t *testing.T) {
		return func(t *testing.T) {
			t.Log(data)
			require := require.New(t)

			var v thing.Thing
			err := yaml.Unmarshal([]byte(data), &v)
			f(require, v, err)
		}
	}

	t.Run("valid", test(`
url: https://github.com
digest: sha256:19ce1ab1f8b4e9de8f5e11885302f3b445dde71dd0d7c4ec0e8f4ace3baecffa`,
		func(require *require.Assertions, v thing.Thing, err error) {
			require.NoError(err)
			require.Equal("https", v.Url.Scheme)
			require.Equal("github.com", v.Url.Host)
			require.NoError(v.Digest.Validate())
			require.Equal(digest.SHA256, v.Digest.Algorithm())
		},
	))
	t.Run("empty", test("{}", func(require *require.Assertions, v thing.Thing, err error) {
		require.ErrorContains(err, "url must be set")
	}))
	t.Run("no digest", test("url: foo", func(require *require.Assertions, v thing.Thing, err error) {
		require.NoError(err)
	}))
	t.Run("invalid digest", test(`
url: foo
digest: sha256:foo`,
		func(require *require.Assertions, v thing.Thing, err error) {
			require.ErrorContains(err, "invalid digest")
		},
	))
}
