package database

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// Кеширование пользователей
type UsersCacheType map[int]Users

var UsersCache = make(UsersCacheType)

func (uc *UsersCacheType) CheckCache(idUsr int) error {
	if _, ok := UsersCache[idUsr]; !ok {
		// Заполняем информацию в кеш из БД
		user := Users{IdUsr: idUsr}
		if err := user.CheckUser(); err != nil {
			return fmt.Errorf("CheckCache:" + err.Error())
		}
		UsersCache[idUsr] = user
	}
	return nil
}

func (uc *UsersCacheType) GetCache(idUsr int) (Users, error) {
	if v, ok := UsersCache[idUsr]; !ok {
		return Users{}, fmt.Errorf("GetCache:User not initialised")
	} else {
		return v, nil
	}
}

// Кеш типов отслеживаний
type TypeTrackingCryptoCache map[int]TypeTrackingCrypto

var TypeTCCache = make(TypeTrackingCryptoCache)

func (ttc *TypeTrackingCryptoCache) CheckAllCache() error {
	if _, ok := TypeTCCache[1]; !ok {
		// Заполняем информацию в кеш из БД
		typeTC := TypeTrackingCrypto{}
		rs, err := typeTC.GetAllTypeInfo()
		if err != nil {
			return fmt.Errorf("CheckCache:" + err.Error())
		}
		for _, v := range rs {
			TypeTCCache[v.(TypeTrackingCrypto).IdTypTrkCrp] = v.(TypeTrackingCrypto)
		}
	}
	return nil
}
func (ttc *TypeTrackingCryptoCache) CheckCache(id int) error {
	if _, ok := TypeTCCache[id]; !ok {
		// Заполняем информацию в кеш из БД
		typeTC := TypeTrackingCrypto{IdTypTrkCrp: id}
		v, err := typeTC.GetTypeInfo()
		if err != nil {
			return fmt.Errorf("CheckCache:" + err.Error())
		}
		TypeTCCache[v.(TypeTrackingCrypto).IdTypTrkCrp] = v.(TypeTrackingCrypto)
	}
	return nil
}

func (ttc *TypeTrackingCryptoCache) GetCache(idType int) (TypeTrackingCrypto, error) {
	if v, ok := TypeTCCache[idType]; !ok {
		return TypeTrackingCrypto{}, fmt.Errorf("GetCache:User not initialised")
	} else {
		return v, nil
	}
}

// Кеш активных отслеживаний
type TrackingCryptoCache map[int]TrackingCrypto

var TCCache = make(TrackingCryptoCache)

func (ttc *TrackingCryptoCache) CheckAllCache() error {
	if _, ok := TCCache[1]; !ok {
		// Заполняем информацию в кеш из БД
		expLst := []Expressions{
			{Key: "OnTrkCrp", Operator: EQ, Value: "true"},
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
			TCCache[subFields.IdTrkCrp] = subFields
		}
	}
	return nil
}
func (ttc *TrackingCryptoCache) GetCache(id int) (TrackingCrypto, error) {
	if v, ok := TCCache[id]; !ok {
		return TrackingCrypto{}, fmt.Errorf("GetCache:User not initialised")
	} else {
		return v, nil
	}
}