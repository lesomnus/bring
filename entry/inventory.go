package entry

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Inventory struct {
	Things Entry
}

func FromPath(p string) (*Inventory, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	v := &Inventory{}
	if err := yaml.NewDecoder(f).Decode(v); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return v, nil
}
