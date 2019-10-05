//Package lock is provided with love by Raph
package lock

import "sync"

type KeyLock struct {
	giantLock sync.RWMutex
	locks     map[string]*sync.Mutex
}

func New() *KeyLock {
	return &KeyLock{
		giantLock: sync.RWMutex{},
		locks:     map[string]*sync.Mutex{},
	}
}

func (k *KeyLock) getLock(key string) *sync.Mutex {
	k.giantLock.RLock()
	if lock, ok := k.locks[key]; ok {
		k.giantLock.RUnlock()
		return lock
	}

	k.giantLock.RUnlock()
	k.giantLock.Lock()

	if lock, ok := k.locks[key]; ok {
		k.giantLock.Unlock()
		return lock
	}

	lock := &sync.Mutex{}
	k.locks[key] = lock
	k.giantLock.Unlock()
	return lock
}

func (k *KeyLock) Lock(key string) {
	k.getLock(key).Lock()
}

func (k *KeyLock) Unlock(key string) {
	k.getLock(key).Unlock()
}

func (k *KeyLock) KeyLocker(key string) sync.Locker {
	return k.getLock(key)
}
