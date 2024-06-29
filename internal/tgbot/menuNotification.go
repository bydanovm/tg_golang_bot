package tgbot

import (
	"strconv"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/services"
	"github.com/sirupsen/logrus"
)

const (
	ChooseCrypto string = "Выберите или введите криптовалюту для отслеживания"
	ChooseSum    string = "Выберите или введите сумму для отслеживания"
	EnterCrypto  string = "Введите криптовалюту для отслеживания"
	EnterSum     string = "Введите сумму для отслеживания"
)

// Функция обработки меню по "Оповещениям"
func menuNotification(update *tgbotapi.Update, keyboardBot *tgBotMenu) (msg interface{}) {
	var ans string
	var keyboard tgbotapi.InlineKeyboardMarkup
	// Обработка callback
	if update.CallbackQuery != nil {
		// Проверка команд
		// Вынести в общее ответ по коллбеку
		// callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		// callback.ShowAlert = true
		// if _, err := bot.AnswerCallbackQuery(callback); err != nil {
		// 	services.Logging.WithFields(logrus.Fields{
		// 		"userId":   update.CallbackQuery.Message.Chat.ID,
		// 		"userName": update.CallbackQuery.Message.From.UserName,
		// 		"type":     "callback_answer",
		// 		"command":  update.CallbackQuery.Data,
		// 	}).Error()
		// }

		// Разберем data callback по структуре command_subcmd_addcmd
		callBackData := strings.Split(update.CallbackQuery.Data, "_")
		// Главное меню оповещений
		if len(callBackData) == 1 && callBackData[0] == SetNotif {
			keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotif), 2)
			// Текст взять из базы (нужен справочник)
			ans = "Здесь можно завести оповещения\n"

			msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID, ans)

			msg_t.ReplyMarkup = &keyboard
			msg = msg_t

		} else if len(callBackData) == 2 && callBackData[0] == SetNotif {
			// Оповещения 2 уровень
			switch update.CallbackQuery.Data {
			case SetNotifCrypto: // Выбрать крипту
				// Текст взять из базы (нужен справочник)
				ans = ChooseCrypto
				// Выбор ТОП-10 криптовалют
				top10cur, err := database.DCCache.GetTop10Cache()
				if err != nil {
					services.Logging.WithFields(logrus.Fields{
						"userId":   update.CallbackQuery.Message.From.ID,
						"userName": update.CallbackQuery.Message.From.UserName,
					}).Error(err)
				}
				var row []tgbotapi.InlineKeyboardButton
				for k, v := range top10cur {
					btn := tgbotapi.NewInlineKeyboardButtonData(v.CryptoName, SetNotifCrypto+"_"+v.CryptoName)
					row = append(row, btn)
					// Делим на N строк по 5 элементов
					if (k+1)%5 == 0 {
						keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
						row = nil
					}
				}
				row = append(row, tgbotapi.NewInlineKeyboardButtonData("Назад", Start))
				row = append(row, tgbotapi.NewInlineKeyboardButtonData("Ввод", SetNotifCryptoEnter))
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

			case SetNotifPrice: // Установить цену
				// Текст взять из базы (нужен справочник)
				ans = ChooseSum
			}

			msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID, ans)
			msg_t.ReplyMarkup = &keyboard

			msg = msg_t

		} else if len(callBackData) == 3 && callBackData[0] == SetNotif {
			// Оповещения 3 уровень
			switch update.CallbackQuery.Data {
			// Ввод своей КВ
			case SetNotifCryptoEnter:
				ans = EnterCrypto
				msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID, ans)
				msg_t.ReplyMarkup = &keyboard
				msg = msg_t

			case SetNotifCriterionMore:
				ans = "Выбран критерий \"Больше\"\n"

			case SetNotifCriterionLess:
				ans = "Выбран критерий \"Меньше\"\n"

			case SetNotifPriceYes:
				// Нажата ДА на последнем этапе, возврат в начало
				// Текст взять из базы (нужен справочник)
				ans = "Оповещение успешно создано\n"

			case SetNotifPriceNo:
				// Нажата НЕТ на последнем этапе, возврат в начало
				// Текст взять из базы (нужен справочник)
				ans = "Вы отменили создание оповещения\n"

			default:
				// Случай, когда пришло 3 аргумента и выбрана КВ через кнопки
				if callBackData[1] == Crypto {
					// Сохранить КВ в мапу
					// Переход к выбору критерия
					keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotifCriterion), 2)

					ans += "Выбрана криптовалюта: " + callBackData[2] + "\nВыберите критерий"
					msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
						update.CallbackQuery.Message.MessageID, ans)
					msg_t.ReplyMarkup = &keyboard
					msg = msg_t

				} else if callBackData[1] == Price {
					// Случай, когда пришло 3 аргумента и выбрана цена через кнопки
					// Проверить аргумент и сохранить цену в мапу
					// Переход к подтверждению
					keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotifPrice), 2)

					ans += "Введена цена: " + callBackData[2] + "\nПодтвердить?"
					msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
						update.CallbackQuery.Message.MessageID, ans)
					msg_t.ReplyMarkup = &keyboard
					msg = msg_t
				}
			}

			// Для выбора критериев переход к выбору или вводу суммы
			if update.CallbackQuery.Data == SetNotifCriterionMore || update.CallbackQuery.Data == SetNotifCriterionLess {
				ans += ChooseSum

				// Тут нужно получить данные о КВ и сформировать меню с предложениями цен
				prices := []int{1001, 1005, 1010, 999, 995, 990}
				var row []tgbotapi.InlineKeyboardButton

				for k, v := range prices {
					btn := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(v), SetNotifPrice+"_"+strconv.Itoa(v))
					row = append(row, btn)
					// Делим на N строк по 5 элементов
					if (k+1)%3 == 0 {
						keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
						row = nil
					}
				}
				row = append(row, tgbotapi.NewInlineKeyboardButtonData("Назад", Start))
				row = append(row, tgbotapi.NewInlineKeyboardButtonData("Ввод", SetNotifPriceEnter))
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

				msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID, ans)
				msg_t.ReplyMarkup = &keyboard

				msg = msg_t

			} else if update.CallbackQuery.Data == SetNotifPriceYes || update.CallbackQuery.Data == SetNotifPriceNo {
				ans += "Здесь можно завести оповещения\n"

				keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotif), 2)
				msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID, ans)
				msg_t.ReplyMarkup = &keyboard
				msg = msg_t
			}
		}
	} else if update.Message.ReplyToMessage != nil {

		switch update.Message.ReplyToMessage.Text {
		case EnterCrypto:
			keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotifCriterion), 2)

			ans = "Выбрана криптовалюта: " + update.Message.Text + "\nВыберите критерий"
			msg_t := tgbotapi.NewMessage(update.Message.Chat.ID,
				ans)
			msg_t.ReplyMarkup = keyboard
			msg = msg_t

		case EnterSum:
			keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotifPrice), 2)

			ans = "Введена сумма: " + update.Message.Text + "\nОтслеживать: {тут то что получилось}\nПодтвердить?"
			msg_t := tgbotapi.NewMessage(update.Message.Chat.ID,
				ans)
			msg_t.ReplyMarkup = keyboard
			msg = msg_t
		}
	}
	return msg
}

// msg_t.ReplyMarkup = tgbotapi.ForceReply{
// 	ForceReply: true,
// }
