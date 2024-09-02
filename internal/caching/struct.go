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

// Кеш структур данных (кеш в кеше)
var structInfoCache = Init[models.StructInfo](0, 0)

type Item[T iCacheble] struct {
	value      T
	created    time.Time
	expiration int64
}

type Cache[T iCacheble] struct {
	mu                sync.RWMutex
	defaultExpiration time.Duration // Время жизни
	cleanupInterval   time.Duration // Интервал очистки
	items             map[interface{}]Item[T]
	keysSort          map[string][]int
	keysMap           []map[interface{}][]interface{}
}

// Инициализация кеша
func Init[T iCacheble](defaultExpiration, cleanupInterval time.Duration) *Cache[T] {
	cache := Cache[T]{
		items:             make(map[interface{}]Item[T]),
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
func (uc *Cache[T]) Get(k interface{}) (res T, ok bool) {
	isRlock := uc.URLockU()
	v, ok := uc.items[k]

	if ok {
		// Не бессрочный И Время жизни не вышло ИЛИ Бессрочный
		if v.expiration > 0 && time.Now().UnixNano() < v.expiration || v.expiration == 0 {
			res = v.value
			// Обновить время жизни
			if v.expiration != 0 {
				v.expiration = time.Now().Add(uc.defaultExpiration).UnixNano()
				isRlock = uc.URUnlock()
				uc.mu.Lock()
				uc.items[k] = v
				uc.mu.Unlock()
			}
		} else {
			ok = false
		}
	}

	if isRlock {
		uc.mu.RUnlock()
	}

	return res, ok
}

func (uc *Cache[T]) GetByIdxInMap(k interface{}, idx int) (res T, ok bool) {
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

// Добавление к кеш без сортировок
func (uc *Cache[T]) SetLite(k interface{}, val T, duration time.Duration) {
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

	uc.items[k] = Item[T]{
		value:      val,
		expiration: expr,
		created:    time.Now(),
	}
}

// Добавление новой записи + перезапись существующей
func (uc *Cache[T]) Set(k interface{}, val T, duration time.Duration) {
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

	structInfo := models.StructInfo{}
	structInfo.GetFieldInfo(val)

	uc.mu.Lock()
	defer uc.mu.Unlock()
	if !ok {
		uc.keysSort[structInfo.StructPKKey] = append(uc.keysSort[structInfo.StructPKKey], k.(int))
		sort.Slice(uc.keysSort[structInfo.StructPKKey], func(i, j int) bool {
			return uc.keysSort[structInfo.StructPKKey][i] < uc.keysSort[structInfo.StructPKKey][j]
		})
		if structInfo.StructSortKey != "" {
			if n, ok := structInfo.StructFieldInfo[structInfo.StructSortKey].StructValue.(int); ok && n > 0 {
				uc.keysSort[structInfo.StructSortKey] = append(uc.keysSort[structInfo.StructSortKey], n)
				sort.Slice(uc.keysSort[structInfo.StructSortKey], func(i, j int) bool {
					return uc.keysSort[structInfo.StructSortKey][i] < uc.keysSort[structInfo.StructSortKey][j]
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
	if structInfo.StructFKKey != "" {
		if uc.keysMap[0] == nil {
			uc.keysMap[0] = make(map[interface{}][]interface{})
		}
		value := structInfo.StructFieldInfo[structInfo.StructFKKey].StructValue
		uc.keysMap[0][value] = append(uc.keysMap[0][value], k.(int))
	}
	// Добавление доп сортировки
	if structInfo.StructSortKey != "" {
		if uc.keysMap[1] == nil {
			uc.keysMap[1] = make(map[interface{}][]interface{})
		}
		value := structInfo.StructFieldInfo[structInfo.StructSortKey].StructValue
		uc.keysMap[1][value] = append(uc.keysMap[1][value], k.(int))
	}

}

// Добавление записи в мапу к существующей
func (uc *Cache[T]) Add(k interface{}, val T) {
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
func (uc *Cache[T]) Update(k interface{}, val T) {
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
func (uc *Cache[T]) Delete(k interface{}) {
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
	uc.items = make(map[interface{}]Item[T])
	uc.keysMap = make([]map[interface{}][]interface{}, 3)
	uc.keysSort = make(map[string][]int)
}

func (uc *Cache[T]) GetCacheCountRecord() int {
	structInfo, ok := structInfoCache.Get(models.GetName(Item[T]{}.value))
	if !ok {
		return 0
	}
	return len(uc.keysSort[structInfo.StructPKKey])
}

func (uc *Cache[T]) GetCacheSortCountRecord() int {
	structInfo, ok := structInfoCache.Get(models.GetName(Item[T]{}.value))
	if !ok {
		return 0
	}
	return len(uc.keysSort[structInfo.StructSortKey])
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
					keys = append(keys, k.(int))
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
