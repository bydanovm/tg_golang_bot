package database

import (
	"fmt"
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
