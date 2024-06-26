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

// Создаем бота
func TelegramBot(chanModules chan models.StatusChannel) {
	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		services.Logging.Panic(err.Error())
	}

	keyboardBot := initMenu()
	// Инициализация меню

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
		// Авторизация пользователя
		userId, err := checkAuthUser(bot, &update)
		if err != nil {
			services.Logging.WithFields(logrus.Fields{
				"module": "tgbot",
			}).Error(err.Error())
		}
		if update.Message == nil && update.InlineQuery != nil {
			// Обработка inline-режима
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
					"userId":   database.UsersCache.GetUserId(update.InlineQuery.From.ID),
					"userName": database.UsersCache.GetUserName(update.InlineQuery.From.ID),
					"type":     "inline",
				}).Info(err)
			}
		} else {
			var command = ""
			var param = ""
			message := []string{}
			if update.Message != nil {
				command = update.Message.Command()
				param = update.Message.CommandArguments()
				if command != "" {
					services.Logging.WithFields(logrus.Fields{
						"userId":   database.UsersCache.GetUserId(update.Message.From.ID),
						"userName": database.UsersCache.GetUserName(update.Message.From.ID),
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
						if userId := database.UsersCache.GetUserId(update.Message.From.ID); userId == 0 {
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
							} else {
								// Кешируем добавленного пользователя
								if err := database.UsersCache.CheckCache(update.Message.From.ID); err != nil {
									services.Logging.Warn(err.Error())
								}
							}
						}
						// Отправлем приветственное сообщение
						ans := "Привет! Я - " + os.Getenv("BOT_NAME") + " помогу тебе знать актуальную информацию по криптовалюте\n" +
							"Используй клавиатуру ниже, чтобы узнать интересующую информацию.\n"
						message = append(message,
							ans)
						msg := tgbotapi.NewMessage(database.UsersCache.GetChatId(update.Message.From.ID),
							ans)
						msg.ReplyMarkup = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkup(), 2)

						bot.Send(msg)
					default:
						prevMenu := database.UsersCache.GetPrevMenu(userId)
						// Предыдущее меню - выбор крипты для отслеживания
						if prevMenu == SetNotifCrypto {

						} else {
							ans := "Команда /" + command + " не найдена.\n" +
								"Воспользуйся командой /" + Start + " для знакомства со мной."
							message = append(message, ans)

							// Отправлем сообщение
							msg := tgbotapi.NewMessage(database.UsersCache.GetChatId(update.Message.From.ID),
								ans)
							bot.Send(msg)
						}
					}
					// Логируем ответ бота на команды
					services.Logging.WithFields(logrus.Fields{
						"userId":   database.UsersCache.GetUserId(update.Message.From.ID),
						"userName": database.UsersCache.GetUserName(update.Message.From.ID),
						"type":     "answer",
					}).Info(message)
				} else {
					// Обработка сообщений с ForceReply
					if update.Message.ReplyToMessage != nil {
						switch update.Message.ReplyToMessage.Text {
						case "Введите криптовалюту для отслеживания":
							ans := "Выбрана криптовалюта: " + update.Message.Text + "\nВыберите критерий"
							msg := tgbotapi.NewMessage(database.UsersCache.GetChatId(int(update.Message.Chat.ID)),
								ans)

							keyboard := MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotifCriterion), 2)
							msg.ReplyMarkup = keyboard

							bot.Send(msg)
						}
					} else {
						if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
							ans := update.Message.Text + " что такое, я такого не знаю.\n" +
								"Воспользуйся командой /" + Start + " для знакомства со мной."
							message = append(message, ans)
							msg := tgbotapi.NewMessage(database.UsersCache.GetChatId(update.Message.From.ID),
								ans)
							bot.Send(msg)
						} else {
							ans := update.Message.Text + " что такое, я такого не знаю.\n" +
								"Воспользуйся командой /" + Start + " для знакомства со мной."
							message = append(message, ans)
							msg := tgbotapi.NewMessage(update.Message.Chat.ID,
								ans)
							bot.Send(msg)
						}
					}
				}
				// Собираем статистику в базу
				if err := database.CollectData(database.UsersCache.GetUserName(update.Message.From.ID),
					database.UsersCache.GetChatId(update.Message.From.ID), update.Message.Text, message); err != nil {
					services.Logging.WithFields(logrus.Fields{
						"userId":   database.UsersCache.GetUserId(update.Message.From.ID),
						"userName": database.UsersCache.GetUserName(update.Message.From.ID),
					}).Error("tgbot:", err.Error())
				}
			} else {
				// Обработка callback
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
						// Получение InlineKeyboard со списком крипты топ 10?
						if callBackData[0] == GetCrypto {
							ans := "Список крипты\n"

							msg := tgbotapi.NewEditMessageText(database.UsersCache.GetChatId(int(update.CallbackQuery.Message.Chat.ID)),
								update.CallbackQuery.Message.MessageID, ans)
							keyboard := tgbotapi.InlineKeyboardMarkup{}
							top10cur, err := database.DCCache.GetTop10Cache()
							if err != nil {
								services.Logging.WithFields(logrus.Fields{
									"userId":   database.UsersCache.GetUserId(int(update.CallbackQuery.Message.Chat.ID)),
									"userName": database.UsersCache.GetUserName(int(update.CallbackQuery.Message.Chat.ID)),
								}).Error(err)
							}
							var row []tgbotapi.InlineKeyboardButton
							for k, v := range top10cur {
								btn := tgbotapi.NewInlineKeyboardButtonData(v.CryptoName, GetCrypto+"_"+v.CryptoName)
								row = append(row, btn)
								// Делим на N строк по 5 элементов
								if (k+1)%5 == 0 {
									keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
									row = nil
								}
							}
							row = append(row, tgbotapi.NewInlineKeyboardButtonData("Назад", Start))
							row = append(row, tgbotapi.NewInlineKeyboardButtonData("Еще", GetCrypto+"_next"))
							keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
							msg.ReplyMarkup = &keyboard
							bot.Send(msg)
						} else if callBackData[0] == Start {
							if userId := database.UsersCache.GetChatId(userId); userId == 0 {
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
								} else {
									// Кешируем добавленного пользователя
									if err := database.UsersCache.CheckCache(update.Message.From.ID); err != nil {
										services.Logging.Warn(err.Error())
									}
								}
							}
							// Отправлем приветственное сообщение
							ans := "Привет! Я - " + os.Getenv("BOT_NAME") + " помогу тебе знать актуальную информацию по криптовалюте\n" +
								"Используй клавиатуру ниже, чтобы узнать интересующую информацию.\n"
							// msg := tgbotapi.NewMessage(database.UsersCache.GetChatId(userId),
							// 	ans)
							msg := tgbotapi.NewEditMessageText(database.UsersCache.GetChatId(int(update.CallbackQuery.Message.Chat.ID)),
								update.CallbackQuery.Message.MessageID, ans)

							keyboard := MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkup(), 2)
							msg.ReplyMarkup = &keyboard

							bot.Send(msg)

						} else if callBackData[0] == SetNotif {
							ans := "Здесь можно завести оповещения\n"
							msg := tgbotapi.NewEditMessageText(database.UsersCache.GetChatId(int(update.CallbackQuery.Message.Chat.ID)),
								update.CallbackQuery.Message.MessageID, ans)

							keyboard := MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotif), 2)
							msg.ReplyMarkup = &keyboard

							bot.Send(msg)
						}
					} else if callBackData[0] == GetCrypto {
						message = coinmarketcup.GetLatest(callBackData[1])
						// Проходим через срез и отправляем каждый элемент пользователю
						for _, val := range message {
							// Логируем ответ бота
							services.Logging.WithFields(logrus.Fields{
								"userId":   database.UsersCache.GetUserId(int(update.CallbackQuery.Message.Chat.ID)),
								"userName": database.UsersCache.GetUserName(int(update.CallbackQuery.Message.Chat.ID)),
								"type":     "callback",
								"command":  GetCrypto,
								"currency": callBackData[1],
							}).Info(val)
							// Отправлем сообщение
							msg := tgbotapi.NewMessage(database.UsersCache.GetChatId(int(update.CallbackQuery.Message.Chat.ID)), val)
							bot.Send(msg)
						}
					} else if callBackData[0] == `next` && callBackData[1] == GetCrypto {
						msg := tgbotapi.NewMessage(database.UsersCache.GetChatId(int(update.CallbackQuery.Message.Chat.ID)),
							"Введите свои криптовалюты")
						bot.Send(msg)
					} else if callBackData[0] == SetNotif {
						// Запоминаем текущее меню
						database.UsersCache.SetPrevMenu(userId, update.CallbackQuery.Data)

						// Оповещения
						switch update.CallbackQuery.Data {
						case SetNotifCrypto:
							ans := "Введите криптовалюту для отслеживания"
							msg := tgbotapi.NewMessage(database.UsersCache.GetChatId(int(update.CallbackQuery.Message.Chat.ID)),
								ans)

							msg.ReplyMarkup = tgbotapi.ForceReply{
								ForceReply: true,
							}

							bot.Send(msg)

						case SetNotifCriterionMore:
							ans := "Выбран критерий \"Больше\"\n"
							msg := tgbotapi.NewEditMessageText(database.UsersCache.GetChatId(int(update.CallbackQuery.Message.Chat.ID)),
								update.CallbackQuery.Message.MessageID, ans)
							bot.Send(msg)
						case SetNotifCriterionLess:
							ans := "Выбран критерий \"Меньше\"\n"
							msg := tgbotapi.NewEditMessageText(database.UsersCache.GetChatId(int(update.CallbackQuery.Message.Chat.ID)),
								update.CallbackQuery.Message.MessageID, ans)
							bot.Send(msg)

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
