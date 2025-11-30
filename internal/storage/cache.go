package storage

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Data[T any] struct {
	value  T
	expiry time.Time
}

type Cache[K comparable, T any] struct {
	data map[K]*Data[T]
	mu   sync.RWMutex
}

type Caches struct {
	UserCache  *Cache[string, string]
	TokenCache *Cache[string, bool]
}

func NewCacheStorage() *Caches {
	return &Caches{
		UserCache:  NewCache[string, string](),
		TokenCache: NewCache[string, bool](),
	}
}

func NewCache[K comparable, T any]() *Cache[K, T] {
	return &Cache[K, T]{
		data: make(map[K]*Data[T]),
	}
}

func (c *Cache[K, T]) Set(key K, val T, expiresAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &Data[T]{
		value:  val,
		expiry: expiresAt,
	}
}

func zeroVal[T any]() T {
	var val T
	return val
}

func (c *Cache[K, T]) Get(key K) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, ok := c.data[key]
	if !ok {
		return zeroVal[T](), false
	}

	if data.expiry.Before(time.Now()) {
		delete(c.data, key)
		fmt.Println("ANTES?", data.expiry.Before(time.Now()), "DATA DE EXPIRACAO", data.expiry, "TEMPO ATUAL", time.Now())
		return zeroVal[T](), false
	}

	return data.value, true
}

func (c *Cache[K, T]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *Cache[K, T]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Just to make sure that the memory of each
	// bucket is clear at a certain level... even
	// though it impacts the perfomance
	runtime.GC()

	c.data = make(map[K]*Data[T])
}
