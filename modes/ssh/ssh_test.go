package ssh

import (
	"testing"
)

const (
	username   = "root"
	privatekey = "E:\\tmp\\id_rsa"
	passphrase = "123456"
	password   = "Cloud12#$"
	ip         = "192.168.31.63"
	port       = 22
)

func TestSSHRun(t *testing.T) {
	ciphers := []string{}
	cmds := []string{"hostname", "whoami"}
	stdout, stderr, err := DO("", ip, port, username, password, "", "", 10, ciphers, cmds)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(stdout)
	t.Log(stderr)
	stdout, stderr, err = DO("", ip, port, username, password, "", "", 10, ciphers, cmds)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(stdout)
	t.Log(stderr)
	stdout, stderr, err = DO("", ip, port, username, password, "", "", 10, ciphers, cmds)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(stdout)
	t.Log(stderr)
	stdout, stderr, err = DO("", ip, port, username, password, "", "", 10, ciphers, cmds)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(stdout)
	t.Log(stderr)
	stdout, stderr, err = DO("", ip, port, username, password, "", "", 10, ciphers, cmds)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(stdout)
	t.Log(stderr)
	return
}
