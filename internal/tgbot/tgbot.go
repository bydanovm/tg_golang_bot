package tgbot

import (
	"fmt"
	"os"
	"reflect"
	"time"

	_ "github.com/lib/pq"
	"github.com/mbydanov/tg_golang_bot/internal/coinmarketcup"
	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/models"
	"github.com/mbydanov/tg_golang_bot/internal/notifications"
	"github.com/mbydanov/tg_golang_bot/internal/services"
	"github.com/sirupsen/logrus"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// Создаем бота
func TelegramBot(statusRetriever chan models.StatusRetriever,
	notifTelegramIn chan models.StatusChannel) {
	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		services.Logging.Panic(err.Error())
	}

	// Устанавливаем время обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Функция отправки сообщения об ошибке из внешних сервисов
	go func(chatID int64) {
		for {
			// Отправляем сообщение об ошибке
			val, ok := <-statusRetriever
			if ok {
				if val.MsgError != nil {
					msg := tgbotapi.NewMessage(chatID, val.MsgError.Error())
					bot.Send(msg)
					services.Logging.Error(val.MsgError.Error())
				}
			}
		}
	}(786751823)

	// Функция получения сообщения от нотификатора
	go func() {
		for {
			val, ok := <-notifTelegramIn
			if ok {
				arr, ok := val.Data.([]notifications.NotificationsCCStruct)
				if ok {
					for _, v := range arr {
						msg := tgbotapi.NewMessage(int64(v.IdChat), v.Event)
						bot.Send(msg)
					}
				}
				if val.Error != nil {
					services.Logging.Error(val.Error.Error())
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
		if update.Message == nil {
			continue
		}
		// Проверяем есть ли пользователь в кеше или базе
		if _, ok := database.UsersCache[update.Message.From.ID]; !ok {
			// Пользователь не в кеше, ищем в БД, если не находим, то добавляем нового
			// Единственная точка входа, где пользователь может добавиться в БД
			user := database.Users{
				IdUsr:     update.Message.From.ID,
				TsUsr:     time.Now(),
				NameUsr:   update.Message.From.UserName,
				FirstName: update.Message.From.FirstName,
				LastName:  update.Message.From.LastName,
				LangCode:  update.Message.From.LanguageCode,
				IsBot:     update.Message.From.IsBot,
				IsBanned:  false,
				ChatIdUsr: update.Message.Chat.ID,
				IdLvlSec:  5}
			// Поиск с последующим добавлением
			if err := user.CheckUser(); err != nil {
				// Отправляем сообщение в лог об ошибке
				services.Logging.Warn(err.Error())
			}
			database.UsersCache[update.Message.From.ID] = user
		}
		// if err := services.CheckUser(update.Message); err != nil {
		// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		// 	bot.Send(msg)
		// }
		// Проверяем что от пользователя пришло именно текстовое сообщение
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			// Логируем запрос в лог
			services.Logging.WithFields(logrus.Fields{
				"userId":   database.UsersCache[update.Message.From.ID].IdUsr,
				"userName": database.UsersCache[update.Message.From.ID].NameUsr,
			}).Info(update.Message.Text)
			switch update.Message.Text {
			case "/start":
				// Отправлем сообщение
				msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr, "Hi, i'm a bot.")
				bot.Send(msg)
			case "/number_of_users":
				if os.Getenv("DB_SWITCH") == "on" {
					// Присваиваем количество пользователей использовавших бота в num переменную
					num, err := database.GetNumberOfUsers()
					if err != nil {
						//Отправлем сообщение
						msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr, "Database error.")
						bot.Send(msg)
					}

					// Создаем строку которая содержит колличество пользователей использовавших бота
					ans := fmt.Sprintf("%d peoples used me", num)

					// Отправлем сообщение
					msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr, ans)
					bot.Send(msg)
				} else {
					// Отправлем сообщение
					msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr, "Database not connected, so i can't say you how many peoples used me.")
					bot.Send(msg)
				}
			default:
				// Проверяем лимит на запросы конкретного пользователя
				message := coinmarketcup.GetLatest(update.Message.Text)
				// message := wiki.WikipediaGET(update.Message.Text)
				if os.Getenv("DB_SWITCH") == "on" {
					// Отправляем username, chat_id, message, answer в БД
					if err := database.CollectData(database.UsersCache[update.Message.From.ID].NameUsr,
						database.UsersCache[update.Message.From.ID].ChatIdUsr, update.Message.Text, message); err != nil {

						// Отправлем сообщение
						msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr, "Database error, but bot still working.")
						bot.Send(msg)
					}
				}

				// Проходим через срез и отправляем каждый элемент пользователю
				for _, val := range message {
					// Логируем ответ бота
					services.Logging.WithFields(logrus.Fields{
						"userId":   database.UsersCache[update.Message.From.ID].IdUsr,
						"userName": database.UsersCache[update.Message.From.ID].NameUsr,
					}).Info(val)
					// Отправлем сообщение
					msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr, val)
					bot.Send(msg)
				}
			}
		} else {
			// Отправлем сообщение
			msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr, "Use the words for search.")
			bot.Send(msg)
		}
	}
}
