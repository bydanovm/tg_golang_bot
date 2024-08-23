package tgbot

import (
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/mbydanov/tg_golang_bot/internal/caching"
	"github.com/mbydanov/tg_golang_bot/internal/database"
)

type UserInfo = database.Users

var MenuCache = caching.Init[MenuInfo](time.Minute*5, time.Second*150)

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
			if ub.Menu.CurrentMenu == GetCrypto ||
				ub.Menu.CurrentMenu == GetCryptoBack ||
				ub.Menu.CurrentMenu == GetCryptoNext {
				ub.Data = []string{GetCryptoCurr, strings.ToUpper(update.Message.Text)}
			} else if ub.Menu.CurrentMenu == SetNotif ||
				ub.Menu.CurrentMenu == SetNotifBack ||
				ub.Menu.CurrentMenu == SetNotifNext {
				ub.Data = []string{SetNotifPrice, strings.ToUpper(update.Message.Text)}
			}
		}
	} else if update.CallbackQuery != nil {
		ub.Data = strings.Split(update.CallbackQuery.Data, "_")
	}

	// Пишем текущее нахождение пользователя в меню
	if len(ub.Data) >= 1 && ub.Data != nil {
		ub.Menu.CurrentMenu = ub.Data[0]
		caching.SetCache(MenuCache, ub.User.IdUsr, ub.Menu, 0)
	}

	return nil
}
