package caching

import "time"

type iCacheble interface {
	// database.Users | database.TrackingCrypto | database.Limits
	any
}
type iCache[T any] interface {
	Get(int) (T, bool)
	Set(int, T, time.Duration)
	Delete(int)
	URLockU() bool
	URUnlock() bool
}
