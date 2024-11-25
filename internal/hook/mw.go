package hook

import "io"

type Mw interface {
	OnStart(next Hook)
	OnSkip(next Hook)
	OnDone(next Hook, r io.Reader)
	OnError(next Hook, err error)
	OnFinish(next Hook)
}

type tiedHook struct {
	hs []Mw
}

func Tie(hs ...Mw) Hook {
	hs = append(hs, &terminalMw{})
	return &tiedHook{hs: hs}
}

func (h *tiedHook) OnStart() {
	curr := h.hs[0]
	curr.OnStart(&tiedHook{hs: h.hs[1:]})
}
func (h *tiedHook) OnSkip() {
	curr := h.hs[0]
	curr.OnSkip(&tiedHook{hs: h.hs[1:]})
}
func (h *tiedHook) OnDone(r io.Reader) {
	curr := h.hs[0]
	curr.OnDone(&tiedHook{hs: h.hs[1:]}, r)
}
func (h *tiedHook) OnError(err error) {
	curr := h.hs[0]
	curr.OnError(&tiedHook{hs: h.hs[1:]}, err)
}
func (h *tiedHook) OnFinish() {
	curr := h.hs[0]
	curr.OnFinish(&tiedHook{hs: h.hs[1:]})
}

type terminalMw struct{}

func (h *terminalMw) OnStart(next Hook)             {}
func (h *terminalMw) OnSkip(next Hook)              {}
func (h *terminalMw) OnDone(next Hook, r io.Reader) {}
func (h *terminalMw) OnError(next Hook, err error)  {}
func (h *terminalMw) OnFinish(next Hook)            {}
