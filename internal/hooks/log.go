package hooks

import (
	"io"
	"log/slog"

	"github.com/lesomnus/bring/internal/task"
)

type LogHook struct {
	T task.Task
	L *slog.Logger
}

func (h *LogHook) OnStart() {
	h.L.Info("start",
		slog.String("from", h.T.Thing.Url.Redacted()),
		slog.String("to", h.T.Dest),
	)
}
func (h *LogHook) OnSkip() {
	h.L.Info("skip")
}
func (h *LogHook) OnDone(r io.Reader) {
	h.L.Info("done")
}
func (h *LogHook) OnError(err error) {
	h.L.Info("error", slog.String("message", err.Error()))
}
func (h *LogHook) OnFinish() {
	h.L.Info("finish")
}
