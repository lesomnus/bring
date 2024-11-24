package hook

type NopHook struct{}

func (h *NopHook) OnStart()          {}
func (h *NopHook) OnSkip()           {}
func (h *NopHook) OnDone()           {}
func (h *NopHook) OnError(err error) {}
func (h *NopHook) OnFinish()         {}
