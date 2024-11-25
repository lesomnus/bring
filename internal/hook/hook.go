package hook

import "io"

type Hook interface {
	// Invoked when the task is started.
	OnStart()
	// Invoked when the digest matches with the file on the local.
	OnSkip()
	// Invoked when the task is succeeded.
	OnDone(r io.Reader)
	// Invoked when the task got an error.
	OnError(err error)
	// Invoked when the task is finished.
	OnFinish()
}

type joinedHook struct {
	hs []Hook
}

func Join(hs ...Hook) Hook {
	return &joinedHook{hs: hs}
}

func (h *joinedHook) OnStart() {
	for _, h := range h.hs {
		h.OnStart()
	}
}
func (h *joinedHook) OnSkip() {
	for _, h := range h.hs {
		h.OnSkip()
	}
}
func (h *joinedHook) OnDone(r io.Reader) {
	for _, h := range h.hs {
		h.OnDone(r)
	}
}
func (h *joinedHook) OnError(err error) {
	for _, h := range h.hs {
		h.OnError(err)
	}
}
func (h *joinedHook) OnFinish() {
	for _, h := range h.hs {
		h.OnFinish()
	}
}
