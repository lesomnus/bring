package hook

import "io"

type forwardMw struct {
	h Hook
}

func Forward(h Hook) Mw {
	return &forwardMw{h: h}
}

func (h *forwardMw) OnStart(next Hook) {
	h.h.OnStart()
	next.OnStart()
}
func (h *forwardMw) OnSkip(next Hook) {
	h.h.OnSkip()
	next.OnSkip()
}
func (h *forwardMw) OnDone(next Hook, r io.Reader) {
	h.h.OnDone(r)
	next.OnDone(r)
}
func (h *forwardMw) OnError(next Hook, err error) {
	h.h.OnError(err)
	next.OnError(err)
}
func (h *forwardMw) OnFinish(next Hook) {
	h.h.OnFinish()
	next.OnFinish()
}
