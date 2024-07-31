package tgbot

import "sync"

type SetNotifCacheStruct struct {
	mu   sync.RWMutex
	Item SetNotifMap
}
type SetNotifStruct struct {
	Crypto      string
	Criterion   string
	Price       float32
	IdTracking  int
	IdCrypto    int
	OffsetNavi  int
	CurrentMenu string
}

// Номер пользователя (чата) - структура
type SetNotifMap map[int]SetNotifStruct

var SetNotifCh = Init()

func Init() *SetNotifCacheStruct {
	items := make(SetNotifMap)

	cache := SetNotifCacheStruct{
		Item: items,
	}

	return &cache
}

func (uc *SetNotifCacheStruct) URLockU() (isRLock bool) {
	uc.mu.RLock()
	isRLock = true
	return isRLock
}
func (uc *SetNotifCacheStruct) URUnlock() (isRLock bool) {
	uc.mu.RUnlock()
	isRLock = false
	return isRLock
}

func (uc *SetNotifCacheStruct) SetCrypto(idUsr int, crypto string) {
	isRLock := uc.URLockU()
	if _, ok := uc.Item[idUsr]; !ok {
		isRLock = uc.URUnlock()
		uc.mu.Lock()
		uc.Item[idUsr] = SetNotifStruct{Crypto: crypto}
		uc.mu.Unlock()
	} else {
		item := uc.Item[idUsr]
		item.Crypto = crypto
		isRLock = uc.URUnlock()
		uc.mu.Lock()
		uc.Item[idUsr] = item
		uc.mu.Unlock()
	}
	if isRLock {
		uc.mu.RUnlock()
	}
}

func (uc *SetNotifCacheStruct) SetCriterion(idUsr int, criterion string) {
	isRLock := uc.URLockU()
	if _, ok := uc.Item[idUsr]; !ok {
		isRLock = uc.URUnlock()
		uc.mu.Lock()
		uc.Item[idUsr] = SetNotifStruct{Criterion: criterion}
		uc.mu.Unlock()
	} else {
		item := uc.Item[idUsr]
		item.Criterion = criterion
		isRLock = uc.URUnlock()
		uc.mu.Lock()
		uc.Item[idUsr] = item
		uc.mu.Unlock()
	}
	if isRLock {
		uc.mu.RUnlock()
	}
}

func (uc *SetNotifCacheStruct) SetPrice(idUsr int, price float32) {
	isRLock := uc.URLockU()
	if _, ok := uc.Item[idUsr]; !ok {
		isRLock = uc.URUnlock()
		uc.mu.Lock()
		uc.Item[idUsr] = SetNotifStruct{Price: price}
		uc.mu.Unlock()
	} else {
		item := uc.Item[idUsr]
		item.Price = price
		isRLock = uc.URUnlock()
		uc.mu.Lock()
		uc.Item[idUsr] = item
		uc.mu.Unlock()
	}
	if isRLock {
		uc.mu.RUnlock()
	}
}

func (uc *SetNotifCacheStruct) GetCrypto(idUsr int) (crypto string) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if v, ok := uc.Item[idUsr]; ok {
		crypto = v.Crypto
	}
	return crypto
}

func (uc *SetNotifCacheStruct) GetCriterion(idUsr int) (criterion string) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if v, ok := uc.Item[idUsr]; ok {
		criterion = v.Criterion
	}
	return criterion
}

func (uc *SetNotifCacheStruct) GetPrice(idUsr int) (price float32) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if v, ok := uc.Item[idUsr]; ok {
		price = v.Price
	}
	return price
}

func (uc *SetNotifCacheStruct) GetObject(idUsr int) (object SetNotifStruct) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if v, ok := uc.Item[idUsr]; ok {
		object = v
	}
	return object
}
