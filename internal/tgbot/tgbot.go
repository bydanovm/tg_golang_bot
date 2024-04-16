package tgbot

import (
	"fmt"
	"os"
	"reflect"
	"strings"
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

const (
	// Comands
	Start         string = "start"           // Начало
	NumberOfUsers string = "number_of_users" // Получить количество активных пользователей
	GetCrypto     string = "getcrypto"       // Получить актуальную информацию по криптовалюте
	SetNotif      string = "setnotif"        // Установить уведомления по изменению цены криптовалюты
)

// Создаем бота
func TelegramBot(chanModules chan models.StatusChannel) {
	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		services.Logging.Panic(err.Error())
	}

	// Устанавливаем время обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// // Функция отправки сообщения об ошибке из внешних сервисов
	// go func(chatID int64) {
	// 	for {
	// 		// Отправляем сообщение об ошибке
	// 		val, ok := <-statusRetriever
	// 		if ok {
	// 			if val.MsgError != nil {
	// 				msg := tgbotapi.NewMessage(chatID, val.MsgError.Error())
	// 				bot.Send(msg)
	// 				services.Logging.Error(val.MsgError.Error())
	// 			}
	// 		}
	// 	}
	// }(786751823)

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
		if update.Message != nil {
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
		}
		if update.Message == nil && update.InlineQuery != nil {
			// Обработка inline-режима
			// q := update.InlineQuery.Query
			// filteredQ := Filter()
		} else {
			var command = ""
			var param = ""
			message := []string{}
			if update.Message != nil {
				command = update.Message.Command()
				param = update.Message.CommandArguments()
				if command != "" {
					services.Logging.WithFields(logrus.Fields{
						"userId":   database.UsersCache[update.Message.From.ID].IdUsr,
						"userName": database.UsersCache[update.Message.From.ID].NameUsr,
						"type":     "command",
						"args":     param,
					}).Info(command)
					// Обработка команд
					// start - Начало
					// number_of_users - Получить количество активных пользователей
					// getcrypto - Получить актуальную информацию по криптовалюте
					// setnotif - Установить уведомления по изменению цены криптовалюты
					switch command {
					case Start:
						// Отправлем приветственное сообщение
						ans := "Привет, я бот"
						message = append(message,
							ans)
						msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr,
							ans)
						bot.Send(msg)
					case NumberOfUsers:
						// Создаем строку которая содержит колличество пользователей использовавших бота
						// Берем из кеша
						ans := fmt.Sprintf("%d пользователь использует бота", len(database.UsersCache))
						message = append(message, ans)
						// Отправлем сообщение
						msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr,
							ans)
						bot.Send(msg)
					case GetCrypto:
						if param != "" {
							message = coinmarketcup.GetLatest(param)
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
						} else {
							msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr, "Выберите криптовалюту")

							keyboard := tgbotapi.InlineKeyboardMarkup{}
							top10cur, err := database.DCCache.GetTop10Cache()
							if err != nil {
								services.Logging.WithFields(logrus.Fields{
									"userId":   database.UsersCache[update.Message.From.ID].IdUsr,
									"userName": database.UsersCache[update.Message.From.ID].NameUsr,
								}).Error(err)
							}
							for _, v := range top10cur {
								var row []tgbotapi.InlineKeyboardButton
								btn := tgbotapi.NewInlineKeyboardButtonData(v.CryptoName, GetCrypto+"_"+v.CryptoName)
								row = append(row, btn)
								keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
							}
							msg.ReplyMarkup = keyboard
							bot.Send(msg)
						}

					default:
						ans := "Команда /" + command + " не найдена"
						message = append(message, ans)

						// Отправлем сообщение
						msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr,
							ans)
						bot.Send(msg)

					}
					// Логируем ответ бота на команды
					services.Logging.WithFields(logrus.Fields{
						"userId":   database.UsersCache[update.Message.From.ID].IdUsr,
						"userName": database.UsersCache[update.Message.From.ID].NameUsr,
						"type":     "answer",
					}).Info(message)
				} else {
					// Обработка обычных сообщений
					// Проверяем что от пользователя пришло именно текстовое сообщение
					if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
						// Логируем запрос в лог
						services.Logging.WithFields(logrus.Fields{
							"userId":   database.UsersCache[update.Message.From.ID].IdUsr,
							"userName": database.UsersCache[update.Message.From.ID].NameUsr,
						}).Info(update.Message.Text)
						// Проверяем лимит на запросы конкретного пользователя
						message := coinmarketcup.GetLatest(update.Message.Text)
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
					} else {
						// Отправлем сообщение
						msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr, "Use the words for search.")
						bot.Send(msg)
					}
				}
				// Собираем статистику в базу
				if err := database.CollectData(database.UsersCache[update.Message.From.ID].NameUsr,
					database.UsersCache[update.Message.From.ID].ChatIdUsr, update.Message.Text, message); err != nil {
					services.Logging.WithFields(logrus.Fields{
						"userId":   database.UsersCache[update.Message.From.ID].IdUsr,
						"userName": database.UsersCache[update.Message.From.ID].NameUsr,
					}).Error("tgbot:", err.Error())
				}
			} else {
				// Обработка callback
				if update.CallbackQuery != nil {
					// Проверка, что пользователь есть в базе (кеше)
					if _, ok := database.UsersCache[update.CallbackQuery.From.ID]; !ok {
						ans := "Мы не нашли Вас в базе. Пожалуйста, воспользуйте сначала командой /start для регистрации"
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ans)
						bot.Send(msg)
						services.Logging.WithFields(logrus.Fields{
							"userId":   update.CallbackQuery.Message.Chat.ID,
							"userName": update.CallbackQuery.Message.From.UserName,
							"type":     "callback",
						}).Error("tgbot:", ans)
					}
					// Проверка команд
					// Разберем data callback по структуре command_cryptocur
					callBackData := strings.Split(update.CallbackQuery.Data, "_")
					if callBackData[0] == GetCrypto {
						message = coinmarketcup.GetLatest(callBackData[1])
						// Проходим через срез и отправляем каждый элемент пользователю
						for _, val := range message {
							// Логируем ответ бота
							services.Logging.WithFields(logrus.Fields{
								"userId":   database.UsersCache[int(update.CallbackQuery.Message.Chat.ID)].IdUsr,
								"userName": database.UsersCache[int(update.CallbackQuery.Message.Chat.ID)].NameUsr,
								"type":     "callback",
								"command":  GetCrypto,
								"currency": callBackData[1],
							}).Info(val)
							// Отправлем сообщение
							msg := tgbotapi.NewMessage(database.UsersCache[int(update.CallbackQuery.Message.Chat.ID)].ChatIdUsr, val)
							bot.Send(msg)
						}
					}
				}
			}
		}
	}
}
