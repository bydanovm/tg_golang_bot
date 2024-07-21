package cacheghost

import (
	"github.com/mbydanov/tg_golang_bot/internal/database"
)

var UsersInfoCache = Init()

// Инициализация
func Init() *UsersCache {
	items := make(UserCacheMap)

	cache := UsersCache{
		Item: items,
	}
	return &cache

}

// *UsersCache - пользовательский кеш

func (uc *UsersCache) URLockU() (isRLock bool) {
	uc.mu.RLock()
	isRLock = true
	return isRLock
}
func (uc *UsersCache) URUnlock() (isRLock bool) {
	uc.mu.RUnlock()
	isRLock = false
	return isRLock
}

func (uc *UsersCache) GetObject(usrId int) (ufc UserFullCache, find bool) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if v, ok := uc.Item[usrId]; ok {
		ufc = v
		find = true
	}
	return ufc, find

}

func (uc *UsersCache) GetUserInfo(usrId int) (usr database.Users) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if v, ok := uc.Item[usrId]; ok {
		usr = v.Users
	}
	return usr
}

func (uc *UsersCache) GetTracking(usrId int) (trk database.TrackingCrypto) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if v, ok := uc.Item[usrId]; ok {
		trk = v.TrackingCrypto
	}
	return trk
}

func (uc *UsersCache) GetLimits(usrId int) (lmt database.Limits) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if v, ok := uc.Item[usrId]; ok {
		lmt = v.Limits
	}
	return lmt
}

// Записываем информацию о пользователе
func (uc *UsersCache) SetUserInfo(usr database.Users) int {
	// Попытка получения всего объекта из кеша
	// Чтобы не перезатереть имеющиеся данные
	if object, find := uc.GetObject(usr.IdUsr); find {
		object.Users = usr
		uc.mu.Lock()
		uc.Item[usr.IdUsr] = object
	} else {
		uc.mu.Lock()
		uc.Item[usr.IdUsr] = UserFullCache{Users: usr}
	}
	defer uc.mu.Unlock()
	return usr.IdUsr
}

// Записываем информацию об отслеживаниях пользователя
func (uc *UsersCache) SetUserTracking(trk database.TrackingCrypto) int {
	// Попытка получения всего объекта из кеша
	// Чтобы не перезатереть имеющиеся данные
	if object, find := uc.GetObject(trk.UserId); find {
		object.TrackingCrypto = trk
		uc.mu.Lock()
		uc.Item[trk.UserId] = object
	} else {
		uc.mu.Lock()
		uc.Item[trk.UserId] = UserFullCache{TrackingCrypto: trk}
	}
	defer uc.mu.Unlock()
	return trk.UserId
}

func (uc *UsersCache) SetUser() {

}
