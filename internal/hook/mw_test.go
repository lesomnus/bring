package hook_test

import (
	"testing"

	"github.com/lesomnus/bring/internal/hook"
	"github.com/stretchr/testify/require"
)

type mockMw struct {
	hook.NopMw
	f func(next hook.Hook)
}

func (h *mockMw) OnStart(next hook.Hook) {
	h.f(next)
}

func TestTie(t *testing.T) {
	v := ""
	hook.Tie(
		&mockMw{f: func(next hook.Hook) { v += "a"; next.OnStart() }},
		&mockMw{f: func(next hook.Hook) { v += "b" }},
		&mockMw{f: func(next hook.Hook) { v += "c"; next.OnStart() }},
	).OnStart()
	require.Equal(t, "ab", v)
}
