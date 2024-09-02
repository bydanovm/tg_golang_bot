package caching

import (
	"time"
)

type iCacheble interface {
	// database.Users | database.TrackingCrypto | database.Limits
	any
}

type iCacher[T iCacheble] interface {
	Get(interface{}) (T, bool)
	GetKeyByIdx(string, int) int
	GetByIdxInMap(interface{}, int) (T, bool)
	GetKeyChain(in interface{}) []interface{}
	GetKeyChainSort(in interface{}) []interface{}
	SetLite(interface{}, T, time.Duration)
	Set(interface{}, T, time.Duration)
	Add(interface{}, T)
	Update(k interface{}, val T)
	Delete(interface{})
	// Pop(int)
	// DropByIdx(int, int)
	DropAll()
	GetCacheCountRecord() int
	GetCacheSortCountRecord() int
	URLockU() bool
	URUnlock() bool
}
