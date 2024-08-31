package caching

import (
	"fmt"
	"sort"
	"strconv"
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

func UpdateCache[T iCacheble](link iCacher[T], k int, object T) {
	link.Update(k, object)
}

func DropAllCache[T iCacheble](link iCacher[T]) {
	link.DropAll()
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
		return retObject, fmt.Errorf("CheckCacheAndWrite:" + err.Error())
	}

	// Проверка наличия в БД
	result, err := database.CheckRecord[T](buffer)
	if err != nil {
		if strings.Contains(err.Error(), "NoRows") {
			// Запись в БД и возврат ответного тела
			result, _, err = database.WriteRecord[T](buffer)
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
func FillCache[T iCacheble](link iCacher[T], records int, offset ...int) error {

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
		SetCache(link, k, *v, 0)
		// link.Set(k, *v, 0)
	}

	return nil
}

// Возврат количества записей в мапе
func GetCacheCountRecord[T iCacheble](link iCacher[T]) int {
	return link.GetCacheCountRecord()
}

// Возврат количества записей в мапе по доп ключу
func GetCacheSortCountRecord[T iCacheble](link iCacher[T]) int {
	return link.GetCacheSortCountRecord()
}

// Возврат ключа по индексу
func GetCacheKeyByIdx[T iCacheble](link iCacher[T], sort string, key int) int {
	return link.GetKeyByIdx(sort, key)
}

// Возврат idx элемента слайса из мапы по ключу k
// idx по умолчанию равен 0 елементу
func GetCacheByIdxInMap[T iCacheble](link iCacher[T], k int, idx ...int) (res T, err error) {
	var idxE int = 0
	for _, v := range idx {
		idxE = v
	}
	object, ok := link.GetByIdxInMap(k, idxE)
	if !ok {
		err = fmt.Errorf("GetCacheByIdxInMap:key:%d:idx:%d:Cache:%t", k, idxE, link)
	}
	return object, err
}

// Возврат связки ключей map[FK][]PK
func GetCacheKeyChain[T iCacheble](link iCacher[T], in interface{}) []interface{} {
	return link.GetKeyChain(in)
}

func GetCacheElementKeyChain[T iCacheble](link iCacher[T], in interface{}) (out interface{}) {
	keyChain := GetCacheKeyChain(link, in)
	if keyChain != nil {
		out = keyChain[0]
	}
	return out
}

// Возврат связки ключей map[FK][]PK
func GetCacheKeyChainSort[T iCacheble](link iCacher[T], in interface{}) []interface{} {
	return link.GetKeyChainSort(in)
}

func GetCacheElementKeyChainSort[T iCacheble](link iCacher[T], in interface{}) (out interface{}) {
	keyChain := GetCacheKeyChainSort(link, in)
	if keyChain != nil {
		out = keyChain[0]
	}
	return out
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

// Возврат n элементов мапы с offset отсортированному по стандартному ключу
func GetCacheOffset[T iCacheble](link iCacher[T], offset int, recordCnt ...int) (out []T, last bool, err error) {
	// Стандартно возвращаем по 10 записей
	var recordCntV int = 10
	for _, v := range recordCnt {
		recordCntV = v
		break
	}

	countRecord := GetCacheCountRecord(link)
	if offset < recordCntV {
		return nil, false, fmt.Errorf("GetCacheOffset:Offset is small")
	} else if offset >= countRecord {
		offset -= (offset - countRecord)
		last = true
	}

	structType := &Item[T]{}
	structType.value = make([]T, 1)
	object := &structType.value[0]
	primaryKey, err := models.GetStructInfoPK(object)

	if countRecord > 1 {
		for i := offset - recordCntV; i < offset; i++ {
			key := GetCacheKeyByIdx(link, primaryKey.StructNameFields, i)
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

// Возврат всех элементов мапы отсортированному по стандартному ключу
func GetCacheAllRecord[T iCacheble](link iCacher[T]) (out []T, err error) {
	countRecord := GetCacheCountRecord(link)

	structType := &Item[T]{}
	structType.value = make([]T, 1)
	object := &structType.value[0]
	primaryKey, err := models.GetStructInfoPK(object)

	if countRecord > 1 {
		for i := 0; i < countRecord; i++ {
			key := GetCacheKeyByIdx(link, primaryKey.StructNameFields, i)
			object, err := GetCacheByIdxInMap(link, key, 0)
			if err != nil {
				return out, fmt.Errorf("GetCacheAllRecord:" + err.Error())
			}
			out = append(out, object)
		}
	} else {
		return out, fmt.Errorf("GetCacheAllRecord:Len cache is zero")
	}

	return out, err

}

// Возврат n элементов мапы с offset отсортированному по дополнительному ключу
func GetCacheOffsetSort[T iCacheble](link iCacher[T], offset int, recordCnt ...int) (out []T, last bool, err error) {
	// Стандартно возвращаем по 10 записей
	var recordCntV int = 10
	for _, v := range recordCnt {
		recordCntV = v
		break
	}

	countRecord := GetCacheSortCountRecord(link)
	if offset < recordCntV {
		return nil, false, fmt.Errorf("GetCacheOffsetSort:Offset is small")
	} else if offset >= countRecord {
		offset -= (offset - countRecord)
		last = true
	}

	structType := &Item[T]{}
	structType.value = make([]T, 1)
	object := &structType.value[0]
	sortKey, err := models.GetStructInfoSort(object)

	if countRecord > 1 {
		for i := offset - recordCntV; i < offset; i++ {
			rank := GetCacheKeyByIdx(link, sortKey.StructNameFields, i)
			keyChainSort := GetCacheElementKeyChainSort(link, rank)
			object, err := GetCacheByIdxInMap(link, keyChainSort.(int), 0)
			if err != nil {
				return out, last, fmt.Errorf("GetCacheOffsetSort:" + err.Error())
			}
			out = append(out, object)
		}
	} else {
		return out, last, fmt.Errorf("GetCacheOffsetSort:Len cache is zero")
	}

	return out, last, err
}

// Обновление записи
func UpdateCacheRecord[T iCacheble](link iCacher[T], k int, object T, cacheOn ...bool) (retObject T, err error) {
	var notFound bool
	var id interface{}
	cache := true
	for _, item := range cacheOn {
		cache = item
		break
	}
	// Проверка на существование объекта
	retObject, err = GetCacheByIdxInMap(link, k)
	if err != nil {
		notFound = true
		// return retObject, fmt.Errorf("UpdateCacheRecord:" + err.Error())
	}

	// Сериализация для отправки
	buffer, err := models.MarshalJSON(object)
	if err != nil {
		return retObject, fmt.Errorf("UpdateCacheRecord:" + err.Error())
	}

	// Проверка наличия в БД
	result, err := database.CheckRecord[T](buffer)
	if err != nil {
		if strings.Contains(err.Error(), "NoRows") {
			// Запись в БД и возврат ответного тела
			result, id, err = database.WriteRecord[T](buffer)
			if err != nil {
				return retObject, fmt.Errorf("UpdateCacheRecord:" + err.Error())
			}
			// // Считывание добавленной записи из БД
			// primaryKey, err := models.GetStructInfoPK(object)
			// if err != nil {
			// 	return retObject, fmt.Errorf("UpdateCacheRecord:" + err.Error())
			// }
			// expLst := []database.Expressions{
			// 	{Key: primaryKey.StructNameFields, Operator: database.EQ, Value: func(interface{}) (out string) {
			// 		if n, ok := id.(int64); ok {
			// 			out = strconv.Itoa(int(n))
			// 		}
			// 		return out
			// 	}(id)},
			// }
			// structType := &Item[T]{}
			// structType.value = make([]T, 1)
			// objectV := &structType.value[0]
			// rs, _, err = database.ReadData(objectV, expLst, 1)
			// if err != nil {
			// 	return retObject, fmt.Errorf("UpdateCacheRecord:" + err.Error())
			// }
		} else {
			return retObject, fmt.Errorf("UpdateCacheRecord:" + err.Error())
		}
	} else {
		// Обновление в БД и возврат ответного тела
		result, err = database.UpdateRecord[T](buffer)
		if err != nil {
			return retObject, fmt.Errorf("UpdateCacheRecord:" + err.Error())
		}
	}

	// Десереализация для обновления в кеше
	data, err := models.UnmarshalJSON[T](result)
	if err != nil {
		return retObject, fmt.Errorf("UpdateCacheRecord:" + err.Error())
	}

	if cache {
		if !notFound {
			// Обновление в кеше
			UpdateCache(link, k, data)
		} else {
			switch idConv := id.(type) {
			case interface{}:
				switch idInt := idConv.(type) {
				case int64:
					SetCache(link, int(idInt), data, 0)
					id = idInt
				default:
					SetCache(link, k, data, 0)
				}
			}
		}
	}

	// Считываем повторно из кеша
	retObject, err = GetCacheByIdxInMap(link, k)
	if err != nil {
		return retObject, fmt.Errorf("UpdateCacheRecord:" + err.Error())
	}

	// Нужна ли проверка на консистентность?

	return retObject, err
}

// Запись в кеш с записью в БД без проверок на существование
func WriteCache[T iCacheble](link iCacher[T], k int, object T, cacheOn ...bool) (retObject T, id int64, err error) {
	cache := true
	for _, item := range cacheOn {
		cache = item
		break
	}
	// Сериализация для отправки
	buffer, err := models.MarshalJSON(object)
	if err != nil {
		return retObject, -1, fmt.Errorf("WriteRecord:" + err.Error())
	}

	// Запись в БД и возврат ответного тела
	result, idw, err := database.WriteRecord[T](buffer)
	if err != nil {
		return retObject, -1, fmt.Errorf("CheckCacheAndWrite:" + err.Error())
	}

	// Десереализация для записи в кеш
	data, err := models.UnmarshalJSON[T](result)
	if err != nil {
		return retObject, -1, fmt.Errorf("CheckCacheAndWrite:" + err.Error())
	}

	// Считывание добавленной записи из БД
	primaryKey, err := models.GetStructInfoPK(object)
	if err != nil {
		return retObject, -1, fmt.Errorf("CheckCacheAndWrite:" + err.Error())
	}
	expLst := []database.Expressions{
		{Key: primaryKey.StructNameFields, Operator: database.EQ, Value: func(interface{}) (out string) {
			if n, ok := idw.(int64); ok {
				out = strconv.Itoa(int(n))
			}
			return out
		}(idw)},
	}
	structType := &Item[T]{}
	structType.value = make([]T, 1)
	objectV := &structType.value[0]
	rs, _, err := database.ReadData(objectV, expLst, 1)
	if err != nil {
		return retObject, -1, fmt.Errorf("CheckCacheAndWrite:" + err.Error())
	}

	// Запись к кеш
	if cache {
		for k, v := range rs {
			switch idConv := idw.(type) {
			case interface{}:
				switch idInt := idConv.(type) {
				case int64:
					SetCache(link, int(idInt), *v, 0)
					id = idInt
				default:
					SetCache(link, k, *v, 0)
				}
			}
		}
	}
	// SetCache(link, int(id), data, 0)

	return data, id, err
}
