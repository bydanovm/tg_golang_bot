package tgbot

import (
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/caching"
	"github.com/mbydanov/tg_golang_bot/internal/database"
)

type UserInfo = database.Users

var MenuCache = caching.Init[SetNotifStruct](time.Minute*5, time.Hour*12)
