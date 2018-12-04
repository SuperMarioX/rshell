package client

import (
	"golang.org/x/crypto/ssh"
	"sync"
)

type safeMap struct {
	sync.Mutex
	data map[string]*ssh.Client
}

func newSafeMap() *safeMap {
	return &safeMap{
		data:  make(map[string]*ssh.Client),
	}
}

func (sm *safeMap) get(key string) *ssh.Client {
	sm.Lock()
	defer sm.Unlock()
	return sm.data[key]
}

func (sm *safeMap) set(key string, value *ssh.Client) {
	sm.Lock()
	defer sm.Unlock()
	sm.data[key] = value
}
