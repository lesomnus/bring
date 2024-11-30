package hook

import "io"

type tapMw struct {
	h Hook
}

func Tap(h Hook) Mw {
	return &tapMw{h: h}
}

func (h *tapMw) OnStart(next Hook) {
	h.h.OnStart()
	next.OnStart()
}
func (h *tapMw) OnSkip(next Hook) {
	h.h.OnSkip()
	next.OnSkip()
}
func (h *tapMw) OnDone(next Hook, r io.Reader) {
	h.h.OnDone(r)
	next.OnDone(r)
}
func (h *tapMw) OnError(next Hook, err error) {
	h.h.OnError(err)
	next.OnError(err)
}
func (h *tapMw) OnFinish(next Hook) {
	h.h.OnFinish()
	next.OnFinish()
}
