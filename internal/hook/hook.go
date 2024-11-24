package hook

type Hook interface {
	// Invoked when the task is started.
	OnStart()
	// Invoked when the digest matches with the file on the local.
	OnSkip()
	// Invoked when the task is succeeded.
	OnDone()
	// Invoked when the task got an error.
	OnError(err error)
	// Invoked when the task is finished.
	OnFinish()
}

type joinedHook struct {
	hs []Hook
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

func (h *joinedHook) OnDone() {
	for _, h := range h.hs {
		h.OnDone()
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

func Join(hs ...Hook) Hook {
	return &joinedHook{hs: hs}
}
