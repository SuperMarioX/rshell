package commands

import (
	"fmt"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/crypt"
)

func Encrypt(text string) (string, error) {
	if text == "" {
		return "", nil
	}

	cfg := options.GetCfg()
	if cfg.Passcrypttype == "aes" {
		return crypt.AesEncrypt(text, cfg)
	} else {
		return "", fmt.Errorf("crypt type[%s] not support", cfg.Passcrypttype)
	}
}
