package database

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"
)

//	type CacheTemplate interface {
//		CheckCache(int) error
//		GetCache(int) error
//	}

type RWMutexUsr struct {
	sync.RWMutex
}

type UserCache struct {
	mu   RWMutexUsr
	Item UsersCacheType
}

// Кеширование пользователей
type UsersCacheType map[int]Users

var UsersCache = Init()

func Init() *UserCache {
	items := make(UsersCacheType)

	cache := UserCache{
		Item: items,
	}

	return &cache
}

func (uc *UserCache) CheckCache(idUsr int) (err error) {
	uc.mu.RLock()
	isLock := true
	if _, ok := uc.Item[idUsr]; !ok {
		uc.mu.RUnlock()
		isLock = false
		// Заполняем информацию в кеш из БД
		user := Users{IdUsr: idUsr}
		if err = user.CheckUser(); err != nil {
			err = fmt.Errorf("CheckCache:" + err.Error())
		} else {
			uc.mu.Lock()
			uc.Item[idUsr] = user
			uc.mu.Unlock()
		}
	}
	if isLock {
		uc.mu.RUnlock()
	}
	return err
}

func (uc *UserCache) GetCache(idUsr int) (Users, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if v, ok := uc.Item[idUsr]; !ok {
		return Users{}, fmt.Errorf("GetCache:User not initialised")
	} else {
		return v, nil
	}
}

func (uc *UserCache) GetCount() (cnt int) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	cnt = len(uc.Item)
	return cnt
}

func (uc *UserCache) GetUserId(idUsr int) (id int) {
	if user, err := uc.GetCache(idUsr); err != nil {
		id = 0
	} else {
		id = user.IdUsr
	}
	return id
}

func (uc *UserCache) GetUserName(idUsr int) (name string) {
	if user, err := uc.GetCache(idUsr); err != nil {
		name = "Not found"
	} else {
		name = user.NameUsr
	}
	return name
}

func (uc *UserCache) GetChatId(idUsr int) (chatId int64) {
	if user, err := uc.GetCache(idUsr); err != nil {
		chatId = int64(user.IdUsr)
	} else {
		chatId = user.ChatIdUsr
	}
	return chatId
}

func (uc *UserCache) GetFLName(idUsr int) (FLName string) {
	if user, err := uc.GetCache(idUsr); err != nil {
		FLName = user.FirstName + " " + user.LastName
	}
	return FLName
}

// Кеш типов отслеживаний
type TypeTrackingCryptoCache map[int]TypeTrackingCrypto
type TypeTrackingCryptoCacheKeys map[string]int

var TypeTCCache = make(TypeTrackingCryptoCache)
var TypeTCCacheKeys = make(TypeTrackingCryptoCacheKeys)

func (ttc *TypeTrackingCryptoCache) CheckAllCache() error {
	t := *ttc
	if _, ok := t[1]; !ok {
		// Заполняем информацию в кеш из БД
		typeTC := TypeTrackingCrypto{}
		rs, err := typeTC.GetAllTypeInfo()
		if err != nil {
			return fmt.Errorf("CheckCache:" + err.Error())
		}
		for _, v := range rs {
			t[v.(TypeTrackingCrypto).IdTypTrkCrp] = v.(TypeTrackingCrypto)
			TypeTCCacheKeys[v.(TypeTrackingCrypto).NameTypeTrkCrp] = v.(TypeTrackingCrypto).IdTypTrkCrp
		}
	}
	return nil
}
func (ttc *TypeTrackingCryptoCache) CheckCache(id int) error {
	t := *ttc
	if _, ok := t[id]; !ok {
		// Заполняем информацию в кеш из БД
		typeTC := TypeTrackingCrypto{IdTypTrkCrp: id}
		v, err := typeTC.GetTypeInfo()
		if err != nil {
			return fmt.Errorf("CheckCache:" + err.Error())
		}
		t[v.(TypeTrackingCrypto).IdTypTrkCrp] = v.(TypeTrackingCrypto)
	}
	return nil
}

func (ttc *TypeTrackingCryptoCache) GetCache(idType int) (TypeTrackingCrypto, error) {
	t := *ttc
	if v, ok := t[idType]; !ok {
		return TypeTrackingCrypto{}, fmt.Errorf("GetCache:User not initialised")
	} else {
		return v, nil
	}
}

// Поиск ИД Типа по Имени
func (ttc *TypeTrackingCryptoCacheKeys) GetCacheIdByName(name string) int {
	t := *ttc
	v, ok := t[name]
	if ok {
		return v
	}
	return 0
}

// Кеш активных отслеживаний
type TrackingCryptoCache map[int]TrackingCrypto

var TCCache = make(TrackingCryptoCache)

func (ttc *TrackingCryptoCache) CheckAllCache() error {
	t := *ttc
	if _, ok := t[1]; !ok {
		// Заполняем информацию в кеш из БД
		expLst := []Expressions{
			// {Key: "OnTrkCrp", Operator: EQ, Value: "true"},
			{Key: "IdTrkCrp", Operator: NotEQ, Value: "0"},
		}
		rs, find, _, err := ReadDataRow(&TrackingCrypto{}, expLst, 0)
		if err != nil {
			return fmt.Errorf("CheckCache:" + err.Error())
		}
		if !find {
			// return fmt.Errorf("CheckCache:not find cache")
			return nil
		}
		subFields := TrackingCrypto{}
		for _, subRs := range rs {
			mapstructure.Decode(subRs, &subFields)
			t[subFields.IdTrkCrp] = subFields
		}
	}
	return nil
}

func (ttc *TrackingCryptoCache) GetCache(id int) (TrackingCrypto, error) {
	t := *ttc
	if v, ok := t[id]; !ok {
		return TrackingCrypto{}, fmt.Errorf("GetCache:User not initialised")
	} else {
		return v, nil
	}
}

