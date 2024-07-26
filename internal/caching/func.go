package caching

import (
	"fmt"
	"strings"
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/models"
)

func GetCache[T iCacheble](k int, link iCacher[T]) ([]T, error) {
	res, ok := link.Get(k)
	if !ok {
		return nil, fmt.Errorf("GetCache:")
	}
	return res, nil
}

func SetCache[T iCacheble](k int, object T, duration time.Duration, link iCacher[T]) {
	link.Set(k, object, duration)
}

// Проверка и получение первого объекта
func CheckCacheAndWrite[T iCacheble](k int, object T, link iCacher[T]) (retObject T, err error) {
	// Первая проверка, если в кеше есть - возращаем обьект
	// Иначе, проверяем в БД
	retObjectList, err := GetCache(k, link)
	if err == nil {
		retObject = retObjectList[0]
		return retObject, nil
	}

	// Сериализация для отправки
	buffer, err := models.MarshalJSON(object)
	if err != nil {
		return retObject, fmt.Errorf("WriteRecord:" + err.Error())
	}

	// Проверка наличия в БД
	result, err := database.CheckRecord[T](buffer)
	if err != nil {
		if strings.Contains(err.Error(), "NoRows") {
			// Запись в БД и возврат ответного тела
			result, err = database.WriteRecord[T](buffer)
			if err != nil {
				return retObject, fmt.Errorf("CheckCacheAndWrite:" + err.Error())
			}
		} else {
			return retObject, fmt.Errorf("CheckCacheAndWrite:" + err.Error())
		}
	}

	// Десереализация для записи в кеш
	data, err := models.UnmarshalJSON[T](result)
	if err != nil {
		return retObject, fmt.Errorf("CheckCacheAndWrite:" + err.Error())
	}

	// Запись к кеш
	SetCache(k, data, 0, link)

	// Считываем повторно из кеша
	retObjectList, err = GetCache(k, link)
	if err != nil {
		return retObject, fmt.Errorf("CheckCacheAndWrite:" + err.Error())
	}
	retObject = retObjectList[0]
	// Нужна ли проверка на консистентность?

	return retObject, err
}
