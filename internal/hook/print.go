package hook

import (
	"fmt"
	"io"
	"time"
	"unicode"

	"github.com/fatih/color"
	"github.com/lesomnus/bring/internal/task"
)

type PrintHook struct {
	o    io.Writer
	task task.Task

	isSettled bool

	t0 time.Time
}

func NewPrintHook(w io.Writer, t task.Task) Hook {
	return &PrintHook{
		o:    w,
		task: t,
	}
}

func (h *PrintHook) print(format string, a ...any) {
	fmt.Fprintf(h.o, format, a...)
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

func (h *PrintHook) header(sym string) string {
	pb := color.New(color.FgHiBlack, color.Faint)
	ph := color.New(color.FgWhite)
	w := countDigits(h.task.Job.NumTasks)
	v := pb.Sprint("[") + ph.Sprintf("%*d", w, h.task.Order) +
		pb.Sprint("/") + fmt.Sprintf("%d", h.task.Job.NumTasks) +
		pb.Sprint("|") + sym +
		pb.Sprint("]")

	return v
}

func (h *PrintHook) path() string {
	return color.New(color.FgHiWhite).Sprint(h.task.Dest)
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
	h.isSettled = true

	sym := color.New(color.FgHiYellow).Sprint("=")
	h.print("%s %s • %s\n", h.header(sym), h.path(), h.dt())
}

func (h *PrintHook) OnDone() {
	h.isSettled = true

	sym := color.New(color.FgHiGreen).Sprint("✓")
	h.print("%s %s • %s\n", h.header(sym), h.path(), h.dt())

}

func (h *PrintHook) OnError(err error) {
	h.isSettled = true

	pr := color.New(color.FgHiRed)
	sym := pr.Sprint("!")
	h.print("%s %s ∙ %s\n\t%s\n", h.header(sym), h.path(), h.dt(), pr.Sprint(err.Error()))
}

func (h *PrintHook) OnFinish() {
	if h.isSettled {
		return
	}

	h.print("%s %s ∙ %s\n", h.header("◦"), h.path(), h.dt())
}
