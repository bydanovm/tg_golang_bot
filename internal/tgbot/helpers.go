package tgbot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/mbydanov/tg_golang_bot/internal/caching"
	"github.com/mbydanov/tg_golang_bot/internal/database"
)

func FindUserIdFromUpdate(update *tgbotapi.Update) (userInfo UserInfo) {
	if update.Message != nil {
		userInfo.IdUsr = update.Message.From.ID
		userInfo.ChatIdUsr = update.Message.Chat.ID
		userInfo.NameUsr = update.Message.From.UserName
		userInfo.FirstName = update.Message.From.FirstName
		userInfo.LastName = update.Message.From.LastName
		userInfo.LangCode = update.Message.From.LanguageCode
		userInfo.IsBot = update.Message.From.IsBot
	} else if update.CallbackQuery != nil {
		userInfo.IdUsr = update.CallbackQuery.From.ID
		userInfo.ChatIdUsr = update.CallbackQuery.Message.Chat.ID
		userInfo.NameUsr = update.CallbackQuery.From.UserName
		userInfo.FirstName = update.CallbackQuery.From.FirstName
		userInfo.LastName = update.CallbackQuery.From.LastName
		userInfo.LangCode = update.CallbackQuery.From.LanguageCode
		userInfo.IsBot = update.CallbackQuery.From.IsBot
	} else if update.InlineQuery != nil {
		userInfo.IdUsr = update.InlineQuery.From.ID
		userInfo.ChatIdUsr = int64(update.InlineQuery.From.ID)
		userInfo.NameUsr = update.InlineQuery.From.UserName
		userInfo.FirstName = update.InlineQuery.From.FirstName
		userInfo.LastName = update.InlineQuery.From.LastName
		userInfo.LangCode = update.InlineQuery.From.LanguageCode
		userInfo.IsBot = update.InlineQuery.From.IsBot
	} else {
		userInfo.IdUsr = -1
		userInfo.ChatIdUsr = -1
		userInfo.NameUsr = ``
		userInfo.FirstName = ``
		userInfo.LastName = ``
		userInfo.LangCode = ``
		userInfo.IsBot = false
	}
	userInfo.IsBanned = false
	userInfo.IdLvlSec = 5
	return userInfo
}

// Единая точка проверки юзера на авторизацию
func checkAuthUser(bot *tgbotapi.BotAPI, update *tgbotapi.Update) (user database.Users, err error) {
	var msg tgbotapi.MessageConfig
	var ans string

	// Определение откуда пришел запрос
	userInfo := FindUserIdFromUpdate(update)
	// Проверка нахождения пользователя в кеше (БД)
	// с возможностью записи в базу нового пользователя
	user, err = caching.CheckCacheAndWrite(caching.UsersCache, userInfo.IdUsr, userInfo)
	if err != nil {
		ans = fmt.Sprintf("Извините. При авторизации возникла какая-то ошибка.\nПопробуйте позже /%s", Start)
		msg = tgbotapi.NewMessage(user.ChatIdUsr,
			ans)
		bot.Send(msg)
	}

	return user, err
}

// Функция проверки передачи информации между функциями, при необходимости достает информацию из менюкеша
func checkCallbackData(update *tgbotapi.Update, size int) ([]string, error) {
	callBackData := strings.Split(update.CallbackQuery.Data, "_")
	if len(callBackData) < size {
		userInfo := FindUserIdFromUpdate(update)
		// Вызов может поступить из другого пункта меню, проверяем наличие в кеше
		menuCache, err := caching.GetCacheByIdxInMap(MenuCache, userInfo.IdUsr, 0)
		if err != nil {
			return nil, err
		}
		if menuCache.IdCrypto != 0 {
			callBackData = append(callBackData, strconv.Itoa(menuCache.IdCrypto))
		} else {
			return nil, err
		}
	}
	return callBackData, nil
}

func clearSetNotifMenuCache(updateBot *UpdateBot) {
	// Очистка данных о КВ
	updateBot.Menu = MenuInfo{}
	caching.SetCache(MenuCache, updateBot.User.IdUsr, updateBot.Menu, 0)
}
