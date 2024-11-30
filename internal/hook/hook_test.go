package hook_test

import (
	"testing"

	"github.com/lesomnus/bring/internal/hook"
	"github.com/stretchr/testify/require"
)

type mockHook struct {
	hook.NopHook
	f func()
}

func (h *mockHook) OnStart() {
	h.f()
}

func TestJoin(t *testing.T) {
	v := ""
	hook.Join(
		&mockHook{f: func() { v += "a" }},
		&mockHook{f: func() { v += "b" }},
		&mockHook{f: func() { v += "c" }},
	).OnStart()
	require.Equal(t, "abc", v)
}
