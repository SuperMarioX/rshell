package filters

import (
	"github.com/luckywinds/rshell/types"
	"strings"
)

func IsBlackCmd(cmd string, bcl []types.BlackCmd) bool {
	cstr := strings.TrimSpace(cmd)
	for _, bstr := range bcl {
		if bstr.Cmd != "" {
			if cstr == bstr.Cmd {
				return true
			}
		}
		if bstr.CmdPrefix != "" {
			if strings.HasPrefix(cstr, bstr.CmdPrefix) {
				return true
			}
		}
	}
	return false
}