package tgbot

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/mbydanov/tg_golang_bot/internal/caching"
)

// Если в хендлер ничего не пришло, то отображаем начало
// Если в хандлер что-то пришло, то вызываем соответствующую функцию из связанного списка
func menuHandler(update *tgbotapi.Update, bot tgbotapi.BotAPI) {
	var msg interface{}

	// Первичная инициализация меню хендлера
	if !keyboardBot.Init {
		keyboardBot.Add(Start, "Главная", "", true, funcMenuStart)
		keyboardBot.Add(GetCrypto, "Узнать курс", Start, true, funcGetCrypto)
		keyboardBot.Add(GetCryptoCurr, "Узнать курс валюты", GetCrypto, false, funcGetCryptoCurr)
		keyboardBot.Add(GetCryptoCurrSetNot, "Установить отслеживание", GetCryptoCurr, true)
		keyboardBot.Add(GetCryptoNext, "Дальше", GetCrypto, true)
		keyboardBot.Add(GetNotif, "Оповещения", Start, true, funcGetNotif)
		keyboardBot.Add(GetNotifId, "Получить отслеживание по ID", GetNotif, false)
		keyboardBot.Add(SetNotif, "Новое отслеживание", GetNotif, true, funcSetNotif)
		keyboardBot.Add(SetNotifPrice, "Установить цену", SetNotif, true, funcSetNotifPrice)
		keyboardBot.Add(GetNotifIdOn, "Отключить", GetNotifId, true)
		keyboardBot.Add(GetNotifIdOff, "Включить", GetNotifId, true)

		keyboardBot.Add(Help, "Справка", Start, true)

		keyboardBot.Init = true
	}

	if update.Message != nil {
		command := update.Message.Command()
		if command != "" {
			pfunc := keyboardBot.GetFunc(command)
			if pfunc != nil {
				ans, keyboard := pfunc(update)
				msg_t := tgbotapi.NewMessage(update.Message.Chat.ID,
					ans)
				msg_t.ReplyMarkup = &keyboard
				msg = msg_t
			}
		}
	} else if update.CallbackQuery != nil {
		// Разберем data callback по структуре command_param
		callBackData := strings.Split(update.CallbackQuery.Data, "_")
		pfunc := keyboardBot.GetFunc(callBackData[0])
		if pfunc != nil {
			ans, keyboard := pfunc(update)
			msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID, ans)
			msg_t.ReplyMarkup = &keyboard
			msg = msg_t
		}
	}

	switch msgConv := msg.(type) {
	case tgbotapi.EditMessageTextConfig:
		bot.Send(msgConv)
	case tgbotapi.MessageConfig:
		bot.Send(msgConv)
	}
}

func funcGetCrypto(update *tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	ans := ChooseGetCrypto
	offset := 10
	callBackData := strings.Split(update.CallbackQuery.Data, "_")
	listCryptoCur, _, _ := caching.GetCacheOffset(caching.CryptoCache, offset)
	listButtons := make([]buttonInfo, 0, 10)
	for _, v := range listCryptoCur {
		listButtons = append(listButtons, buttonInfo{v.CryptoName, GetCryptoCurr + `_` + strconv.Itoa(v.CryptoId)})
	}
	keyboard := ConvertToButtonInlineKeyboard(listButtons, callBackData[0], 3)
	return ans, keyboard
}

func funcGetCryptoCurr(update *tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	var ans string
	callBackData := strings.Split(update.CallbackQuery.Data, "_")
	id, err := strconv.Atoi(callBackData[1])
	if err != nil {

	}
	crypto, err := caching.GetCacheByIdxInMap(caching.CryptoCache, id, 0)
	if err == nil {
		ans = fmt.Sprintf("1 %s = "+FormatFloatToString(crypto.CryptoLastPrice)+" %s",
			crypto.CryptoName,
			crypto.CryptoLastPrice,
			"$")
	} else {
		ans = fmt.Sprintf("%s %s", crypto.CryptoName, "не найдена в базе")
	}
	keyboard := MenuToInlineFromNode(callBackData[0], 2)
	return ans, keyboard
}

