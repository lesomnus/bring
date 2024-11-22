package thing_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/lesomnus/bring/thing"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestThingYamlParse(t *testing.T) {
	http_url, _ := url.Parse("https://github.com")
	digest := digest.Digest("sha256:19ce1ab1f8b4e9de8f5e11885302f3b445dde71dd0d7c4ec0e8f4ace3baecffa")

	test := func(data string, f func(require *require.Assertions, v thing.Thing, err error)) func(t *testing.T) {
		return func(t *testing.T) {
			t.Log(data)
			require := require.New(t)

			var v thing.Thing
			err := yaml.Unmarshal([]byte(data), &v)

			// Thing is not fail to unmarshal unless it is not a map.
			require.NoError(err)

			f(require, v, err)
		}
	}

	t.Run("empty", test("{}", func(require *require.Assertions, v thing.Thing, err error) {
		require.Equal(thing.Thing{}, v)
	}))
	t.Run("url", test(fmt.Sprintf("url: %s", http_url.String()), func(require *require.Assertions, v thing.Thing, err error) {
		require.Equal(v.Url.String(), http_url.String())
	}))
	t.Run("url in invalid form", test(`
url:
  foo: bar
`, func(require *require.Assertions, v thing.Thing, err error) {
		require.Equal(v.Url.Scheme, thing.FailInvalid)
	}))
	t.Run("digest by string", test(fmt.Sprintf("digest: %s", digest.String()), func(require *require.Assertions, v thing.Thing, err error) {
		require.Equal(digest, v.Digest)
	}))
	t.Run("digest by object", test(fmt.Sprintf(`
digest:
  algo: sha256
  value: %s
`,
		digest.Encoded(),
	), func(require *require.Assertions, v thing.Thing, err error) {
		require.Equal(digest, v.Digest)
	}))
	t.Run("digest in invalid form", test(`
digest:
  - foo
  - bar
`, func(require *require.Assertions, v thing.Thing, err error) {
		require.Equal(v.Digest.Algorithm().String(), thing.FailInvalid)
	}))
}
