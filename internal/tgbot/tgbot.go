package tgbot

import (
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/models"
	"github.com/mbydanov/tg_golang_bot/internal/notifications"
	"github.com/mbydanov/tg_golang_bot/internal/services"
	"github.com/sirupsen/logrus"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// Создаем бота
func TelegramBot(chanModules chan models.StatusChannel) {
	var msg interface{}
	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		services.Logging.Panic(err.Error())
	}

	// keyboardBot := initMenu()
	// Инициализация меню

	// Устанавливаем время обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Функция получения сообщений от модулей
	go func() {
		for {
			v, ok := <-chanModules
			if ok {
				if v.Start {
					if v.Module == models.RetrieverCoins {
						// Отправка обратно в канал для нотификатора
						chanModules <- v
					} else if v.Module == models.Notificator {
						arr, ok := v.Data.([]notifications.NotificationsCCStruct)
						if ok {
							for _, v := range arr {
								msg := tgbotapi.NewMessage(int64(v.IdChat), v.Event)
								bot.Send(msg)
							}
						}
					}
				}
				if v.Error != nil {
					services.Logging.WithFields(logrus.Fields{
						"module": v.Module,
					}).Error(v.Error.Error())
				}
			}
		}
	}()

	// Получаем обновления от бота
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		services.Logging.Panic(err.Error())
	}

	for update := range updates {
		// Авторизация пользователя
		_, err := checkAuthUser(bot, &update)
		if err != nil {
			services.Logging.WithFields(logrus.Fields{
				"module": "tgbot",
			}).Error(err.Error())
		}
		msg = menuGetCrypto(&update, keyboardBot)
		if msg == nil {
			msg = menuNotification(&update, keyboardBot)
		}

		var command = ""

		if update.Message != nil {
			command = update.Message.Command()
			if command != "" {
				// Обработка команд
				// start - Начало
				// number_of_users - Получить количество активных пользователей
				// getcrypto - Получить актуальную информацию по криптовалюте
				// setnotif - Установить уведомления по изменению цены криптовалюты
				switch command {
				case Start:
					// Проверяем есть ли пользователь в кеше или базе
					msg = menuStart(&update, keyboardBot)
				default:
					ans := "Команда /" + command + " не найдена.\n" +
						"Воспользуйся командой /" + Start + " для знакомства со мной."

					// Отправлем сообщение
					msg = tgbotapi.NewMessage(database.UsersCache.GetChatId(update.Message.From.ID),
						ans)
				}
			}
		}
		if update.CallbackQuery != nil {
			// Проверка команд
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			callback.ShowAlert = true
			if _, err := bot.AnswerCallbackQuery(callback); err != nil {
				services.Logging.WithFields(logrus.Fields{
					"userId":   update.CallbackQuery.Message.Chat.ID,
					"userName": update.CallbackQuery.Message.From.UserName,
					"type":     "callback_answer",
					"command":  update.CallbackQuery.Data,
				}).Error()
			}
			// Разберем data callback по структуре command_cryptocur
			callBackData := strings.Split(update.CallbackQuery.Data, "_")
			// Получение инфо о крипте
			if len(callBackData) == 1 {
				if callBackData[0] == Start {
					msg = menuStart(&update, keyboardBot)
				}
			}
		}

		switch msgConv := msg.(type) {
		case tgbotapi.EditMessageTextConfig:
			bot.Send(msgConv)
		case tgbotapi.MessageConfig:
			bot.Send(msgConv)
		}
	}
}

