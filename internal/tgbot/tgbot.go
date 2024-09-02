package tgbot

import (
	"os"
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/exchange"
	"github.com/mbydanov/tg_golang_bot/internal/models"
	"github.com/mbydanov/tg_golang_bot/internal/notifications"
	"github.com/mbydanov/tg_golang_bot/internal/services"
	"github.com/sirupsen/logrus"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func newBot(token string) (*tgbotapi.BotAPI, error) {
	defer func() {
		if x := recover(); x != nil {
			services.Logging.Errorf("Panic:newBot:%v\n", x)
		}
	}()
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		services.Logging.WithFields(logrus.Fields{
			"module": "tgBot",
			"type":   "newBot",
			"status": "error",
		}).Errorf("%+v", err)
	} else {
		services.Logging.WithFields(logrus.Fields{
			"module": "tgBot",
			"type":   "newBot",
			"status": "ok",
		}).Info()
	}
	return bot, err
}

func getUpdates(bot *tgbotapi.BotAPI, u tgbotapi.UpdateConfig) (updates tgbotapi.UpdatesChannel, err error) {
	defer func() {
		if x := recover(); x != nil {
			services.Logging.Errorf("Panic:getUpdates:%v\n", x)
		}
	}()

	updates, err = bot.GetUpdatesChan(u)
	if err != nil {
		services.Logging.WithFields(logrus.Fields{
			"module": "tgBot",
			"type":   "getUpdates",
			"status": "error",
		}).Errorf("%+v", err)
	} else {
		services.Logging.WithFields(logrus.Fields{
			"module": "tgBot",
			"type":   "getUpdates",
			"status": "ok",
		}).Info()
	}
	return updates, err
}

func TelegramBot() {
	// Создаем бота
	var bot *tgbotapi.BotAPI
	var err error
	for {
		bot, err = newBot(os.Getenv("TOKEN"))
		if bot != nil && err == nil {
			break
		}
		<-time.After(time.Second * 3)
	}

	// Устанавливаем время обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Функция получения сообщений от модулей
	go func() {
		for v := range exchange.Exchange.ReadChannel(exchange.NotificationTGBot) {
			if v.Start {
				if v.Module == models.Notificator {
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
	}()

	// Получаем обновления от бота
	var updates tgbotapi.UpdatesChannel
	for {
		updates, err = getUpdates(bot, u)
		if updates != nil && err == nil {
			break
		}
		<-time.After(time.Second * 3)
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
