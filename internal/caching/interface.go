package caching

import "time"

type iCacheble interface {
	// database.Users | database.TrackingCrypto | database.Limits
	any
}
type iCacher[T any] interface {
	Get(int) ([]T, bool)
	GetKeyByIdx(string, int) int
	GetByIdxInMap(int, int) (T, bool)
	GetKeyChain(in interface{}) []interface{}
	GetKeyChainSort(in interface{}) []interface{}
	Set(int, T, time.Duration)
	Add(int, T)
	Update(k int, val T)
	Delete(int)
	Pop(int)
	DropByIdx(int, int)
	DropAll()
	GetCacheCountRecord() int
	GetCacheSortCountRecord() int
	URLockU() bool
	URUnlock() bool
}
