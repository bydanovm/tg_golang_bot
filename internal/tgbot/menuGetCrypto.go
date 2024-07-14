package tgbot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/mbydanov/tg_golang_bot/internal/coinmarketcup"
	"github.com/mbydanov/tg_golang_bot/internal/database"
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
			ans, keyboard = GetCryptoFunc(update.Message.Text)

			msg_t := tgbotapi.NewMessage(update.Message.Chat.ID,
				ans)
			msg_t.ReplyMarkup = &keyboard
			msg = msg_t
		}
	}

	return msg

}

func GetCryptoFunc(crypto string) (ans string, keyboard tgbotapi.InlineKeyboardMarkup) {
	// Получение данных из базы
	cryptos, err := coinmarketcup.GetLatestStruct(crypto)
	isFind := false
	if err == nil {
		for _, v := range cryptos {
			if v.Find {
				isFind = true
				ans += fmt.Sprintf("1 %s = "+FormatFloatToString(v.Crypto.CryptoLastPrice)+" %s",
					v.Crypto.CryptoName,
					v.Crypto.CryptoLastPrice,
					"$")
			} else {
				ans += fmt.Sprintf("%s %s", v.Crypto.CryptoName, "не найдена в базе")
			}
		}
	}

	var row []tgbotapi.InlineKeyboardButton
	row = append(row, tgbotapi.NewInlineKeyboardButtonData("Назад", GetCrypto))
	if isFind && len(cryptos) == 1 {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("Установить отслеживание", SetNotifCrypto+"_"+crypto))
	} else {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("Ввод", GetCryptoEnter))
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	return ans, keyboard
}
func GetCryptoListOffset(offset int) (keyboard tgbotapi.InlineKeyboardMarkup) {
	listCryptoCur, lastList, _ := database.DCCache.GetCryptoOffset(offset)
	// if err != nil {
	// 	services.Logging.WithFields(logrus.Fields{
	// 		"userId":   update.CallbackQuery.Message.From.ID,
	// 		"userName": update.CallbackQuery.Message.From.UserName,
	// 	}).Error(err)
	// }
	var row []tgbotapi.InlineKeyboardButton
	for k, v := range listCryptoCur {
		btn := tgbotapi.NewInlineKeyboardButtonData(v.CryptoName, GetCrypto+"_"+v.CryptoName)
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
	row = append(row, tgbotapi.NewInlineKeyboardButtonData("Ввод", GetCryptoEnter))
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	return keyboard
}
