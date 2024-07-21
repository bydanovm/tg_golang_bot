package cacheghost

import (
	"sync"

	"github.com/mbydanov/tg_golang_bot/internal/database"
)

type UserFullCache struct {
	database.Users
	database.TrackingCrypto
	database.Limits
}

type UserCacheMap map[int]UserFullCache

type UsersCache struct {
	mu   sync.RWMutex
	Item UserCacheMap
}
