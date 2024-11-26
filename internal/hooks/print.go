package hooks

import (
	"fmt"
	"io"
	"strings"
	"time"
	"unicode"

	"github.com/fatih/color"
	"github.com/lesomnus/bring/internal/hook"
	"github.com/lesomnus/bring/internal/task"
)

type PrintHook struct {
	hook.NopHook
	T task.Task
	O io.Writer

	t0  time.Time
	sym string
	err error
}

func countDigits(n int) int {
	if n == 0 {
		return 1
	}
	count := 0
	for n > 0 {
		n /= 10
		count++
	}
	return count
}

func (h *PrintHook) writeStrings(b *strings.Builder, ss ...string) {
	for _, s := range ss {
		b.WriteString(s)
	}
}

func (h *PrintHook) path() string {
	return color.New(color.FgHiWhite).Sprint(h.T.Dest)
}

func (h *PrintHook) dt() string {
	t1 := time.Now()
	dt := t1.Sub(h.t0)

	pw := color.New(color.FgWhite, color.Faint)
	pb := color.New(color.FgHiBlack, color.Faint)
	pc := pw // First always a digit.
	s := []rune(dt.String())
	v := ""
	l := true // last was digit?
	p := 0
	for i, r := range s {
		lc := unicode.IsDigit(r)
		if l != lc {
			v += pc.Sprint(string(s[p:i]))
			// Swap paint
			if pc == pw {
				pc = pb
			} else {
				pc = pw
			}
			p = i
		}
		l = lc
	}
	// Last always not a digit.
	v += pb.Sprint(string(s[p:]))

	return v
}

func (h *PrintHook) OnStart() {
	h.t0 = time.Now()
}

func (h *PrintHook) OnSkip() {
	h.sym = color.New(color.FgHiYellow).Sprint("=")
}

func (h *PrintHook) OnDone(r io.Reader) {
	h.sym = color.New(color.FgHiGreen).Sprint("✓")
}

func (h *PrintHook) OnError(err error) {
	h.sym = color.New(color.FgHiRed).Sprint("!")
	h.err = err
}

func (h *PrintHook) OnFinish() {
	b := strings.Builder{}

	pb_ := color.New(color.FgHiBlack, color.Faint)
	pb := func(s string) string {
		return pb_.Sprint(s)
	}

	w := countDigits(h.T.Job.NumTasks)
	n := color.New(color.FgWhite).Sprintf("%*d", w, h.T.Order)
	N := fmt.Sprintf("%d", h.T.Job.NumTasks)

	// [n/N|sym] path • dt
	//             error message
	h.writeStrings(&b, pb("["), n, pb("/"), N, pb("|"), h.sym, pb("] "), h.path(), " • ", h.dt(), "\n")
	if h.err != nil {
		msg := color.New(color.FgHiRed).Sprint(h.err.Error())
		h.writeStrings(&b, strings.Repeat(" ", w*2+5+1+2), msg, "\n")
	}

	fmt.Fprint(h.O, b.String())
}
