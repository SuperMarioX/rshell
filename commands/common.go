package commands

import (
	"fmt"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/crypt"
	"github.com/luckywinds/rshell/types"
)

func getGroupAuthbyGroupname(groupname string) (types.Hostgroup, types.Auth, error) {
	if groupname == "" {
		return types.Hostgroup{}, types.Auth{}, fmt.Errorf("groupname[%s] empty", groupname)
	}

	hg, err := options.GetHostgroupByname(groupname)
	if err != nil {
		return types.Hostgroup{}, types.Auth{}, fmt.Errorf("group[%s] not found", groupname)
	}

	au, err := options.GetAuthByname(hg.Authmethod)
	if err != nil {
		return types.Hostgroup{}, types.Auth{}, fmt.Errorf("auth[%s] not found", hg.Authmethod)
	}

	if len(hg.Ips) == 0 {
		return types.Hostgroup{}, types.Auth{}, fmt.Errorf("hostgroup[%s] hosts empty", groupname)
	}

	return hg, au, nil
}

func getPlainPass(pass string, cfg types.Cfg) (string, error) {
	text, err := crypt.AesDecrypt(pass, cfg)
	if err != nil {
		return "", err
	}
	return text, nil
}