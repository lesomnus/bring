package config

import "time"

type BringConfig struct {
	BringTimeout time.Duration `yaml:"bring_timeout"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
}
