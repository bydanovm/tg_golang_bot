package caching

import "time"

type iCacheble interface {
	// database.Users | database.TrackingCrypto | database.Limits
	any
}
type iCacher[T any] interface {
	Get(int) ([]T, bool)
	GetByIdx(int, int) (T, bool)
	Set(int, T, time.Duration)
	Add(int, T)
	Delete(int)
	Pop(int)
	DropByIdx(int, int)
	URLockU() bool
	URUnlock() bool
}