func funcGetNotif(update *tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	ans := "Текущие оповещения"
	// Вывести от новых к старым в формате Валюта - Значение
	trackings, _ := caching.GetCacheRecordsKeyChain(caching.TrackingCache, update.CallbackQuery.From.ID, true)
	// Создание списка кнопок
	listButtons := make([]buttonInfo, 0, 10)
	for _, v := range trackings {
		infoCurrency, _ := caching.GetCacheByIdxInMap(caching.CryptoCache, v.DctCrpId, 0)
		listButtons = append(listButtons, buttonInfo{infoCurrency.CryptoName + " - " + fmt.Sprintf(FormatFloatToString(v.ValTrkCrp), v.ValTrkCrp) + " $", GetNotifId + "_" + fmt.Sprintf("%d", v.IdTrkCrp)})
	}
	keyboard := ConvertToButtonInlineKeyboard(listButtons, GetNotif, 3)
	return ans, keyboard
}

func funcSetNotif(update *tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	callBackData := strings.Split(update.CallbackQuery.Data, "_")
	offset := 10
	ans := ChooseGetCrypto

	listCryptoCur, _, _ := caching.GetCacheOffset(caching.CryptoCache, offset)
	listButtons := make([]buttonInfo, 0, 10)
	for _, v := range listCryptoCur {
		listButtons = append(listButtons, buttonInfo{v.CryptoName, SetNotifPrice + `_` + strconv.Itoa(v.CryptoId)})
	}
	keyboard := ConvertToButtonInlineKeyboard(listButtons, callBackData[0], 3)
	return ans, keyboard
}

func funcSetNotifPrice(update *tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	callBackData := strings.Split(update.CallbackQuery.Data, "_")

	// Возможно переключить на пакет caching
	// Пишем в кеш крипту
	SetNotifCh.SetCrypto(int(update.CallbackQuery.Message.Chat.ID), callBackData[1])

	ans := ChooseSum

	// Возможно переключить на пакет caching
	// Получаем из кеша нашу крипту
	crypto, err := strconv.Atoi(SetNotifCh.GetCrypto(int(update.CallbackQuery.Message.Chat.ID)))
	if err != nil {
		// Здесь должна быть обработка ошибки
	}

	infoCurrency, err := caching.GetCacheByIdxInMap(caching.CryptoCache, crypto, 0)
	if err != nil {
		// Здесь должна быть обработка ошибки
		ans += fmt.Sprintf("%s %s", "Выбранная криптовалюта", "не найдена в базе")
	}
	ans += "Выбрана криптовалюта: " + infoCurrency.CryptoName + "\n"

	prices := []PriceInfo{}

	prices = append(prices, PriceInfo{1, infoCurrency.CryptoLastPrice * 1.01})
	prices = append(prices, PriceInfo{5, infoCurrency.CryptoLastPrice * 1.05})
	prices = append(prices, PriceInfo{10, infoCurrency.CryptoLastPrice * 1.1})
	prices = append(prices, PriceInfo{-1, infoCurrency.CryptoLastPrice * 0.99})
	prices = append(prices, PriceInfo{-5, infoCurrency.CryptoLastPrice * 0.95})
	prices = append(prices, PriceInfo{-10, infoCurrency.CryptoLastPrice * 0.9})

	listButtons := make([]buttonInfo, 0, 10)
	for _, v := range prices {
		price := fmt.Sprintf("%f", v.Price)
		n := fmt.Sprintf(FormatFloatToString(v.Price)+" (%+d%%)", v.Price, v.Koeff)
		listButtons = append(listButtons, buttonInfo{n, SetNotifPrice + `_` + price})
	}
	keyboard := ConvertToButtonInlineKeyboard(listButtons, callBackData[0], 3)

	return ans, keyboard
}

func funcMenuStart(update *tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	ans := "Привет! Я - " + os.Getenv("BOT_NAME") + " помогу тебе знать актуальную информацию по криптовалюте\n" +
		"Используй кнопки ниже, чтобы узнать интересующую информацию.\n"
	keyboard := MenuToInlineFromNode(Start, 2)
	return ans, keyboard
}
