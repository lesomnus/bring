package http_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/lesomnus/bring/bringer/http"
	"github.com/lesomnus/bring/thing"
	"github.com/stretchr/testify/require"
)

type mockTransport struct {
	h http.Handler
}

func (t *mockTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	t.h.ServeHTTP(w, r)
	return w.Result(), nil
}

func TestHttpBringer(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		require := require.New(t)

		data := []byte("Royale with cheese")
		b := HttpBringer(WithTransport(&mockTransport{h: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(data))
		})}))

		f, err := b.Bring(context.Background(), thing.Thing{})
		if err == nil {
			defer f.Close()
		}
		require.NoError(err)

		v, err := io.ReadAll(f)
		require.NoError(err)
		require.Equal(v, data)
	})
	t.Run("not 200", func(t *testing.T) {
		require := require.New(t)

		data := []byte("Royale with cheese")
		b := HttpBringer(WithTransport(&mockTransport{h: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(data))
		})}))

		f, err := b.Bring(context.Background(), thing.Thing{})
		if err == nil {
			defer f.Close()
		}
		require.ErrorContains(err, "not 200")
	})
}
