package caching

import "time"

func GetCache[T iCacheble](k int, link iCache[T]) (T, bool) {
	return link.Get(k)
}

func SetCache[T iCacheble](k int, object T, duration time.Duration, link iCache[T]) {
	link.Set(k, object, duration)
}

func CheckCacheAndWrite[T iCacheble](k int, object T, link iCache[T]) (retObject T, ok bool) {
	retObject, ok = GetCache[T](k, link)

	if !ok {
		// Запрос к БД
		ok = true

		// Запрос успешный -> запись в кеш
		if ok {
			SetCache(k, object, 0, link)
		}
	}

	return retObject, ok
}
