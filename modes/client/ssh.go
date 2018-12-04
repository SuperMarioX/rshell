package client

import (
	"fmt"
	"github.com/luckywinds/rshell/pkg/checkers"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"time"
)

var dialcache = newSafeMap()

func New(groupname, host string, port int, user, pass, keyname, passphrase string, timeout int, ciphers []string) (*ssh.Client, error) {
	if groupname == "" {
		groupname = "DEFAULT"
	}
	if !checkers.IsIpv4(host) || port <= 0 || port > 65535 || user == "" {
		return nil, fmt.Errorf("host[%s] or port[%d] or user[%s] illegal", host, port, user)
	}
	if pass == "" && keyname == "" {
		return nil, fmt.Errorf("pass and keyname can not be empty")
	}
	if timeout < 0 || timeout > 600 {
		return nil, fmt.Errorf("timeout[%d] illegal", timeout)
	}
	if len(ciphers) == 0 {
		ciphers = []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"}
	}

	cachekey := groupname + "/" + host + ":" + fmt.Sprintf("%d", port)
	if v := dialcache.get(cachekey); v != nil {
		return v, nil
	}

	var err error
	auth := make([]ssh.AuthMethod, 0)
	if pass != "" {
		auth = append(auth, ssh.Password(pass))
	}
	if keyname != "" {
		var (
			pemBytes []byte
			signer   ssh.Signer
		)
		pemBytes, err = ioutil.ReadFile(keyname)
		if err != nil {
			return nil, err
		}
		if passphrase == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(passphrase))
		}
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: time.Duration(timeout) * time.Second,
		Config: ssh.Config{
			Ciphers: ciphers,
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), clientConfig)
	if err != nil {
		return nil, err
	}

	dialcache.set(cachekey, client)
	return client, nil
}
