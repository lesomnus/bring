package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Dest   string
	Things Entry
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
