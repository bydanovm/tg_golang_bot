package caching

import (
	"sort"
	"sync"
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/models"
)

var UsersCache = Init[database.Users](time.Hour, time.Minute*30)

var CryptoCache = Init[database.DictCrypto](time.Hour*24, time.Hour*12)

var TrackingCache = Init[database.TrackingCrypto](time.Hour*24, time.Hour*12)

var TrackingTypeCache = Init[database.TypeTrackingCrypto](time.Hour*24*365, 0)

var LimitsCache = Init[database.Limits](time.Hour*24, time.Hour*12)

var LimitsDictCache = Init[database.LimitsDict](time.Hour*24*365, 0)

// Кеш коин маркетов
var CoinMarketsCache = Init[database.CoinMarkets](time.Hour*24*365, 0)
var CoinMarketsEndpointCache = Init[database.CoinMarketsEndpoint](time.Hour*24*365, 0)
var CoinMarketsHandCache = Init[database.CoinMarketsHand](time.Hour*24*365, 0)

// Временный кеш с ценами КВ
var CryptoPricesCache = Init[database.Cryptoprices](0, 0)

type Item[T iCacheble] struct {
	value      T
	created    time.Time
	expiration int64
}

type Cache[T iCacheble] struct {
	mu                sync.RWMutex
	defaultExpiration time.Duration // Время жизни
	cleanupInterval   time.Duration // Интервал очистки
	items             map[int]Item[T]
	keysSort          map[string][]int
	keysMap           []map[interface{}][]interface{}
}

