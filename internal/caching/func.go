package caching

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/models"
)

func GetCache[T iCacheble](link iCacher[T], k int) ([]T, error) {
	res, ok := link.Get(k)
	if !ok {
		return nil, fmt.Errorf("GetCache:")
	}
	return res, nil
}

func SetCache[T iCacheble](link iCacher[T], k int, object T, duration time.Duration) {
	link.Set(k, object, duration)
}

// Проверка и получение первого объекта
func CheckCacheAndWrite[T iCacheble](link iCacher[T], k int, object T) (retObject T, err error) {
	// Первая проверка, если в кеше есть - возращаем обьект
	// Иначе, проверяем в БД
	retObjectList, err := GetCache(link, k)
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
	SetCache(link, k, data, 0)

	// Считываем повторно из кеша
	retObjectList, err = GetCache(link, k)
	if err != nil {
		return retObject, fmt.Errorf("CheckCacheAndWrite:" + err.Error())
	}
	retObject = retObjectList[0]
	// Нужна ли проверка на консистентность?

	return retObject, err
}

// Функция наполнения кеша из БД
func FillCache[T iCacheble](link iCacher[T], records int, offset int) error {

	structType := &Item[T]{}
	structType.value = make([]T, 1)
	object := &structType.value[0]

	if records == 0 {
		records = 100
	}

	// Определение PK в структуре, по которому будет произведен поиск
	primaryKey, err := models.GetStructInfoPK(object)
	if err != nil {
		return fmt.Errorf("CheckCache:" + err.Error())
	}

	expLst := []database.Expressions{
		{Key: primaryKey.StructNameFields, Operator: database.NotEQ, Value: "0"},
	}
	rs, _, err := database.ReadData(object, expLst, records)
	if err != nil {
		return fmt.Errorf("CheckCache:" + err.Error())
	}

	for k, v := range rs {
		link.Set(k, *v, 0)
	}

	return nil
}

// Возврат количества записей в мапе
func GetCacheCountRecord[T iCacheble](link iCacher[T]) int {
	return link.GetCacheCountRecord()
}

// Возврат ключа по индексу
func GetCacheKeyByIdx[T iCacheble](link iCacher[T], key int) int {
	return link.GetKeyByIdx(key)
}

// Возврат idx элемента слайса из мапы по ключу k
func GetCacheByIdxInMap[T iCacheble](link iCacher[T], k int, idx int) (res T, err error) {
	object, ok := link.GetByIdxInMap(k, idx)
	if !ok {
		err = fmt.Errorf("GetCacheByIdxInMap:key:%d:idx:%d:Cache:%t", k, idx, link)
	}
	return object, err
}

// Возврат связки ключей map[FK][]PK
func GetCacheKeyChain[T iCacheble](link iCacher[T], in interface{}) []interface{} {
	return link.GetKeyChain(in)
}

func GetCacheElementKeyChain[T iCacheble](link iCacher[T], in interface{}) interface{} {
	return link.GetKeyChain(in)[0]
}

// Возврат записей по связке ключей map[FK][]PK с возможностью сортировки
func GetCacheRecordsKeyChain[T iCacheble](link iCacher[T], in interface{}, sorting bool) (out []T, err error) {
	// Сортируем полученные ключи
	keyChain := GetCacheKeyChain(link, in)
	sort.Slice(keyChain, func(i, j int) bool {
		iElem := keyChain[i]
		jElem := keyChain[j]
		switch a := iElem.(type) {
		case int:
			if b, ok := jElem.(int); ok {
				if sorting {
					return a > b
				} else {
					return a < b
				}
			}
		case string:
			if b, ok := jElem.(string); ok {
				if sorting {
					return a > b
				} else {
					return a < b
				}
			}
		}
		return false
	})
	for _, v := range keyChain {
		convV, ok := v.(int)
		if !ok {
			return out, fmt.Errorf("GetCacheRecordsKeyChain:TypeConversionError")
		}
		tracking, err := GetCacheByIdxInMap(link, convV, 0)
		if err != nil {
			return out, fmt.Errorf("GetCacheRecordsKeyChain:" + err.Error())
		} else {
			out = append(out, tracking)
		}
	}
	return out, err
}

// Возврат 10 элементов мапы с offset отсортированному по ключу
func GetCacheOffset[T iCacheble](link iCacher[T], offset int) (out []T, last bool, err error) {
	countRecord := GetCacheCountRecord(link)
	if offset < 10 {
		return nil, false, fmt.Errorf("GetCacheOffset:Offset is small")
	} else if offset >= countRecord {
		offset -= (offset - countRecord)
		last = true
	}

	if countRecord > 1 {
		for i := offset - 10; i < offset; i++ {
			key := GetCacheKeyByIdx(link, i)
			object, err := GetCacheByIdxInMap(link, key, 0)
			if err != nil {
				return out, last, fmt.Errorf("GetCacheOffset:" + err.Error())
			}
			out = append(out, object)
		}
	} else {
		return out, last, fmt.Errorf("GetCacheOffset:Len cache is zero")
	}

	return out, last, err
}
