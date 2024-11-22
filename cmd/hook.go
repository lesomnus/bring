package cmd

import (
	"fmt"
	"os"
	"time"
	"unicode"

	"github.com/fatih/color"
)

type ExecuteHook interface {
	OnStart()
	OnSkip()
	OnDone()
	OnFinish()
	OnError(err error)
}

type StdIoPrinterHook struct {
	ctx ExecuteContext
	t0  time.Time

	isSettled bool
}

func NewStdIoPrinterHook(ctx ExecuteContext) ExecuteHook {
	return &StdIoPrinterHook{
		ctx:       ctx,
		isSettled: false,
	}
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

func (h *StdIoPrinterHook) header(sym string) string {
	pb := color.New(color.FgHiBlack, color.Faint)
	ph := color.New(color.FgWhite)
	w := countDigits(h.ctx.N)
	v := pb.Sprint("[") + ph.Sprintf("%*d", w, h.ctx.I) +
		pb.Sprint("/") + fmt.Sprintf("%d", h.ctx.N) +
		pb.Sprint("|") + sym +
		pb.Sprint("]")

	return v
}

func (h *StdIoPrinterHook) path() string {
	return color.New(color.FgHiWhite).Sprint(h.ctx.Path)
}

func (h *StdIoPrinterHook) dt() string {
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

func (h *StdIoPrinterHook) OnStart() {
	h.t0 = time.Now()
}

func (h *StdIoPrinterHook) OnSkip() {
	h.isSettled = true

	sym := color.New(color.FgHiYellow).Sprint("=")
	fmt.Fprintf(os.Stdout, "%s %s • %s\n", h.header(sym), h.path(), h.dt())
}

func (h *StdIoPrinterHook) OnDone() {
	h.isSettled = true

	sym := color.New(color.FgHiGreen).Sprint("✓")
	fmt.Fprintf(os.Stdout, "%s %s • %s\n", h.header(sym), h.path(), h.dt())

}

func (h *StdIoPrinterHook) OnError(err error) {
	pr := color.New(color.FgHiRed)
	sym := pr.Sprint("!")
	fmt.Fprintf(os.Stdout, "%s %s • %s\n\t%s\n", h.header(sym), h.path(), h.dt(), pr.Sprint(err.Error()))
}

func (h *StdIoPrinterHook) OnFinish() {
	if h.isSettled {
		return
	}

	sym := "◦"
	fmt.Fprintf(os.Stdout, "%s %s\n", h.header(sym), h.path())
}