// Инициализация кеша
func Init[T iCacheble](defaultExpiration, cleanupInterval time.Duration) *Cache[T] {
	cache := Cache[T]{
		items:             make(map[int]Item[T]),
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		keysSort:          make(map[string][]int),
		keysMap:           make([]map[interface{}][]interface{}, 3),
		// keysMap:           make([]map[interface{}][]interface{}, 3),
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
func (uc *Cache[T]) Get(k int) (res T, ok bool) {
	isRlock := uc.URLockU()
	v, ok := uc.items[k]

	if ok {
		// Не бессрочный И Время жизни не вышло ИЛИ Бессрочный
		if v.expiration > 0 && time.Now().UnixNano() < v.expiration || v.expiration == 0 {
			res = v.value
			// Обновить время жизни
			v.expiration = time.Now().Add(uc.defaultExpiration).UnixNano()
			isRlock = uc.URUnlock()
			uc.mu.Lock()
			uc.items[k] = v
			uc.mu.Unlock()
		} else {
			ok = false
		}
	}

	if isRlock {
		uc.mu.RUnlock()
	}

	return res, ok
}

func (uc *Cache[T]) GetByIdxInMap(k int, idx int) (res T, ok bool) {
	isRlock := uc.URLockU()
	v, ok := uc.items[k]

	if ok {
		// Не бессрочный И Время жизни не вышло ИЛИ Бессрочный
		if v.expiration > 0 && time.Now().UnixNano() < v.expiration || v.expiration == 0 {
			res = v.value
			// Обновить время жизни
			v.expiration = time.Now().Add(uc.defaultExpiration).UnixNano()
			isRlock = uc.URUnlock()
			uc.mu.Lock()
			uc.items[k] = v
			uc.mu.Unlock()
		} else {
			ok = false
		}
	}

	if isRlock {
		uc.mu.RUnlock()
	}

	return res, ok
}
func (uc *Cache[T]) GetKeyByIdx(sort string, idx int) (key int) {
	if v, ok := uc.keysSort[sort]; ok {
		if len(v) > idx {
			key = v[idx]
		}
	}
	return uc.keysSort[sort][idx]
}

// Возврат связки ключей map[FK][]PK
func (uc *Cache[T]) GetKeyChain(in interface{}) []interface{} {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	v, ok := uc.keysMap[0][in]
	if !ok {
		return nil
	}

	return v
}

// Возврат связки ключей map[FK][]PK для доп сортировки
func (uc *Cache[T]) GetKeyChainSort(in interface{}) []interface{} {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	v, ok := uc.keysMap[1][in]
	if !ok {
		return nil
	}

	return v
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
	// Проверка на нахождение ключа в слайсе
	uc.mu.RLock()
	_, ok := uc.items[k]
	uc.mu.RUnlock()
	primaryKey, _ := models.GetStructInfoPK(val)
	primaryKey.StructValue = k
	uc.mu.Lock()
	defer uc.mu.Unlock()
	if !ok {
		uc.keysSort[primaryKey.StructNameFields] = append(uc.keysSort[primaryKey.StructNameFields], k)
		sort.Slice(uc.keysSort[primaryKey.StructNameFields], func(i, j int) bool {
			return uc.keysSort[primaryKey.StructNameFields][i] < uc.keysSort[primaryKey.StructNameFields][j]
		})
		sortKey, err := models.GetStructInfoSort(val)
		if err == nil {
			if n, ok := sortKey.StructValue.(int); ok && n > 0 {
				uc.keysSort[sortKey.StructNameFields] = append(uc.keysSort[sortKey.StructNameFields], n)
				sort.Slice(uc.keysSort[sortKey.StructNameFields], func(i, j int) bool {
					return uc.keysSort[sortKey.StructNameFields][i] < uc.keysSort[sortKey.StructNameFields][j]
				})
			}
		}
	}

	uc.items[k] = Item[T]{
		value:      val,
		expiration: expr,
		created:    time.Now(),
	}

	// Добавление ключа связки PK <-> FK
	foreignKey, err := models.GetStructInfoFK(val)
	if err == nil {
		// map[interface{}][]interface{}
		if uc.keysMap[0] == nil {
			uc.keysMap[0] = make(map[interface{}][]interface{})
		}
		uc.keysMap[0][foreignKey.StructValue] = append(uc.keysMap[0][foreignKey.StructValue], primaryKey.StructValue)
	}
	// Добавление доп сортировки
	sortKey, err := models.GetStructInfoSort(val)
	if err == nil {
		if uc.keysMap[1] == nil {
			uc.keysMap[1] = make(map[interface{}][]interface{})
		}
		uc.keysMap[1][sortKey.StructValue] = append(uc.keysMap[1][sortKey.StructValue], primaryKey.StructValue)
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
		// item.value = append(item.value, val)
		item.value = val
		uc.items[k] = item
	}
}

// Обновление записи в кеше
func (uc *Cache[T]) Update(k int, val T) {
	uc.mu.RLock()
	item, ok := uc.items[k]
	uc.mu.RUnlock()
	if ok {
		uc.mu.Lock()
		defer uc.mu.Unlock()
		item.expiration = time.Now().Add(uc.defaultExpiration).UnixNano()
		// valT := []T{}
		// valT = append(valT, val)
		item.value = val
		uc.items[k] = item
	}
}

// Удаление всей записи
func (uc *Cache[T]) Delete(k int) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	if v, ok := uc.items[k]; ok {
		primaryKey, _ := models.GetStructInfoPK(v.value)
		delete(uc.items, k)
		// Удаляем ключ из слайса
		if v1, ok1 := uc.keysSort[primaryKey.StructNameFields]; ok1 {
			for idx, val := range v1 {
				if val == k {
					uc.keysSort[primaryKey.StructNameFields] = append(uc.keysSort[primaryKey.StructNameFields][:idx], uc.keysSort[primaryKey.StructNameFields][idx+1:]...)
					break
				}
			}
		}

	}
}

// Удаление последней записи из слайса в мапе
// func (uc *Cache[T]) Pop(k int) {
// 	uc.mu.RLock()
// 	item, ok := uc.items[k]
// 	uc.mu.RUnlock()
// 	if ok {
// 		// if len(item.value) > 0 {
// 			uc.mu.Lock()
// 			defer uc.mu.Unlock()
// 			item.value = item.value[:len(item.value)-1]
// 			uc.items[k] = item
// 		// }
// 	}
// }

// Удаление конкретного элемента из слайса в мапе
// func (uc *Cache[T]) DropByIdx(k int, idx int) {
// 	uc.mu.RLock()
// 	item, ok := uc.items[k]
// 	uc.mu.RUnlock()
// 	if ok {
// 		if len(item.value) > 0 && len(item.value) > idx {
// 			uc.mu.Lock()
// 			defer uc.mu.Unlock()
// 			item.value = append(item.value[:idx], item.value[idx+1:]...)
// 			uc.items[k] = item
// 		}
// 	}
// }

func (uc *Cache[T]) DropAll() {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	uc.items = make(map[int]Item[T])
	uc.keysMap = make([]map[interface{}][]interface{}, 3)
	uc.keysSort = make(map[string][]int)
}

func (uc *Cache[T]) GetCacheCountRecord() int {
	structType := &Item[T]{}
	// structType.value = make([]T, 1)
	object := &structType.value
	primaryKey, _ := models.GetStructInfoPK(object)
	return len(uc.keysSort[primaryKey.StructNameFields])
}

func (uc *Cache[T]) GetCacheSortCountRecord() int {
	structType := &Item[T]{}
	// structType.value = make([]T, 1)
	object := &structType.value
	sortKey, _ := models.GetStructInfoSort(object)
	return len(uc.keysSort[sortKey.StructNameFields])
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
			for _, k := range keys {
				uc.Delete(k)
			}

		}
		if keys := expiredKeys(); len(keys) != 0 {
			clearItems(keys)
		}
	}
}