func Filter(dcs []database.DictCrypto, fn func(dc database.DictCrypto) bool) []database.DictCrypto {
	var filtered []database.DictCrypto
	for _, v := range dcs {
		if fn(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// Единая точка проверки юзера на авторизацию
func checkAuthUser(bot *tgbotapi.BotAPI, update *tgbotapi.Update) (userId int, err error) {
	var msg tgbotapi.MessageConfig
	var chatId int64
	var userName string
	// Определение откуда пришел запрос
	if update.Message != nil {
		userId = update.Message.From.ID
		chatId = update.Message.Chat.ID
		userName = update.Message.From.UserName
	} else if update.CallbackQuery != nil {
		userId = update.CallbackQuery.From.ID
		chatId = update.CallbackQuery.Message.Chat.ID
		userName = update.CallbackQuery.Message.From.UserName
	} else if update.InlineQuery != nil {
		userId = update.InlineQuery.From.ID
		userName = update.InlineQuery.From.UserName
	} else {
		userId = 0
		chatId = 0
		err = fmt.Errorf("tgbot:checkAuthUser:Message is nil")
	}

	// Проверка нахождения пользователя в базе
	if err := database.UsersCache.CheckCache(userId); err != nil {
		services.Logging.WithFields(logrus.Fields{
			"module":   "tgbot",
			"userId":   userId,
			"userName": userName,
		}).Error(err.Error())
	}

	// Доп. проверка на получение идентификатора пользователя
	if ok := database.UsersCache.GetUserId(userId); ok == 0 {
		ans := fmt.Sprintf("Чтобы начать работу с ботом введите команду /%s", Start)
		// Отправлем сообщение
		msg = tgbotapi.NewMessage(chatId,
			ans)
		bot.Send(msg)
	}

	return userId, err
}
func FindUserIdFromUpdate(update *tgbotapi.Update) (userInfo UserInfo) {
	if update.Message != nil {
		userInfo.UserId = update.Message.From.ID
		userInfo.ChatId = update.Message.Chat.ID
		userInfo.UserName = update.Message.From.UserName
		userInfo.FirstName = update.Message.From.FirstName
		userInfo.LastName = update.Message.From.LastName
		userInfo.LanguageCode = update.Message.From.LanguageCode
		userInfo.IsBot = update.Message.From.IsBot
	} else if update.CallbackQuery != nil {
		userInfo.UserId = update.CallbackQuery.From.ID
		userInfo.ChatId = update.CallbackQuery.Message.Chat.ID
		userInfo.UserName = update.CallbackQuery.Message.From.UserName
		userInfo.FirstName = update.CallbackQuery.Message.From.FirstName
		userInfo.LastName = update.CallbackQuery.Message.From.LastName
		userInfo.LanguageCode = update.CallbackQuery.Message.From.LanguageCode
		userInfo.IsBot = update.CallbackQuery.Message.From.IsBot
	} else {
		userInfo.UserId = -1
		userInfo.ChatId = -1
		userInfo.UserName = ``
		userInfo.FirstName = ``
		userInfo.LastName = ``
		userInfo.LanguageCode = ``
		userInfo.IsBot = false
	}
	return userInfo
}
func menuStart(update *tgbotapi.Update, keyboardBot *tgBotMenu) (msg interface{}) {
	var ans string
	var keyboard tgbotapi.InlineKeyboardMarkup

	userInfo := FindUserIdFromUpdate(update)
	// Проверяем есть ли пользователь в кеше или базе
	if userId := database.UsersCache.GetUserId(userInfo.UserId); userId == 0 {
		// Пользователь не в кеше, ищем в БД, если не находим, то добавляем нового
		// Единственная точка входа, где пользователь может добавиться в БД
		user := database.Users{
			IdUsr:     userInfo.UserId,
			TsUsr:     time.Now(),
			NameUsr:   userInfo.UserName,
			FirstName: userInfo.FirstName,
			LastName:  userInfo.LastName,
			LangCode:  userInfo.LanguageCode,
			IsBot:     userInfo.IsBot,
			IsBanned:  false,
			ChatIdUsr: userInfo.ChatId,
			IdLvlSec:  5}
		// Поиск с последующим добавлением
		if err := user.CheckUser(); err != nil {
			// Отправляем сообщение в лог об ошибке
			services.Logging.Warn(err.Error())
		} else {
			// Кешируем добавленного пользователя
			if err := database.UsersCache.CheckCache(userInfo.UserId); err != nil {
				services.Logging.Warn(err.Error())
			}
		}
	}
	// Отправлем приветственное сообщение
	ans = "Привет! Я - " + os.Getenv("BOT_NAME") + " помогу тебе знать актуальную информацию по криптовалюте\n" +
		"Используй кнопки ниже, чтобы узнать интересующую информацию.\n"

	keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkup(), 2)

	if update.CallbackQuery != nil &&
		update.CallbackQuery.Message.From.UserName == os.Getenv("BOT_NAME") {
		msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID, ans)
		msg_t.ReplyMarkup = &keyboard
		msg = msg_t
	} else if update.Message != nil &&
		update.Message.From.UserName == os.Getenv("BOT_NAME") {
		msg_t := tgbotapi.NewEditMessageText(update.Message.Chat.ID,
			update.Message.MessageID, ans)
		msg_t.ReplyMarkup = &keyboard
		msg = msg_t
	} else {
		msg_t := tgbotapi.NewMessage(database.UsersCache.GetChatId(userInfo.UserId),
			ans)
		msg_t.ReplyMarkup = &keyboard
		msg = msg_t
	}

	return msg
}