func (ttc *TrackingCryptoCache) GetCacheLastId() int {
	t := *ttc
	var maxId int
	for k, _ := range t {
		if k > maxId {
			maxId = k
		}
	}
	return maxId + 1
}

// Кеш словаря критовалют
type DictCryptoCache map[int]interface{}
type DictCryptoCacheKeys map[string]int // Словарь symbol - Id

var DCCache = make(DictCryptoCache)
var DCCacheKeys = make(DictCryptoCacheKeys)

func (dcc *DictCryptoCache) CheckAllCache() error {
	d := *dcc
	if _, ok := d[1]; !ok {
		// Заполняем информацию в кеш из БД
		expLst := []Expressions{
			{Key: "CryptoId", Operator: NotEQ, Value: "0"},
		}
		rs, find, _, err := ReadDataRow(&DictCrypto{}, expLst, 0)
		if err != nil {
			return fmt.Errorf("CheckCache:" + err.Error())
		}
		if !find {
			// return fmt.Errorf("CheckCache:not find cache")
			return nil
		}
		subFields := DictCrypto{}
		for _, subRs := range rs {
			mapstructure.Decode(subRs, &subFields)
			d[subFields.CryptoId] = subFields
			DCCacheKeys[subFields.CryptoName] = subFields.CryptoId
		}
	}
	return nil
}
func (dcc *DictCryptoCache) GetCache(id int) (DictCrypto, error) {
	d := *dcc
	if v, ok := d[id].(DictCrypto); !ok {
		return DictCrypto{}, fmt.Errorf("GetCache:Crypto not initialised")
	} else {
		return v, nil
	}
}
func (dcc *DictCryptoCache) GetAllCache() (DCout []DictCrypto) {
	d := *dcc
	if len(d) > 1 {
		for _, v := range d {
			DCout = append(DCout, v.(DictCrypto))
		}
		return DCout
	}
	return []DictCrypto{}
}
func (dcc *DictCryptoCache) GetTop10Cache() (DCout []DictCrypto, err error) {
	d := *dcc
	if len(d) > 1 {
		cnt := 0
		for _, v := range d {
			if cnt == 10 {
				break
			}
			if v.(DictCrypto).CryptoCounter >= 1 {
				DCout = append(DCout, v.(DictCrypto))
				cnt++
			}
		}
		if cnt < 10 {
			for i := cnt; i < 10; i++ {
				for _, v := range d {
					isFound := false
					for _, v1 := range DCout {
						if v.(DictCrypto).CryptoId == v1.CryptoId {
							isFound = true
						}
					}
					if !isFound {
						DCout = append(DCout, v.(DictCrypto))
						break
					}
				}
			}
		}
		return DCout, err
	} else {
		return []DictCrypto{}, fmt.Errorf("GetCache:Crypto not initialised")
	}
}

// Поиск ИД криптовалюты по Имени в словаре
func (dcck *DictCryptoCacheKeys) GetCacheIdByName(name string) int {
	dK := *dcck
	v, ok := dK[strings.ToUpper(name)]
	if ok {
		return v
	}
	return 0
}

// Кеширование лимитов
type LimitsCache map[int]Limits

var LmtCache = make(LimitsCache)

func (lmt *LimitsCache) InitCache() error {
	l := *lmt
	if _, ok := l[1]; !ok {
		// Заполняем информацию в кеш из БД
		expLst := []Expressions{
			{Key: "IdLmt", Operator: NotEQ, Value: "0"},
		}
		rs, find, _, err := ReadDataRow(&Limits{}, expLst, 0)
		if err != nil {
			return fmt.Errorf("InitCache:" + err.Error())
		}
		if !find {
			return nil
		}
		subFields := Limits{}
		for _, subRs := range rs {
			mapstructure.Decode(subRs, &subFields)
			l[subFields.IdLmt] = subFields
		}
	}
	return nil
}

func (lmt *LimitsCache) GetCacheById(id int) (Limits, error) {
	l := *lmt
	if v, ok := l[id]; !ok {
		return Limits{}, fmt.Errorf("GetCacheById:Limit not initialized")
	} else {
		return v, nil
	}
}

func (lmt *LimitsCache) GetCacheLastId() int {
	l := *lmt
	var maxId int
	if _, ok := l[1]; ok {
		for maxId = range l {
			break
		}
		for n := range l {
			if n > maxId {
				maxId = n
			}
		}
	}
	return maxId + 1
}

// Кеш словаря лимитов
type LimitsCacheKeys map[string]LimitsDict // Словарь symbolName - Id
var LmtCacheKeys = make(LimitsCacheKeys)

func (lmtK *LimitsCacheKeys) InitCache() error {
	l := *lmtK
	if _, ok := l["LMT001"]; !ok {
		// Заполняем информацию в кеш из БД
		expLst := []Expressions{
			{Key: "IdLmtDct", Operator: NotEQ, Value: "0"},
		}
		rs, find, _, err := ReadDataRow(&LimitsDict{}, expLst, 0)
		if err != nil {
			return fmt.Errorf("InitCache:" + err.Error())
		}
		if !find {
			return nil
		}
		subFields := LimitsDict{}
		for _, subRs := range rs {
			mapstructure.Decode(subRs, &subFields)
			l[subFields.NameLmtDct] = subFields
		}
	}
	return nil
}

func (lmtK *LimitsCacheKeys) GetCacheById(symbol string) (LimitsDict, error) {
	lK := *lmtK
	if v, ok := lK[symbol]; !ok {
		return LimitsDict{}, fmt.Errorf("GetCacheById:Limit not initialized")
	} else {
		return v, nil
	}
}
