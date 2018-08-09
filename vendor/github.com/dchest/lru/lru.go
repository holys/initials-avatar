// Copyright 2013-2015 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lru implements Least Recently Used cache algorithm.
//
// Cache capacity can be optionally limited by both size in bytes and a number
// of items. Items can keep track of their modification and access time.
//
// Example:
//
//  // Create a 1 MB cache.
//  c := lru.New(lru.Config{ MaxBytes: 1024*1024 })
//  ...
//  // Query cache and insert item if it's not there.
//  var x someType
//  v, ok := c.Get("key")
//  if !ok {
//  	// Item is not in cache, fetch it from main storage...
//  	x = ...
//  	// ...then add it to cache.
//  	c.Set("key", x, x.Size())
//  } else {
//  	x = v.(someType)
//  }
//  // x now contains the value.
//
package lru

import (
	"container/list"
	"math"
	"sync"
	"time"
)

// Key is any comparable value.
type Key interface{}

type Cache struct {
	sync.Mutex

	m map[Key]*list.Element
	l *list.List

	size int64

	config Config
}

type Config struct {
	// Maximum number of items the cache can contain (unlimited by default).
	MaxItems int

	// Maximum byte capacity of cache (unlimited by default).
	MaxBytes int64

	// Track item access time (false by default).
	TrackAccessTime bool

	// Item expiration duration.
	//
	// An item is removed from cache when trying to get it if the given
	// time passed since its modification time.
	//
	// Set to zero for no expiration (default).
	Expires time.Duration

	// Function to call when an item is dropped or removed from cache
	// (nil by default).
	//
	// The handler function is given a copy of the removed item as an argument.
	// During its execution the cache is locked, and must not be accessed or
	// modified from the handler function.
	RemoveHandler func(Item)
}

// Item represents a cached item with additional information.
type Item struct {
	Key        Key
	Value      interface{}
	Size       int64     // byte size of value
	ModTime    time.Time // when this item was added to cache
	AccessTime time.Time // when this item was last accessed
}

// New returns a cache instance configured with the given Config.
func New(config Config) *Cache {
	return &Cache{
		m:      make(map[Key]*list.Element, config.MaxItems),
		l:      list.New(),
		config: config,
	}
}

// Reset clears the cache.
func (c *Cache) Reset() {
	c.Lock()
	defer c.Unlock()
	c.m = make(map[Key]*list.Element, c.config.MaxItems)
	c.l.Init()
	c.size = 0
}

// enforceCapacity ensures that the cache is not over its capacity,
// dropping items if needed.
// Cache must be locked.
func (c *Cache) enforceCapacity() {
	// Check maxItems.
	if c.config.MaxItems > 0 {
		for c.l.Len() > c.config.MaxItems {
			if !c.dropTail() {
				break
			}
		}
	}
	// Check maxBytes.
	if c.config.MaxBytes > 0 {
		for c.size > c.config.MaxBytes {
			if !c.dropTail() {
				break
			}
		}
	}
}

// Set sets or updates a cache item for the given key to the given value,
// and the given value size in bytes.
func (c *Cache) Set(key Key, value interface{}, size int64) {
	c.Lock()
	defer c.Unlock()

	var it *Item
	elem, ok := c.m[key]
	if ok {
		// Element exists, push it to front,
		// and then update value etc.
		it = elem.Value.(*Item)
		c.size -= it.Size
		c.l.MoveToFront(elem)
	} else {
		// Insert new element.
		it = &Item{Key: key}
		elem = c.l.PushFront(it)
		c.m[key] = elem
	}
	// Set or update item content.
	it.Value = value
	it.Size = size
	if it.Size < 0 {
		panic("cache: value has negative size")
	}
	it.ModTime = time.Now()
	it.AccessTime = it.ModTime

	// Update cache size.
	if it.Size >= math.MaxInt64-c.size {
		panic("cache: value size is too big")
	}
	c.size += it.Size

	c.enforceCapacity()
}

