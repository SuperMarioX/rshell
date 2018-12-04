package commands

import (
	"fmt"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/crypt"
)

func Decrypt(text string) (string, error) {
	if text == "" {
		return "", nil
	}

	cfg := options.GetCfg()
	if cfg.Passcrypttype == "aes" {
		return crypt.AesDecrypt(text, cfg)
	} else {
		return "", fmt.Errorf("crypt type[%s] not support", cfg.Passcrypttype)
	}
}
