package caching

import (
	"sync"
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/database"
)

type item[T iCacheble] struct {
	value      []T
	created    time.Time
	expiration int64
}

var UsersCache = Init[database.Users](time.Minute, time.Minute*2)

type Cache[T iCacheble] struct {
	mu                sync.RWMutex
	defaultExpiration time.Duration // Время жизни
	cleanupInterval   time.Duration // Интервал очистки
	items             map[int]item[T]
}

func Init[T iCacheble](defaultExpiration, cleanupInterval time.Duration) *Cache[T] {
	cache := Cache[T]{
		items:             make(map[int]item[T]),
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}
	if cleanupInterval > 0 {
		cache.startCleaner()
	}

	return &cache
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

// Получение существующей записи
func (uc *Cache[T]) Get(k int) (res []T, ok bool) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	v, ok := uc.items[k]

	if ok {
		// Не бессрочный И Время жизни не вышло ИЛИ Бессрочный
		if v.expiration > 0 && time.Now().UnixNano() < v.expiration || v.expiration == 0 {
			res = v.value
		} else {
			ok = false
		}
	}

	return res, ok
}

func (uc *Cache[T]) GetByIdx(k int, idx int) (res T, ok bool) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	v, ok := uc.items[k]

	if ok {
		// Не бессрочный И Время жизни не вышло ИЛИ Бессрочный
		if v.expiration > 0 && time.Now().UnixNano() < v.expiration || v.expiration == 0 {
			res = v.value[idx]
		} else {
			ok = false
		}
	}

	return res, ok
}

// Добавление новой записи + перезапись существующей
func (uc *Cache[T]) Set(k int, val T, duration time.Duration) {
	var expr int64

	if duration == 0 {
		duration = uc.defaultExpiration
	}
	if duration > 0 {
		// Время истечения кеша
		expr = time.Now().Add(duration).UnixNano()
	}

	uc.mu.Lock()
	defer uc.mu.Unlock()

	uc.items[k] = item[T]{
		value:      []T{val},
		expiration: expr,
		created:    time.Now(),
	}
}

// Добавление записи в мапу к существующей
func (uc *Cache[T]) Add(k int, val T) {
	uc.mu.RLock()
	item, ok := uc.items[k]
	uc.mu.RUnlock()
	if ok {
		uc.mu.Lock()
		defer uc.mu.Unlock()
		item.value = append(item.value, val)
		uc.items[k] = item
	}
}

// Удаление всей записи
func (uc *Cache[T]) Delete(k int) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	if _, ok := uc.items[k]; ok {
		delete(uc.items, k)
	}
}

// Удаление последней записи из слайса в мапе
func (uc *Cache[T]) Pop(k int) {
	uc.mu.RLock()
	item, ok := uc.items[k]
	uc.mu.RUnlock()
	if ok {
		if len(item.value) > 0 {
			uc.mu.Lock()
			defer uc.mu.Unlock()
			item.value = item.value[:len(item.value)-1]
			uc.items[k] = item
		}
	}
}

// Удаление конкретного элемента из слайса в мапе
func (uc *Cache[T]) DropByIdx(k int, idx int) {
	uc.mu.RLock()
	item, ok := uc.items[k]
	uc.mu.RUnlock()
	if ok {
		if len(item.value) > 0 && len(item.value) > idx {
			uc.mu.Lock()
			defer uc.mu.Unlock()
			item.value = append(item.value[:idx], item.value[idx+1:]...)
			uc.items[k] = item
		}
	}
}

func (uc *Cache[T]) startCleaner() {
	go uc.cleaner()
}

func (uc *Cache[T]) cleaner() {
	for {
		// ожидание
		<-time.After(uc.cleanupInterval)
		if uc.items == nil {
			return
		}

		expiredKeys := func() (keys []int) {
			uc.mu.RLock()
			defer uc.mu.RUnlock()

			for k, i := range uc.items {
				if time.Now().UnixNano() > i.expiration && i.expiration > 0 {
					keys = append(keys, k)
				}
			}

			return keys
		}
		clearItems := func(keys []int) {
			uc.mu.Lock()
			defer uc.mu.Unlock()
			for _, k := range keys {
				delete(uc.items, k)
			}
		}
		if keys := expiredKeys(); len(keys) != 0 {
			clearItems(keys)
		}
	}
}
