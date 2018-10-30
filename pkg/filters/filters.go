package filters

import (
	"fmt"
	"github.com/luckywinds/rshell/types"
	"strings"
)

func IsBlackCmd(cmd string, bcl []types.BlackCmd) bool {
	cstr := strings.TrimSpace(cmd)
	for _, bstr := range bcl {
		if bstr.Cmd != "" {
			if cstr == bstr.Cmd {
				fmt.Printf("DANGEROUS: Get a blacklist command [%s], DO NOT RUN.\n", strings.TrimSpace(cstr))
				return true
			}
		}
		if bstr.CmdPrefix != "" {
			if strings.HasPrefix(cstr, bstr.CmdPrefix) {
				fmt.Printf("DANGEROUS: Get a blacklist command [%s], DO NOT RUN.\n", strings.TrimSpace(cstr))
				return true
			}
		}
	}
	return false
}