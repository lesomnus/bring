package config

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type confCtxKey struct{}

type Config struct {
	Dest string
	Log  LogConfig

	Things Entry
}

func New() *Config {
	return &Config{
		Dest: "./inventory/",
		Log: LogConfig{
			Enabled: true,
			Format:  "simple",
			Level:   "warn",
		},
	}
}

func From(ctx context.Context) *Config {
	v, ok := ctx.Value(confCtxKey{}).(*Config)
	if ok {
		return nil
	}

	return v
}

func Into(ctx context.Context, conf *Config) context.Context {
	return context.WithValue(ctx, confCtxKey{}, conf)
}

func LoadFromFilepath(p string) (*Config, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	v := &Config{}
	if err := yaml.NewDecoder(f).Decode(v); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return v, nil
}
