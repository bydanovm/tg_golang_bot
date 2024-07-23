package caching

import (
	"sync"
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/database"
)

type item[T any] struct {
	value      map[int]T
	created    time.Time
	expiration int64
}

var UsersCache = Init[database.Users](time.Minute, time.Minute*2)

type Cache[T any] struct {
	mu                sync.RWMutex
	defaultExpiration time.Duration // Время жизни
	cleanupInterval   time.Duration // Интервал очистки
	items             item[T]
}

func Init[T any](defaultExpiration, cleanupInterval time.Duration) *Cache[T] {
	return &Cache[T]{
		items: item[T]{
			value: make(map[int]T),
		},
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}
	// return &Cache[T]{
	// 	items:             item[T]{
	// 		make(map[int]T)
	// 	},
}

func (uc *Cache[T]) URLockU() (isRLock bool) {
	uc.mu.RLock()
	isRLock = true
	return isRLock
}
func (uc *Cache[T]) URUnlock() (isRLock bool) {
	uc.mu.RUnlock()
	isRLock = false
	return isRLock
}

func (uc *Cache[T]) Get(k int) (T, bool) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	v, ok := uc.items.value[k]
	return v, ok
}

func (uc *Cache[T]) Set(k int, val T) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	uc.items.value[k] = val
}

func (uc *Cache[T]) Delete(k int) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	delete(uc.items.value, k)
}
