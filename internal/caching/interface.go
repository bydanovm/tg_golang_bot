package caching

import "time"

type iCacheble interface {
	// database.Users | database.TrackingCrypto | database.Limits
	any
}
type iCacher[T any] interface {
	Get(int) ([]T, bool)
	GetKeyByIdx(int) int
	GetByIdxInMap(int, int) (T, bool)
	GetKeyChain(in interface{}) []interface{}
	Set(int, T, time.Duration)
	Add(int, T)
	Delete(int)
	Pop(int)
	DropByIdx(int, int)
	GetCacheCountRecord() int
	URLockU() bool
	URUnlock() bool
}
