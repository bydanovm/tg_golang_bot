package caching

func GetCache[T any](k int, link iCache[T]) (T, bool) {
	return link.Get(k)
}

func SetCache[T any](k int, object T, link iCache[T]) {
	link.Set(k, object)
}

func CheckCacheAndWrite[T any](k int, object T, link iCache[T]) (retObject T, ok bool) {
	if _, ok := GetCache[T](k, link); !ok {
		// Запрос к БД
		ok = true

		// Запрос успешный -> запись в кеш
		if ok {
			SetCache(k, object, link)
		}
	}

	return retObject, ok
}
