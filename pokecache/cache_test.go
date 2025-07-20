package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	c := NewCache(5 * time.Minute)

	key := "testKey"
	val := []byte("testValue")

	c.Add(key, val)
	got, ok := c.Get(key)
	if !ok {
		t.Fatalf("expected key %s to exist", key)
	}
	if string(got) != string(val) {
		t.Fatalf("expected %s, got %s", val, got)
	}
}

func TestCacheExpiration(t *testing.T) {
	interval := 10 * time.Millisecond
	c := NewCache(interval)

	key := "testKey"
	val := []byte("testValue")

	c.Add(key, val)
	time.Sleep(interval + 5*time.Millisecond)

	_, ok := c.Get(key)
	if ok {
		t.Fatalf("expected key %s to be expired", key)
	}
}

func TestCacheReapLoop(t *testing.T) {
	interval := 10 * time.Millisecond
	c := NewCache(interval)

	key := "testKey"
	val := []byte("testValue")

	c.Add(key, val)
	time.Sleep(interval + 5*time.Millisecond)

	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.entry[key]; exists {
		t.Fatalf("expected key %s to be reaped", key)
	}
}