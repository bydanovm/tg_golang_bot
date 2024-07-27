package tgbot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/mbydanov/tg_golang_bot/internal/caching"
)

const (
	ChooseGetCrypto string = "Выберите или введите криптовалюту"
	EnterGetCrypto  string = "Введите криптовалюту"
)

func menuGetCrypto(update *tgbotapi.Update, keyboardBot *tgBotMenu) (msg interface{}) {
	var ans string
	var keyboard tgbotapi.InlineKeyboardMarkup
	// Обработка callback
	if update.CallbackQuery != nil {
		// Разберем data callback по структуре command_subcmd_addcmd
		callBackData := strings.Split(update.CallbackQuery.Data, "_")
		// Главное меню оповещений
		if len(callBackData) == 1 && callBackData[0] == GetCrypto {
			// Текст взять из базы (нужен справочник)
			ans = ChooseGetCrypto

			// Выбор ТОП-10 криптовалют
			// top10cur, err := database.DCCache.GetTop10Cache()
			offset := 10
			keyboard = GetCryptoListOffset(offset)

			msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID, ans)
			msg_t.ReplyMarkup = &keyboard

			msg = msg_t

		} else if len(callBackData) == 2 && callBackData[0] == GetCrypto {
			switch update.CallbackQuery.Data {
			// Ввод КВ вручную
			case GetCryptoEnter:
				ans = EnterGetCrypto
				msg_t := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					ans)
				msg_t.ReplyMarkup = tgbotapi.ForceReply{
					ForceReply: true,
				}
				msg = msg_t

			// case GetCryptoYet:
			// Когда пришла в ответе КВ из кнопки

			default:
				ans, keyboard = GetCryptoFunc(callBackData[1])
				msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID, ans)
				msg_t.ReplyMarkup = &keyboard
				msg = msg_t
			}

		} else if len(callBackData) == 3 && callBackData[0] == GetCrypto {
			if callBackData[1] == Yet {
				ans = ChooseGetCrypto

				offset, _ := strconv.Atoi(callBackData[2])
				keyboard = GetCryptoListOffset(offset)

				msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID, ans)
				msg_t.ReplyMarkup = &keyboard

				msg = msg_t

			}
		}
	} else if update.Message.ReplyToMessage != nil {

		switch update.Message.ReplyToMessage.Text {
		case EnterGetCrypto:
			// Ручной ввод отключен, так как нет словаря валюта <-> ид
			// ans, keyboard = GetCryptoFunc(update.Message.Text)

			// msg_t := tgbotapi.NewMessage(update.Message.Chat.ID,
			// 	ans)
			// msg_t.ReplyMarkup = &keyboard
			// msg = msg_t
		}
	}

	return msg

}

func GetCryptoFunc(idStr string) (ans string, keyboard tgbotapi.InlineKeyboardMarkup) {
	// Получение данных из кеша
	id, err := strconv.Atoi(idStr)
	if err != nil {

	}
	crypto, err := caching.GetCacheByIdxInMap(caching.CryptoCache, id, 0)
	if err == nil {
		ans += fmt.Sprintf("1 %s = "+FormatFloatToString(crypto.CryptoLastPrice)+" %s",
			crypto.CryptoName,
			crypto.CryptoLastPrice,
			"$")
	} else {
		ans += fmt.Sprintf("%s %s", crypto.CryptoName, "не найдена в базе")
	}

	var row []tgbotapi.InlineKeyboardButton
	row = append(row, tgbotapi.NewInlineKeyboardButtonData("Назад", GetCrypto))
	// if len(crypto) == 1 {
	row = append(row, tgbotapi.NewInlineKeyboardButtonData("Установить отслеживание", SetNotifCrypto+"_"+crypto.CryptoName))
	// } else {
	// row = append(row, tgbotapi.NewInlineKeyboardButtonData("Ввод", GetCryptoEnter))
	// }
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	return ans, keyboard
}
func GetCryptoListOffset(offset int) (keyboard tgbotapi.InlineKeyboardMarkup) {
	listCryptoCur, lastList, _ := caching.GetCacheOffset(caching.CryptoCache, offset)

	var row []tgbotapi.InlineKeyboardButton
	for k, v := range listCryptoCur {
		btn := tgbotapi.NewInlineKeyboardButtonData(v.CryptoName, GetCrypto+"_"+strconv.Itoa(v.CryptoId))
		row = append(row, btn)
		// Делим на N строк по 5 элементов
		if (k+1)%5 == 0 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			row = nil
		}
	}

	if offset > 10 {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("Назад", GetCryptoYet+`_`+strconv.Itoa(offset-10)))
	} else {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("Назад", Start))
	}

	if !lastList {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("Ещё", GetCryptoYet+`_`+strconv.Itoa(offset+10)))
	}
	// row = append(row, tgbotapi.NewInlineKeyboardButtonData("Ввод", GetCryptoEnter))
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	return keyboard
}
