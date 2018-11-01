package checkers

import (
	"github.com/luckywinds/rshell/types"
	"github.com/scylladb/go-set/strset"
	"net"
)

func IsDuplicate(ss []string) bool {
	tempset := strset.New()
	for _, value := range ss {
		tempset.Add(value)
	}
	if tempset.Size() != len(ss) {
		return true
	}
	return false
}

func IsIpv4(s string) bool {
	temp := net.ParseIP(s)
	if temp == nil {
		return false
	}
	if temp.To4() == nil {
		return false
	}
	return true
}

func CheckHostgroupSize(h types.Hostgroup, max int) bool {
	if len(h.Ips) > max {
		return false
	}
	return true
}