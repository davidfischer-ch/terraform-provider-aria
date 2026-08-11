// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"
	"time"
)

func TestRWMutexKVLockUnlock(t *testing.T) {
	ctx := t.Context()
	kv := NewRWMutexKV()

	kv.Lock(ctx, "a")
	kv.Unlock(ctx, "a")

	// Re-locking the same key after unlocking must not block.
	kv.Lock(ctx, "a")
	kv.Unlock(ctx, "a")
}

func TestRWMutexKVSeparateKeysAreIndependent(t *testing.T) {
	ctx := t.Context()
	kv := NewRWMutexKV()

	kv.Lock(ctx, "a")
	defer kv.Unlock(ctx, "a")

	done := make(chan struct{})
	go func() {
		kv.Lock(ctx, "b")
		kv.Unlock(ctx, "b")
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("locking key \"b\" blocked while an unrelated key \"a\" was held")
	}
}

func TestRWMutexKVConcurrentReaders(t *testing.T) {
	ctx := t.Context()
	kv := NewRWMutexKV()

	kv.RLock(ctx, "a")
	defer kv.RUnlock(ctx, "a")

	done := make(chan struct{})
	go func() {
		kv.RLock(ctx, "a")
		kv.RUnlock(ctx, "a")
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("a second reader blocked; RLock must allow concurrent readers")
	}
}