// getItemPtr returns a pointer to the item under the given key,
// updates its access time, and moves it to the front of the list.
// If there is no such key, it returns nil, false.
// Cache must be locked.
func (c *Cache) getItemPtr(key Key) (it *Item, ok bool) {
	elem, ok := c.m[key]
	if !ok {
		return nil, false
	}
	it = elem.Value.(*Item)
	// Check for expiration.
	if c.config.Expires > 0 && time.Now().Sub(it.ModTime) > c.config.Expires {
		// Item expired, delete it.
		c.removeElement(elem)
		return nil, false
	}
	// Update access time.
	if c.config.TrackAccessTime {
		it.AccessTime = time.Now()
	}
	c.l.MoveToFront(elem)
	return it, true
}

// Get returns a value of item cached under the given key.
// If there is no such key in the cache, it returns nil, false.
func (c *Cache) Get(key Key) (value interface{}, ok bool) {
	c.Lock()
	defer c.Unlock()
	if it, ok := c.getItemPtr(key); ok {
		return it.Value, true
	}
	return nil, false
}

// GetItem returns a copy of item cached under the given key.
func (c *Cache) GetItem(key Key) (it Item, ok bool) {
	c.Lock()
	defer c.Unlock()
	if it, ok := c.getItemPtr(key); ok {
		return *it, true
	}
	return
}

// OldestItem returns a copy of least recently used item. If cache contains no
// items, the second return value is false. Accessing oldest item via this
// function doesn't change its access time or the order of items in cache.
func (c *Cache) OldestItem() (it Item, ok bool) {
	c.Lock()
	defer c.Unlock()
	if tail := c.l.Back(); tail != nil {
		return *(tail.Value.(*Item)), true
	}
	return
}

// removeElement removes a given list element from cache
// and calls removeHander if any.
// Cache must be locked.
func (c *Cache) removeElement(elem *list.Element) {
	it := c.l.Remove(elem).(*Item)
	delete(c.m, it.Key)
	// Update cached byte size.
	c.size -= it.Size
	// Call handler if it exists.
	if c.config.RemoveHandler != nil {
		c.config.RemoveHandler(*it)
	}
}

// dropTail drops the least accessed item from cache.
// Cache must be locked.
func (c *Cache) dropTail() bool {
	tail := c.l.Back()
	if tail == nil {
		return false
	}
	c.removeElement(tail)
	return true
}

// Remove deletes an item with the given key from cache.
func (c *Cache) Remove(key string) bool {
	c.Lock()
	defer c.Unlock()
	elem, ok := c.m[key]
	if !ok {
		return false
	}
	c.removeElement(elem)
	return true
}

// Len returns the number of items in cache.
func (c *Cache) Len() int {
	c.Lock()
	defer c.Unlock()
	return c.l.Len()
}

// Size returns the size of all values in cache.
func (c *Cache) Size() int64 {
	c.Lock()
	defer c.Unlock()
	return c.size
}

// Reconfigure sets a new cache configuration.
func (c *Cache) Reconfigure(newConfig Config) {
	c.Lock()
	defer c.Unlock()
	c.config = newConfig
	c.enforceCapacity()
}

// Config returns a copy of cache configuration.
func (c *Cache) Config() Config {
	c.Lock()
	defer c.Unlock()
	return c.config
}

// Items returns a slice of copies of all items in cache.
func (c *Cache) Items() []Item {
	c.Lock()
	defer c.Unlock()
	items := make([]Item, 0, c.l.Len())
	for elem := c.l.Front(); elem != nil; elem = elem.Next() {
		items = append(items, *elem.Value.(*Item))
	}
	return items
}

// SetBytes is like Set, but accepts a slice of bytes for value,
// and sets the size to its length.
func (c *Cache) SetBytes(key Key, value []byte) {
	c.Set(key, value, int64(len(value)))
}

// GetBytes is like Get, but returns a bytes slice value.
// It panics if value is not a byte slice.
func (c *Cache) GetBytes(key Key) (value []byte, ok bool) {
	v, ok := c.Get(key)
	if !ok {
		return nil, false
	}
	return v.([]byte), true
}
