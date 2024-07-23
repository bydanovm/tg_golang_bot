package caching

import "time"

func GetCache[T iCacheble](k int, link iCacher[T]) ([]T, bool) {
	return link.Get(k)
}

func SetCache[T iCacheble](k int, object T, duration time.Duration, link iCacher[T]) {
	link.Set(k, object, duration)
}

func CheckCacheAndWrite[T iCacheble](k int, object T, link iCacher[T]) (retObject []T, ok bool) {
	retObject, ok = GetCache[T](k, link)

	if !ok {
		// Запрос к БД
		ok = true

		// Запрос успешный -> запись в кеш
		if ok {
			SetCache(k, object, 0, link)
		}

		retObject, ok = GetCache[T](k, link)
	}

	return retObject, ok
}
