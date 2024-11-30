package hook_test

import (
	"testing"

	"github.com/lesomnus/bring/internal/hook"
	"github.com/stretchr/testify/require"
)

func TestTap(t *testing.T) {
	v := ""
	hook.Tie(
		hook.Tap(&mockHook{f: func() { v += "a" }}),
		hook.Tap(&mockHook{f: func() { v += "b" }}),
		hook.Tap(&mockHook{f: func() { v += "c" }}),
	).OnStart()
	require.Equal(t, "abc", v)
}
