// Copyright 2013-2015 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lru

import (
	"strconv"
	"testing"
	"time"
)

func getIntCheck(t *testing.T, c *Cache, key string, expectedValue int) {
	v, ok := c.Get(key)
	if !ok {
		t.Fatalf("cache doesn't contain %q", key)
	}
	x, ok := v.(int)
	if !ok {
		t.Fatalf("retrieved value %v is not int", x)
	}
	if x != expectedValue {
		t.Fatalf("wrong stored value: expected %d, got %d", expectedValue, x)
	}
}

func setInt(c *Cache, key string, value int) {
	c.Set(key, value, 0)
}

func TestMaxItems(t *testing.T) {
	c := New(Config{MaxItems: 3})
	setInt(c, "one", 1)
	setInt(c, "two", 2)
	setInt(c, "three", 3)

	// Try getting every, making sure we get "two" first,
	// so that it will later be pushed off the cache.

	getIntCheck(t, c, "two", 2)
	getIntCheck(t, c, "one", 1)
	getIntCheck(t, c, "three", 3)

	// Add new value, causing cache to drop "two".
	setInt(c, "four", 4)

	// Check that it was dropped...
	if v, ok := c.Get("two"); ok {
		t.Fatalf("cache didn't drop \"two\" (got %v)", v)
	}

	// ...and other elements are still there.
	getIntCheck(t, c, "one", 1)
	getIntCheck(t, c, "three", 3)
	getIntCheck(t, c, "four", 4)

	// Check that "one" is the oldest item.
	it, _ := c.OldestItem()
	if it.Key != "one" {
		t.Fatalf("oldest item is %s, expected %s", it.Key, "one")
	}

	// Replace element's value.
	setInt(c, "four", 100)
	getIntCheck(t, c, "four", 100)
}

func TestMaxBytes(t *testing.T) {
	removeCalled := 0
	c := New(Config{
		MaxBytes:      1000,
		RemoveHandler: func(it Item) { removeCalled++ },
	})
	b := make([]byte, 100)
	// Add 1100 bytes.
	for i := 0; i < 11; i++ {
		c.SetBytes(strconv.Itoa(i), b)
	}
	// Ensure there's no 0th item, as it should be dropped.
	if _, ok := c.GetBytes("0"); ok {
		t.Fatalf("cache didn't drop 0th item")
	}
	// Ensure items 1-10 exist.
	for i := 1; i < 11; i++ {
		if _, ok := c.GetBytes(strconv.Itoa(i)); !ok {
			t.Fatalf("cache item %d doesn't exist", i)
		}
	}
	if removeCalled != 1 {
		t.Fatalf("removeHandler was called %d times, expected %d", removeCalled, 1)
	}
}

func TestExpiration(t *testing.T) {
	c := New(Config{Expires: 1 * time.Millisecond})
	c.Set("hello", "world", 0)
	time.Sleep(2 * time.Millisecond)
	_, ok := c.Get("hello")
	if ok {
		t.Fatalf("didn't remove expired item")
	}
}

func benchSet(b *testing.B, config Config) {
	c := New(config)
	bs := make([]byte, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.SetBytes(strconv.Itoa(i), bs)
	}
}

func BenchmarkSet(b *testing.B) {
	benchSet(b, Config{MaxItems: 1000})
}

func benchGet(b *testing.B, config Config) {
	c := New(config)
	c.SetBytes("test", make([]byte, 100))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get("test")
	}
}

func BenchmarkGet(b *testing.B) {
	benchGet(b, Config{MaxItems: 1000})
}

func BenchmarkGetTrackTime(b *testing.B) {
	benchGet(b, Config{MaxItems: 1000, TrackAccessTime: true})
}
