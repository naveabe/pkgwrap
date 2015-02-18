package initscript

import (
	"fmt"
)

type BasicRunnable struct {
	Path string            `json:"path" yaml:"path"`
	Args string            `json:"args" yaml:"args"`
	Env  map[string]string `json:"env" yaml:"env"`
}

func (b *BasicRunnable) Command() string {
	envStr := ""
	if b.Env != nil {
		for k, v := range b.Env {
			envStr += fmt.Sprintf("%s=%s ", k, v)
		}
	}
	return fmt.Sprintf("%s%s %s", envStr, b.Path, b.Args)
}
