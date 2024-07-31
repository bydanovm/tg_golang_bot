package tgbot

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/mbydanov/tg_golang_bot/internal/caching"
	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/models"
	"github.com/mbydanov/tg_golang_bot/internal/notifications"
	"github.com/mbydanov/tg_golang_bot/internal/services"
	"github.com/sirupsen/logrus"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// Создаем бота
func TelegramBot(chanModules chan models.StatusChannel) {
	// var msg interface{}
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
		user, err := checkAuthUser(bot, &update)
		if err != nil {
			services.Logging.WithFields(logrus.Fields{
				"module": "tgbot",
				"user":   user.IdUsr,
			}).Error(err.Error())
		}
		menuHandler(&update, *bot)
		// msg = menuGetCrypto(&update, keyboardBot)
		// if msg == nil {
		// 	msg = menuNotification(&update, keyboardBot)
		// }

		// var command = ""

		// if update.Message != nil {
		// 	command = update.Message.Command()
		// 	if command != "" {
		// 		// Обработка команд
		// 		// start - Начало
		// 		// number_of_users - Получить количество активных пользователей
		// 		// getcrypto - Получить актуальную информацию по криптовалюте
		// 		// setnotif - Установить уведомления по изменению цены криптовалюты
		// 		switch command {
		// 		case Start:
		// 			// Проверяем есть ли пользователь в кеше или базе
		// 			// msg = menuStart(&update, keyboardBot)
		// 		default:
		// 			// ans := "Команда /" + command + " не найдена.\n" +
		// 			// 	"Воспользуйся командой /" + Start + " для знакомства со мной."

		// 			// Отправлем сообщение
		// 			// msg = tgbotapi.NewMessage(database.UsersCache.GetChatId(update.Message.From.ID),
		// 			// 	ans)
		// 		}
		// 	}
		// }
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
			// callBackData := strings.Split(update.CallbackQuery.Data, "_")
			// Получение инфо о крипте
			// if len(callBackData) == 1 {
			// 	if callBackData[0] == Start {
			// 		msg = menuStart(&update, keyboardBot)
			// 	}
			// }
		}

		// switch msgConv := msg.(type) {
		// case tgbotapi.EditMessageTextConfig:
		// 	bot.Send(msgConv)
		// case tgbotapi.MessageConfig:
		// 	bot.Send(msgConv)
		// }
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
