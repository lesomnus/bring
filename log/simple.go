package log

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type SimpleHandler struct {
	slog.Handler
	w io.Writer
	m *sync.Mutex
}

func NewSimpleHandler(w io.Writer, opts *slog.HandlerOptions) *SimpleHandler {
	var o slog.HandlerOptions
	if opts != nil {
		o = *opts
	}

	o.ReplaceAttr = replaceAttr(o.ReplaceAttr)

	return &SimpleHandler{
		Handler: slog.NewTextHandler(w, &o),

		w: w,
		m: &sync.Mutex{},
	}
}

func replaceAttr(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.TimeKey:
			fallthrough
		case slog.LevelKey:
			fallthrough
		case slog.MessageKey:
			return slog.Attr{}
		}

		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

var (
	p_time  = color.New(color.FgWhite, color.Faint)
	p_debug = color.New(color.FgHiGreen)
	p_info  = color.New(color.FgHiBlue)
	p_warn  = color.New(color.FgHiYellow)
	p_error = color.New(color.FgHiRed)
	p_msg   = color.New(color.FgHiWhite)
)

func (h *SimpleHandler) Handle(ctx context.Context, r slog.Record) error {
	b := strings.Builder{}
	b.WriteString(p_time.Sprint(r.Time.Format("15:04:05.000")))

	var sym string
	switch r.Level {
	case slog.LevelDebug:
		sym = p_debug.Sprint(" ? ")
	case slog.LevelInfo:
		sym = p_info.Sprint(" i ")
	case slog.LevelWarn:
		sym = p_warn.Sprint(" ! ")
	case slog.LevelError:
		sym = p_error.Sprint(" x ")
	default:
		sym = p_error.Sprint("   ")
	}
	b.WriteString(sym)
	b.WriteString(p_msg.Sprint(r.Message))
	b.WriteString(" ")

	h.m.Lock()
	defer h.m.Unlock()

	h.w.Write([]byte(b.String()))
	err := h.Handler.Handle(ctx, r)
	return err
}

func (h *SimpleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SimpleHandler{Handler: h.Handler.WithAttrs(attrs), w: h.w, m: h.m}
}

func (h *SimpleHandler) WithGroup(name string) slog.Handler {
	return &SimpleHandler{Handler: h.Handler.WithGroup(name), w: h.w, m: h.m}
}
