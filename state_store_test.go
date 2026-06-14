package ratelimiter

import (
	"testing"
	"time"
)

func testMissingKey(t *testing.T, store StateStore) {
	t.Helper()

	state, version, exists, err := store.Get("test-key")

	if state != nil || version != 0 || exists != false || err != nil {
		t.Fatalf("empty key lookup resulted in conflicting values, %v, %d, %t, %v", state, version, exists, err)
	}
}

func testExistingKey(t *testing.T, store StateStore) {
	t.Helper()
	state := fixedWindowState{
		time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
		10,
	}

	ok, err := store.CompareAndSwap("key-1", 0, state)

	if !ok {
		t.Fatalf("cas returned false on a fresh store")
	}

	if err != nil {
		t.Fatal(err)
	}

	fetchedState, version, exists, err := store.Get("key-1")

	if err != nil {
		t.Fatal(err)
	}

	if version != 1 {
		t.Fatalf("Get returned wrong version, expected 1, got %d", version)
	}

	if !exists {
		t.Fatalf("Get return false exists flag on existing value.")
	}

	if fetchedState != state {
		t.Fatalf("Get returned wrong state, expected %v, get %v", state, fetchedState)
	}
}

func TestMissingKey(t *testing.T) {
	memStore := NewMemoryStore()
	ShardedMemoryStore, err := NewShardedMemoryStore(10)

	if err != nil {
		t.Fatal(err)
	}

	testMissingKey(t, memStore)
	testMissingKey(t, ShardedMemoryStore)
}

func TestExistingKey(t *testing.T) {
	memStore := NewMemoryStore()
	ShardedMemoryStore, err := NewShardedMemoryStore(10)

	if err != nil {
		t.Fatal(err)
	}

	testExistingKey(t, memStore)
	testExistingKey(t, ShardedMemoryStore)
}

// TODO: CAS on missing key with version 0 stores it, returns true, version becomes 1
// TODO: CAS on missing key with wrong version returns false
// TODO: CAS on existing key with correct version updates it, returns true, version increments
// TODO: CAS on existing key with stale version returns false
// TODO: concurrent access doesn't corrupt state
