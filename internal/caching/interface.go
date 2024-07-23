package caching

import "github.com/mbydanov/tg_golang_bot/internal/database"

type iCacheble interface {
	database.Users | database.TrackingCrypto | database.Limits
}
type iCache[T any] interface {
	Get(int) (T, bool)
	Set(int, T)
	Delete(int)
	URLockU() bool
	URUnlock() bool
}
