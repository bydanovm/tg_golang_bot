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
		if update.Message == nil && update.InlineQuery != nil {
			// Обработка inline-режима
			// Кешируем пользователя
			if err := database.UsersCache.CheckCache(update.InlineQuery.From.ID); err != nil {
				services.Logging.WithFields(logrus.Fields{
					"module":   "tgbot",
					"userId":   update.InlineQuery.From.ID,
					"userName": update.InlineQuery.From.UserName,
				}).Error(err.Error())
			}
			query := update.InlineQuery.Query
			filteredCrypto := Filter(database.DCCache.GetAllCache(), func(dc database.DictCrypto) bool {
				return strings.Index(strings.ToUpper(dc.CryptoName), strings.ToUpper(query)) >= 0
			})
			// Логирование
			services.Logging.WithFields(logrus.Fields{
				"userId":   update.InlineQuery.From.ID,
				"userName": update.InlineQuery.From.UserName,
				"query":    query,
				"filtered": filteredCrypto,
			}).Info()

			var articles []interface{}
			if len(filteredCrypto) == 0 {
				// Если ничего не найдено - выводим топ 10
				top10cur, err := database.DCCache.GetTop10Cache()
				if err != nil {
					services.Logging.WithFields(logrus.Fields{
						"userId":   update.InlineQuery.From.ID,
						"userName": update.InlineQuery.From.UserName,
					}).Error(err)
				}
				for _, v := range top10cur {
					text := fmt.Sprintf("Криптовалюта: %s\nЦена: %.9f %s",
						v.CryptoName,
						v.CryptoLastPrice,
						"USD",
					)
					msg := tgbotapi.NewInlineQueryResultArticleMarkdown(v.CryptoName, v.CryptoName, text)
					articles = append(articles, msg)
				}
			} else {
				for k, v := range filteredCrypto {
					text := fmt.Sprintf("Криптовалюта: %s\nЦена: %.9f %s",
						v.CryptoName,
						v.CryptoLastPrice,
						"USD",
					)
					msg := tgbotapi.NewInlineQueryResultArticleMarkdown(v.CryptoName, v.CryptoName, text)
					articles = append(articles, msg)
					if k >= 10 {
						break
					}
				}

			}
			inlineConfig := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     0,
				Results:       articles,
			}
			_, err = bot.AnswerInlineQuery(inlineConfig)
			if err != nil {
				services.Logging.WithFields(logrus.Fields{
					"userId":   database.UsersCache[update.InlineQuery.From.ID].IdUsr,
					"userName": database.UsersCache[update.InlineQuery.From.ID].NameUsr,
					"type":     "inline",
				}).Info(err)
			}
		} else {
			var command = ""
			var param = ""
			message := []string{}
			if update.Message != nil {
				// Кешируем пользователя
				if err := database.UsersCache.CheckCache(update.Message.From.ID); err != nil {
					services.Logging.WithFields(logrus.Fields{
						"module":   "tgbot",
						"userId":   update.Message.From.ID,
						"userName": update.Message.From.UserName,
					}).Error(err.Error())
				}

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
						// Отправлем приветственное сообщение
						ans := "Привет! Давай я немного расскажу о себе.\n" +
							"Я умею выдавать цену по интересующей тебя криптовалюте. Для этого используй команду /" + GetCrypto +
							" плюс мнемоника криптовалюты. Например: /" + GetCrypto + " BTC. Если ты не знаешь к какой криптовалюте обратиться," +
							" используй просто команду /" + GetCrypto + ", она выведет тебе список 10 самых используемых у меня валют.\n" +
							"Ещё ты можешь обратиться ко мне из любого чата командой @" + os.Getenv("BOT_NAME") +
							" и я выведу тебе 10 любых криптовалют с их ценами. Так же я могу выдать тебе отсортированный список " +
							"криптовалют, для этого после команды @" + os.Getenv("BOT_NAME") + " достаточно ввести любую букву из " +
							"мнемоники криптовалюты.\n" +
							"Еще я хочу научиться опрашивать много много разных API и показывать тебе уведомления о изменении цен, " +
							"но пока я еще маленький и еще многому учусь.\n"
						message = append(message,
							ans)
						msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr,
							ans)
						bot.Send(msg)
					case NumberOfUsers:
						if _, ok := database.UsersCache[update.Message.From.ID]; !ok {
							ans := fmt.Sprintf("Я тебя не знаю, давай сначала познакомимся.\nВведи команду /%s", Start)
							message = append(message, ans)
							// Отправлем сообщение
							msg := tgbotapi.NewMessage(update.Message.Chat.ID,
								ans)
							bot.Send(msg)
						} else {
							// Создаем строку которая содержит колличество пользователей использовавших бота
							// Берем из кеша
							ans := fmt.Sprintf("%d пользователь использует бота", len(database.UsersCache))
							message = append(message, ans)
							// Отправлем сообщение
							msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr,
								ans)
							bot.Send(msg)
						}
					case GetCrypto:
						if _, ok := database.UsersCache[update.Message.From.ID]; !ok {
							ans := fmt.Sprintf("Я тебя не знаю, давай сначала познакомимся.\nВведи команду /%s", Start)
							message = append(message, ans)
							// Отправлем сообщение
							msg := tgbotapi.NewMessage(update.Message.Chat.ID,
								ans)
							bot.Send(msg)
						} else {

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

						}

					default:
						if _, ok := database.UsersCache[update.Message.From.ID]; !ok {
							ans := fmt.Sprintf("Я тебя не знаю, давай сначала познакомимся.\nВведи команду /%s", Start)
							message = append(message, ans)
							// Отправлем сообщение
							msg := tgbotapi.NewMessage(update.Message.Chat.ID,
								ans)
							bot.Send(msg)
						} else {
							ans := "Команда /" + command + " не найдена.\n" +
								"Воспользуйся командой /" + Start + " для знакомства со мной."
							message = append(message, ans)

							// Отправлем сообщение
							msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr,
								ans)
							bot.Send(msg)
						}

					}
					// Логируем ответ бота на команды
					services.Logging.WithFields(logrus.Fields{
						"userId":   database.UsersCache[update.Message.From.ID].IdUsr,
						"userName": database.UsersCache[update.Message.From.ID].NameUsr,
						"type":     "answer",
					}).Info(message)
				} else {
					if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
						if _, ok := database.UsersCache[update.Message.From.ID]; !ok {
							ans := fmt.Sprintf("Я тебя не знаю, давай сначала познакомимся.\nВведи команду /%s", Start)
							message = append(message, ans)
							msg := tgbotapi.NewMessage(update.Message.Chat.ID,
								ans)
							bot.Send(msg)
						} else {
							ans := update.Message.Text + " что такое, я такого не знаю.\n" +
								"Воспользуйся командой /" + Start + " для знакомства со мной."
							message = append(message, ans)
							msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr,
								ans)
							bot.Send(msg)
						}
					} else {
						ans := update.Message.Text + " что такое, я такого не знаю.\n" +
							"Воспользуйся командой /" + Start + " для знакомства со мной."
						message = append(message, ans)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID,
							ans)
						bot.Send(msg)
					}
					// Обработка обычных сообщений
					// Проверяем что от пользователя пришло именно текстовое сообщение
					// if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
					// 	// Логируем запрос в лог
					// 	services.Logging.WithFields(logrus.Fields{
					// 		"userId":   database.UsersCache[update.Message.From.ID].IdUsr,
					// 		"userName": database.UsersCache[update.Message.From.ID].NameUsr,
					// 	}).Info(update.Message.Text)
					// 	// Проверяем лимит на запросы конкретного пользователя
					// 	message := coinmarketcup.GetLatest(update.Message.Text)
					// 	if os.Getenv("DB_SWITCH") == "on" {
					// 		// Отправляем username, chat_id, message, answer в БД
					// 		if err := database.CollectData(database.UsersCache[update.Message.From.ID].NameUsr,
					// 			database.UsersCache[update.Message.From.ID].ChatIdUsr, update.Message.Text, message); err != nil {

					// 			// Отправлем сообщение
					// 			msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr, "Database error, but bot still working.")
					// 			bot.Send(msg)
					// 		}
					// 	}

					// 	// Проходим через срез и отправляем каждый элемент пользователю
					// 	for _, val := range message {
					// 		// Логируем ответ бота
					// 		services.Logging.WithFields(logrus.Fields{
					// 			"userId":   database.UsersCache[update.Message.From.ID].IdUsr,
					// 			"userName": database.UsersCache[update.Message.From.ID].NameUsr,
					// 		}).Info(val)
					// 		// Отправлем сообщение
					// 		msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr, val)
					// 		bot.Send(msg)
					// 	}
					// } else {
					// 	// Отправлем сообщение
					// 	msg := tgbotapi.NewMessage(database.UsersCache[update.Message.From.ID].ChatIdUsr, "Use the words for search.")
					// 	bot.Send(msg)
					// }
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
						ans := "Я не нашел тебя в своей базе. Пожалуйста, воспользуйте сначала командой /start для знакомства."
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ans)
						bot.Send(msg)
						services.Logging.WithFields(logrus.Fields{
							"userId":   update.CallbackQuery.Message.Chat.ID,
							"userName": update.CallbackQuery.Message.From.UserName,
							"type":     "callback",
						}).Error("tgbot:", ans)
					} else {
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
