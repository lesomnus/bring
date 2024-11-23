package config

import (
	"log/slog"
	"os"

	"github.com/lesomnus/bring/log"
)

type LogConfig struct {
	Enabled bool
	Format  string // "text" | "json" | "simple"
	Level   string // "error" | "warn" | "info" | "debug"
}

func (c *LogConfig) Logger() (l *slog.Logger) {
	if !c.Enabled {
		return log.Discard
	}

	opt := &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}
	switch c.Level {
	case "debug":
		opt.Level = slog.LevelDebug
	case "info":
		opt.Level = slog.LevelInfo
	case "":
	case "warn":
		opt.Level = slog.LevelWarn
	case "error":
		opt.Level = slog.LevelError

	default:
		defer func() {
			l.Warn("unknown log level", slog.String("value", c.Level))
		}()
	}

	o := os.Stderr

	var h slog.Handler = log.NewSimpleHandler(o, opt)
	switch c.Format {
	case "text":
		h = slog.NewTextHandler(o, opt)
	case "json":
		h = slog.NewJSONHandler(o, opt)
	case "":
	case "simple":
		h = log.NewSimpleHandler(o, opt)

	default:
		defer func() {
			l.Warn("unknown log format", slog.String("value", c.Format))
		}()
	}

	l = slog.New(h)
	return
}
