package tgbot

import (
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/mbydanov/tg_golang_bot/internal/caching"
	"github.com/mbydanov/tg_golang_bot/internal/database"
)

type UserInfo = database.Users

var MenuCache = caching.Init[MenuInfo](time.Minute*5, time.Hour*12)

type MenuInfo struct {
	Crypto       string
	Criterion    string
	Price        float32
	CurrentPrice float32
	IdTracking   int
	IdCrypto     int
	IdCriterion  int
	OffsetNavi   int
	CurrentMenu  string
}
type UpdateBot struct {
	User UserInfo
	Data []string
	Menu MenuInfo
}

func (ub *UpdateBot) FillInfo(update *tgbotapi.Update) (err error) {
	ub.User = FindUserIdFromUpdate(update)
	ub.Menu, _ = caching.GetCacheByIdxInMap(MenuCache, ub.User.IdUsr)

	// Определение события
	if update.Message != nil {
		command := update.Message.Command()
		if command != "" {
			ub.Data = append(ub.Data, command)
		} else if update.Message.Text != "" {
			if ub.Menu.CurrentMenu == GetCrypto {
				ub.Data = []string{GetCryptoCurr, update.Message.Text}
			}
		}
	} else if update.CallbackQuery != nil {
		ub.Data = strings.Split(update.CallbackQuery.Data, "_")
	}

	return nil
}
