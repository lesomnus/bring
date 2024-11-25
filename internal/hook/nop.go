package hook

import "io"

type NopHook struct{}

func (h *NopHook) OnStart()           {}
func (h *NopHook) OnSkip()            {}
func (h *NopHook) OnDone(r io.Reader) {}
func (h *NopHook) OnError(err error)  {}
func (h *NopHook) OnFinish()          {}

type NopMw struct{}

func (h *NopMw) OnStart(next Hook)             { next.OnStart() }
func (h *NopMw) OnSkip(next Hook)              { next.OnSkip() }
func (h *NopMw) OnDone(next Hook, r io.Reader) { next.OnDone(r) }
func (h *NopMw) OnError(next Hook, err error)  { next.OnError(err) }
func (h *NopMw) OnFinish(next Hook)            { next.OnFinish() }
