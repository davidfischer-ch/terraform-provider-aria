// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Copied and the modified from the helpers that were removed from Terraform Plugin SDKv2
//
// See: https://developer.hashicorp.com/terraform/plugin/sdkv2/guides/...
//      v2-upgrade-guide#removal-of-helper-mutexkv-package

// RWMutexKV is a simple key/value store for arbitrary read-write mutexes.
//
// It can be used to serialize changes across arbitrary collaborators that share knowledge of the
// keys they must serialize on. The lock can be held by an arbitrary number of readers or a single
// writer. See https://pkg.go.dev/sync#RWMutex for further details.
type RWMutexKV struct {
	lock  sync.Mutex
	store map[string]*sync.RWMutex
}

// Returns a properly initialized RWMutexKV.
func NewRWMutexKV() *RWMutexKV {
	return &RWMutexKV{
		store: make(map[string]*sync.RWMutex),
	}
}

// Lock the named mutex for writing purposes (exclusive).
// If the lock is already locked for reading or writing, blocks until the lock is available.
func (self *RWMutexKV) Lock(ctx context.Context, key string) {
	tflog.Debug(ctx, "Write Locking "+key)
	self.get(key).Lock()
	tflog.Debug(ctx, "Write Locked "+key)
}

// Release a single lock for writing.
// It is a run-time error if the mutex is not locked for writing.
func (self *RWMutexKV) Unlock(ctx context.Context, key string) {
	tflog.Debug(ctx, "Write Unlocking "+key)
	self.get(key).Unlock()
	tflog.Debug(ctx, "Write Unlocked "+key)
}

// Lock the named mutex for reading purposes (concurrent friendly).
// If the lock is already locked for writing, blocks until the lock is available.
func (self *RWMutexKV) RLock(ctx context.Context, key string) {
	tflog.Debug(ctx, "Read Locking "+key)
	self.get(key).RLock()
	tflog.Debug(ctx, "Read Locked "+key)
}

// Release a single lock for reading.
// It is a run-time error if the mutex is not locked for reading.
func (self *RWMutexKV) RUnlock(ctx context.Context, key string) {
	tflog.Debug(ctx, "Read Unlocking "+key)
	self.get(key).RUnlock()
	tflog.Debug(ctx, "Read Unlocked "+key)
}

// Return a mutex for the given key, no guarantee of its lock status.
func (self *RWMutexKV) get(key string) *sync.RWMutex {
	self.lock.Lock()
	defer self.lock.Unlock()
	mutex, ok := self.store[key]
	if !ok {
		mutex = &sync.RWMutex{}
		self.store[key] = mutex
	}
	return mutex
}
