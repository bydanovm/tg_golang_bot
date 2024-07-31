package tgbot

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	ChooseCrypto string = "Выберите или введите криптовалюту для отслеживания\n"
	ChooseSum    string = "Выберите или введите сумму для отслеживания\n"
	EnterCrypto  string = "Введите криптовалюту для отслеживания\n"
	EnterSum     string = "Введите сумму для отслеживания\n"
)

type PriceInfo struct {
	Koeff int
	Price float32
}

// Функция обработки меню по "Оповещениям"
func menuNotification(update *tgbotapi.Update, keyboardBot *tgBotMenu) (msg interface{}) {
	// var ans string
	// var keyboard tgbotapi.InlineKeyboardMarkup
	// Обработка callback
	// if update.CallbackQuery != nil {
	// 	// Проверка команд
	// 	// Вынести в общее ответ по коллбеку
	// 	// callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	// 	// callback.ShowAlert = true
	// 	// if _, err := bot.AnswerCallbackQuery(callback); err != nil {
	// 	// 	services.Logging.WithFields(logrus.Fields{
	// 	// 		"userId":   update.CallbackQuery.Message.Chat.ID,
	// 	// 		"userName": update.CallbackQuery.Message.From.UserName,
	// 	// 		"type":     "callback_answer",
	// 	// 		"command":  update.CallbackQuery.Data,
	// 	// 	}).Error()
	// 	// }

	// 	// Разберем data callback по структуре command_subcmd_addcmd
	// 	callBackData := strings.Split(update.CallbackQuery.Data, "_")
	// 	// Главное меню оповещений
	// 	if len(callBackData) == 1 {
	// 		switch callBackData[0] {
	// 		case SetNotif:
	// 			keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotif), 2)
	// 			// Текст взять из базы (нужен справочник)
	// 			ans = "Здесь можно завести оповещения\n"

	// 		case GetNotif:
	// 			ans = "Текущие оповещения\n"
	// 			// Вывести от новых к старым в формате Валюта - Значение
	// 			trackings, _ := caching.GetCacheRecordsKeyChain(caching.TrackingCache, update.CallbackQuery.From.ID, true)
	// 			// Создание списка кнопок
	// 			listButtons := make([]buttonInfo, 0, 10)
	// 			for _, v := range trackings {
	// 				// infoCurrency, _ := database.DCCache.GetCache(v.DctCrpId)
	// 				infoCurrency, _ := caching.GetCacheByIdxInMap(caching.CryptoCache, v.DctCrpId, 0)
	// 				listButtons = append(listButtons, buttonInfo{infoCurrency.CryptoName + " - " + fmt.Sprintf(FormatFloatToString(v.ValTrkCrp), v.ValTrkCrp) + " $", GetNotifId + "_" + fmt.Sprintf("%d", v.IdTrkCrp)})
	// 			}
	// 			keyboard = ConvertToButtonInlineKeyboard(listButtons, GetNotif, 3)
	// 		}
	// 		msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
	// 			update.CallbackQuery.Message.MessageID, ans)

	// 		msg_t.ReplyMarkup = &keyboard
	// 		msg = msg_t

	// 	} else if len(callBackData) == 2 {
	// 		if callBackData[0] == SetNotif {
	// 			// Оповещения 2 уровень
	// 			switch update.CallbackQuery.Data {
	// 			case SetNotifCrypto: // Выбрать крипту
	// 				// Текст взять из базы (нужен справочник)
	// 				ans = ChooseCrypto
	// 				// Выбор ТОП-10 криптовалют
	// 				offset := 10
	// 				keyboard = GetCryptoListOffset(offset, SetNotifCrypto, SetNotif)

	// 			case SetNotifPrice: // Установить цену
	// 				// Текст взять из базы (нужен справочник)
	// 				ans = ChooseSum

	// 			}
	// 		} else if callBackData[0] == GetNotifId {
	// 			// Мои оповещения -> выбрано Оповещение
	// 			id, _ := strconv.Atoi(callBackData[1])
	// 			infoTracking, _ := caching.GetCacheByIdxInMap(caching.TrackingCache, id, 0)
	// 			infoCurrency, _ := caching.GetCacheByIdxInMap(caching.CryptoCache, infoTracking.DctCrpId, 0)
	// 			// Считвание типа триггера
	// 			ans = fmt.Sprintf("Выбрано отслеживание по %s\nТекущая стоимость 1 %s = "+FormatFloatToString(infoTracking.ValTrkCrp)+" $\nТип триггера: %d\nСрабатывание триггера: "+FormatFloatToString(infoTracking.ValTrkCrp)+"\n", infoCurrency.CryptoName, infoCurrency.CryptoName, infoCurrency.CryptoLastPrice, infoTracking.TypTrkCrpId, infoTracking.ValTrkCrp)

	// 			keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(GetNotifId), 1)
	// 		}

	// 		msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
	// 			update.CallbackQuery.Message.MessageID, ans)
	// 		msg_t.ReplyMarkup = &keyboard

	// 		msg = msg_t

	// 	} else if len(callBackData) == 3 && callBackData[0] == SetNotif {
	// 		// Оповещения 3 уровень
	// 		switch update.CallbackQuery.Data {
	// 		// Ввод своей КВ
	// 		case SetNotifCryptoEnter:
	// 			ans = EnterCrypto
	// 			msg_t := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
	// 				ans)
	// 			msg_t.ReplyMarkup = tgbotapi.ForceReply{
	// 				ForceReply: true,
	// 			}
	// 			msg = msg_t

	// 		// Ввод своей суммы отслеживания
	// 		case SetNotifPriceEnter:
	// 			ans = EnterSum
	// 			msg_t := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
	// 				ans)
	// 			msg_t.ReplyMarkup = tgbotapi.ForceReply{
	// 				ForceReply: true,
	// 			}
	// 			msg = msg_t

	// 		case SetNotifCriterionMore:
	// 			ans = "Выбран критерий \"Больше\"\n"
	// 			SetNotifCh.SetCriterion(int(update.CallbackQuery.Message.Chat.ID), "+")

	// 			// msg = tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
	// 			// 	update.CallbackQuery.Message.MessageID, ans)

	// 		case SetNotifCriterionLess:
	// 			ans = "Выбран критерий \"Меньше\"\n"
	// 			SetNotifCh.SetCriterion(int(update.CallbackQuery.Message.Chat.ID), "-")

	// 		// msg = tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
	// 		// 	update.CallbackQuery.Message.MessageID, ans)
	// 		case SetNotifPriceYes:
	// 			// Нажата ДА на последнем этапе, возврат в начало
	// 			// Определить КВ по мнемонике
	// 			idCrpt := database.DCCacheKeys.GetCacheIdByName(SetNotifCh.GetCrypto(int(update.CallbackQuery.Message.Chat.ID)))
	// 			if idCrpt == 0 {
	// 				ans += fmt.Sprintf("Криптовалюта %s не найдена.\nИсправьте команду и повторите запрос\n", SetNotifCh.GetCrypto(int(update.CallbackQuery.Message.Chat.ID)))
	// 			}
	// 			// Найти Тип отслеживания
	// 			idType := 0
	// 			if SetNotifCh.GetCriterion(int(update.CallbackQuery.Message.Chat.ID)) == "+" {
	// 				idType = database.TypeTCCacheKeys.GetCacheIdByName("RAISE_V")
	// 			} else if SetNotifCh.GetCriterion(int(update.CallbackQuery.Message.Chat.ID)) == "-" {
	// 				idType = database.TypeTCCacheKeys.GetCacheIdByName("FALL_V")
	// 			} else {
	// 				ans += "Неверный тип отслеживания\nИсправьте команду и повторите запрос\n"
	// 			}
	// 			// Установка лимита
	// 			limit := database.Limits{}
	// 			if ans == "" {
	// 				limit = database.Limits{
	// 					IdLmt:       database.LmtCache.GetCacheLastId(),
	// 					ValAvailLmt: database.LmtCacheKeys["LMT003"].StdValLmt,
	// 					ActiveLmt:   true,
	// 					UserId:      update.CallbackQuery.From.ID,
	// 					LtmDctId:    database.LmtCacheKeys["LMT003"].IdLmtDct,
	// 				}
	// 				if err := limit.SetLimit(); err != nil {
	// 					ans += fmt.Sprintf("tgbot:%s\n", err.Error())
	// 				}
	// 			}

	// 			tracking := database.TrackingCrypto{
	// 				IdTrkCrp:    database.TCCache.GetCacheLastId(),
	// 				DctCrpId:    idCrpt,
	// 				TypTrkCrpId: idType,
	// 				LmtId:       limit.IdLmt,
	// 				UserId:      update.CallbackQuery.From.ID,
	// 				ValTrkCrp:   SetNotifCh.GetPrice(int(update.CallbackQuery.Message.Chat.ID)),
	// 				OnTrkCrp:    true,
	// 			}
	// 			if err := tracking.SetTracking(); err != nil {
	// 				ans += fmt.Sprintf("tgbot:%s\n", err.Error())
	// 			} else {
	// 				ans += fmt.Sprintf("Отслеживание по криптовалюте %s успешно добавлено\n", SetNotifCh.GetCrypto(int(update.CallbackQuery.Message.Chat.ID)))
	// 			}
	// 		case SetNotifPriceNo:
	// 			// Нажата НЕТ на последнем этапе, возврат в начало
	// 			// Текст взять из базы (нужен справочник)
	// 			ans = "Отслеживание не сохранено\n"

	// 		default:
	// 			// Случай, когда пришло 3 аргумента и выбрана КВ через кнопки
	// 			if callBackData[1] == Crypto {
	// 				// Сохранить КВ в мапу
	// 				// Переход к выбору критерия
	// 				// keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotifCriterion), 2)
	// 				// // Запись в кеш выбранной КВ
	// 				// SetNotifCh.SetCrypto(int(update.CallbackQuery.Message.Chat.ID), callBackData[2])

	// 				ans = "Выбрана криптовалюта: " + callBackData[2] + "\n"
	// 				// msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
	// 				// 	update.CallbackQuery.Message.MessageID, ans)
	// 				// msg_t.ReplyMarkup = &keyboard
	// 				// msg = msg_t

	// 				SetNotifCh.SetCrypto(int(update.CallbackQuery.Message.Chat.ID), callBackData[2])

	// 				ans += ChooseSum

	// 				crypto := SetNotifCh.GetCrypto(int(update.CallbackQuery.Message.Chat.ID))
	// 				cryptos, _ := coinmarketcup.GetLatestStruct(crypto)
	// 				isFind := false
	// 				prices := []PriceInfo{}

	// 				if len(cryptos) == 1 {
	// 					if cryptos[0].Find {
	// 						isFind = true
	// 						// Нужно вычислять количество знаков динамически
	// 						prices = append(prices, PriceInfo{1, cryptos[0].Crypto.CryptoLastPrice * 1.01})
	// 						prices = append(prices, PriceInfo{5, cryptos[0].Crypto.CryptoLastPrice * 1.05})
	// 						prices = append(prices, PriceInfo{10, cryptos[0].Crypto.CryptoLastPrice * 1.1})
	// 						prices = append(prices, PriceInfo{-1, cryptos[0].Crypto.CryptoLastPrice * 0.99})
	// 						prices = append(prices, PriceInfo{-5, cryptos[0].Crypto.CryptoLastPrice * 0.95})
	// 						prices = append(prices, PriceInfo{-10, cryptos[0].Crypto.CryptoLastPrice * 0.9})
	// 					} else {
	// 						ans += fmt.Sprintf("%s %s", cryptos[0].Crypto.CryptoName, "не найдена в базе")
	// 					}
	// 				}

	// 				var row []tgbotapi.InlineKeyboardButton
	// 				if isFind {
	// 					for k, v := range prices {
	// 						price := fmt.Sprintf("%f", v.Price)
	// 						n := fmt.Sprintf(FormatFloatToString(v.Price)+" (%+d%%)", v.Price, v.Koeff)
	// 						btn := tgbotapi.NewInlineKeyboardButtonData(n, SetNotifPrice+"_"+price)
	// 						row = append(row, btn)
	// 						// Делим на N строк по 3 элемента
	// 						if (k+1)%3 == 0 {
	// 							keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	// 							row = nil
	// 						}
	// 					}
	// 				}
	// 				row = append(row, tgbotapi.NewInlineKeyboardButtonData("Назад", SetNotifCrypto))
	// 				row = append(row, tgbotapi.NewInlineKeyboardButtonData("Ввод", SetNotifPriceEnter))
	// 				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	// 				msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
	// 					update.CallbackQuery.Message.MessageID, ans)
	// 				msg_t.ReplyMarkup = &keyboard

	// 				msg = msg_t

	// 			} else if callBackData[1] == Price {
	// 				// Случай, когда пришло 3 аргумента и выбрана цена через кнопки
	// 				// Переход к подтверждению
	// 				keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotifPrice), 2)
	// 				// Запись в кеш выбранной цены
	// 				n, _ := strconv.ParseFloat(callBackData[2], 32)
	// 				SetNotifCh.SetPrice(int(update.CallbackQuery.Message.Chat.ID), float32(n))

	// 				// Считывание из кеша всего объекта
	// 				setNotif := SetNotifCh.GetObject(int(update.CallbackQuery.Message.Chat.ID))

	// 				ans += fmt.Sprintf("Создается отсеживание:\nВалюта - %s\nКритерий - %s\nЦена - %.9f", setNotif.Crypto, setNotif.Criterion, setNotif.Price)
	// 				msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
	// 					update.CallbackQuery.Message.MessageID, ans)
	// 				msg_t.ReplyMarkup = &keyboard
	// 				msg = msg_t
	// 			}
	// 		}

	// 		// Для выбора критериев переход к выбору или вводу суммы
	// 		if update.CallbackQuery.Data == SetNotifCriterionMore || update.CallbackQuery.Data == SetNotifCriterionLess {
	// 			// ans := ChooseSum

	// 			// crypto := SetNotifCh.GetCrypto(int(update.CallbackQuery.Message.Chat.ID))
	// 			// cryptos, _ := coinmarketcup.GetLatestStruct(crypto)
	// 			// isFind := false
	// 			// prices := []PriceInfo{}

	// 			// if len(cryptos) == 1 {
	// 			// 	if cryptos[0].Find {
	// 			// 		isFind = true
	// 			// 		// Нужно вычислять количество знаков динамически
	// 			// 		prices = append(prices, PriceInfo{1, cryptos[0].Crypto.CryptoLastPrice * 1.01})
	// 			// 		prices = append(prices, PriceInfo{5, cryptos[0].Crypto.CryptoLastPrice * 1.05})
	// 			// 		prices = append(prices, PriceInfo{10, cryptos[0].Crypto.CryptoLastPrice * 1.1})
	// 			// 		prices = append(prices, PriceInfo{-1, cryptos[0].Crypto.CryptoLastPrice * 0.99})
	// 			// 		prices = append(prices, PriceInfo{-5, cryptos[0].Crypto.CryptoLastPrice * 0.95})
	// 			// 		prices = append(prices, PriceInfo{-10, cryptos[0].Crypto.CryptoLastPrice * 0.9})
	// 			// 	} else {
	// 			// 		ans += fmt.Sprintf("%s %s", cryptos[0].Crypto.CryptoName, "не найдена в базе")
	// 			// 	}
	// 			// }

	// 			// var row []tgbotapi.InlineKeyboardButton
	// 			// if isFind {
	// 			// 	for k, v := range prices {
	// 			// 		price := fmt.Sprintf("%f", v.Price)
	// 			// 		n := fmt.Sprintf(FormatFloatToString(v.Price)+" (%+d%%)", v.Price, v.Koeff)
	// 			// 		btn := tgbotapi.NewInlineKeyboardButtonData(n, SetNotifPrice+"_"+price)
	// 			// 		row = append(row, btn)
	// 			// 		// Делим на N строк по 3 элемента
	// 			// 		if (k+1)%3 == 0 {
	// 			// 			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	// 			// 			row = nil
	// 			// 		}
	// 			// 	}
	// 			// }
	// 			// row = append(row, tgbotapi.NewInlineKeyboardButtonData("Назад", Start))
	// 			// row = append(row, tgbotapi.NewInlineKeyboardButtonData("Ввод", SetNotifPriceEnter))
	// 			// keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	// 			// msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
	// 			// 	update.CallbackQuery.Message.MessageID, ans)
	// 			// msg_t.ReplyMarkup = &keyboard

	// 			// msg = msg_t

	// 		} else if update.CallbackQuery.Data == SetNotifPriceYes || update.CallbackQuery.Data == SetNotifPriceNo {
	// 			ans += "Здесь можно завести оповещения"
	// 			keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotif), 2)
	// 			msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
	// 				update.CallbackQuery.Message.MessageID, ans)
	// 			msg_t.ReplyMarkup = &keyboard
	// 			msg = msg_t
	// 		}
	// 	} else if len(callBackData) == 4 && callBackData[0] == SetNotif && callBackData[1] == Crypto && callBackData[2] == Yet {
	// 		// Пришло 4 аргумента для меню выбора КВ
	// 		ans = ChooseGetCrypto

	// 		offset, _ := strconv.Atoi(callBackData[3])
	// 		keyboard = GetCryptoListOffset(offset, SetNotifCrypto, SetNotif)

	// 		msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
	// 			update.CallbackQuery.Message.MessageID, ans)
	// 		msg_t.ReplyMarkup = &keyboard

	// 		msg = msg_t
	// 	}
	// } else if update.Message.ReplyToMessage != nil {

	// 	switch update.Message.ReplyToMessage.Text {
	// 	case EnterCrypto:
	// 		keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotifCriterion), 2)

	// 		ans = "Выбрана криптовалюта: " + update.Message.Text + "\nВыберите критерий"
	// 		// Запись в кеш введенной КВ
	// 		SetNotifCh.SetCrypto(int(update.Message.Chat.ID), update.Message.Text)

	// 		msg_t := tgbotapi.NewMessage(update.Message.Chat.ID,
	// 			ans)
	// 		msg_t.ReplyMarkup = &keyboard
	// 		msg = msg_t

	// 	case EnterSum:

	// 		n, err := strconv.ParseFloat(update.Message.Text, 32)
	// 		if err != nil {
	// 			ans += "Ошибка преобразования цены. Измените цену и повторите заного\n"

	// 			var row []tgbotapi.InlineKeyboardButton
	// 			row = append(row, tgbotapi.NewInlineKeyboardButtonData("Назад", SetNotifCrypto))
	// 			row = append(row, tgbotapi.NewInlineKeyboardButtonData("Повтор ввода", SetNotifPriceEnter))
	// 			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	// 		} else {
	// 			keyboard = MenuToInlineKeyboard(keyboardBot.GetMainMenuInlineMarkupFromNode(SetNotifPrice), 2)
	// 			// Запись в кеш введенной цены
	// 			SetNotifCh.SetPrice(int(update.Message.Chat.ID), float32(n))

	// 			// Считывание из кеша всего объекта
	// 			setNotif := SetNotifCh.GetObject(int(update.Message.Chat.ID))

	// 			ans += fmt.Sprintf("Создается отсеживание:\nВалюта - %s\nКритерий - %s\nЦена - %.9f", setNotif.Crypto, setNotif.Criterion, setNotif.Price)

	// 		}
	// 		msg_t := tgbotapi.NewMessage(update.Message.Chat.ID,
	// 			ans)
	// 		msg_t.ReplyMarkup = &keyboard
	// 		msg = msg_t
	// 	}
	// }
	return msg
}
