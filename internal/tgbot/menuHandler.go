package tgbot

import (
	"fmt"
	"os"
	"strconv"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/mbydanov/tg_golang_bot/internal/caching"
	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/services"
	"github.com/sirupsen/logrus"
)

const (
	cChooseCrypto    string = "Выберите или введите криптовалюту для отслеживания\n"
	cChooseSum       string = "Выберите или введите сумму для отслеживания\n"
	cEnterCrypto     string = "Введите криптовалюту для отслеживания\n"
	cEnterSum        string = "Введите сумму для отслеживания\n"
	cChooseGetCrypto string = "Выберите или введите криптовалюту"
	cEnterGetCrypto  string = "Введите криптовалюту"
	сGetNotifIdOn    string = "Отслеживание включено"
	сGetNotifIdOff   string = "Отслеживания отключено"
)

type PriceInfo struct {
	Koeff int
	Price float32
}

func menuHandler(update *tgbotapi.Update, bot tgbotapi.BotAPI) {
	var msg interface{}
	var err error
	var ans string
	var keyboard tgbotapi.InlineKeyboardMarkup
	var updateBot UpdateBot

	// Первичная инициализация меню хендлера
	if !keyboardBot.Init {
		keyboardBot.Add(Start, "Главная", "", KeyboardSettings{visible: true, multipage: false}, funcMenuStart)
		keyboardBot.Add(GetCrypto, "Узнать курс", Start, KeyboardSettings{visible: true, multipage: true}, funcGetCrypto)
		keyboardBot.Add(GetCryptoCurr, "Узнать курс валюты (инв)", GetCrypto, KeyboardSettings{visible: false, multipage: false}, funcGetCryptoCurr)

		keyboardBot.Add(GetCryptoCurrSetNotif, "Установить отслеживание", GetCryptoCurr, KeyboardSettings{visible: true, multipage: false}, funcSetNotifPrice)
		keyboardBot.Add(GetCryptoCurrSetNotifPriceEnter, "Цена введена - подтвердить (инв)", GetCryptoCurrSetNotif, KeyboardSettings{visible: false, multipage: false}, funcSetNotifPriceEnter)
		keyboardBot.Add(GetCryptoCurrSetNotifNo, "Отменить", GetCryptoCurrSetNotifPriceEnter, KeyboardSettings{visible: true, multipage: false}, funcSetNotifNo)
		keyboardBot.Add(GetCryptoCurrSetNotifNoOk, "Мои отслеживания", GetCryptoCurrSetNotifNo, KeyboardSettings{visible: true, multipage: false}, funcGetNotif)
		keyboardBot.Add(GetCryptoCurrSetNotifYes, "Подтвердить", GetCryptoCurrSetNotifPriceEnter, KeyboardSettings{visible: true, multipage: false}, funcSetNotifYes)
		keyboardBot.Add(GetCryptoCurrSetNotifYesOk, "Мои отслеживания", GetCryptoCurrSetNotifYes, KeyboardSettings{visible: true, multipage: false}, funcGetNotif)

		keyboardBot.Add(GetCryptoBack, "Назад", GetCrypto, KeyboardSettings{visible: true, multipage: false}, funcGetCrypto)
		keyboardBot.Add(GetCryptoNext, "Дальше", GetCrypto, KeyboardSettings{visible: true, multipage: false}, funcGetCrypto)
		keyboardBot.Add(GetNotif, "Отслеживания", Start, KeyboardSettings{visible: true, multipage: false}, funcGetNotif)
		keyboardBot.Add(GetNotifId, "Получить отслеживание по ID", GetNotif, KeyboardSettings{visible: false, multipage: false}, funcGetNotifId)
		keyboardBot.Add(GetNotifIdOn, "Включить", GetNotifId, KeyboardSettings{visible: false, multipage: false, homeBack: true}, funcGetNotifIdOn)
		keyboardBot.Add(GetNotifIdOnOk, "Мои отслеживания", GetNotifIdOn, KeyboardSettings{visible: true, multipage: false}, funcGetNotif)
		keyboardBot.Add(GetNotifIdOff, "Отключить", GetNotifId, KeyboardSettings{visible: false, multipage: false, homeBack: true}, funcGetNotifIdOff)
		keyboardBot.Add(GetNotifIdOffOk, "Мои отслеживания", GetNotifIdOff, KeyboardSettings{visible: true, multipage: false}, funcGetNotif)

		keyboardBot.Add(SetNotif, "Новое отслеживание", GetNotif, KeyboardSettings{visible: true, multipage: true}, funcSetNotif)
		keyboardBot.Add(SetNotifBack, "Назад", SetNotif, KeyboardSettings{visible: true, multipage: false}, funcSetNotif)
		keyboardBot.Add(SetNotifNext, "Дальше", SetNotif, KeyboardSettings{visible: true, multipage: false}, funcSetNotif)

		keyboardBot.Add(SetNotifPrice, "Установить цену", SetNotif, KeyboardSettings{visible: false, multipage: false}, funcSetNotifPrice)
		keyboardBot.Add(SetNotifPriceEnter, "Цена введена - подтвердить (инв)", SetNotifPrice, KeyboardSettings{visible: false, multipage: false}, funcSetNotifPriceEnter)
		keyboardBot.Add(SetNotifNo, "Отменить", SetNotifPriceEnter, KeyboardSettings{visible: true, multipage: false, homeBack: true}, funcSetNotifNo)
		keyboardBot.Add(SetNotifNoMyNotif, "Мои отслеживания", SetNotifNo, KeyboardSettings{visible: true, multipage: false}, funcGetNotif)
		keyboardBot.Add(SetNotifNoNewNotif, "Новое отслеживание", SetNotifNo, KeyboardSettings{visible: true, multipage: true}, funcSetNotif)
		keyboardBot.Add(SetNotifYes, "Подтвердить", SetNotifPriceEnter, KeyboardSettings{visible: true, multipage: false, homeBack: true}, funcSetNotifYes)
		keyboardBot.Add(SetNotifYesMyNotif, "Мои отслеживания", SetNotifYes, KeyboardSettings{visible: true, multipage: false}, funcGetNotif)

		keyboardBot.Add(Help, "Справка", Start, KeyboardSettings{visible: true, multipage: false})

		keyboardBot.Init = true
	}

	updateBot.FillInfo(update)

	pfunc := keyboardBot.GetFunc(updateBot.Data)
	if pfunc != nil {
		ans, keyboard, err = pfunc(&updateBot)
		if update.Message != nil {
			msg_t := tgbotapi.NewMessage(update.Message.Chat.ID,
				ans)
			msg_t.ReplyMarkup = &keyboard
			msg = msg_t
		} else {
			msg_t := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID, ans)
			msg_t.ReplyMarkup = &keyboard
			msg = msg_t
		}
	}

	switch msgConv := msg.(type) {
	case tgbotapi.EditMessageTextConfig:
		bot.Send(msgConv)
	case tgbotapi.MessageConfig:
		bot.Send(msgConv)
	}

	// Логирование
	// Получить команду
	var command, callBackData, text string
	if update.Message != nil {
		command = update.Message.Command()
		text = update.Message.Text
	}
	// Получить каллбек
	if update.CallbackQuery != nil {
		callBackData = update.CallbackQuery.Data
	}
	services.Logging.WithFields(logrus.Fields{
		"module": "menuHandler",
		"user": logrus.Fields{
			"userInfo": updateBot.User,
			"userCtrl": logrus.Fields{
				"menuCache":    updateBot.Menu,
				"command":      command,
				"callBackData": callBackData,
				"text":         text,
			},
		},
		"error": func(err error) interface{} {
			// if err
			return err
		}(err),
	}).Info()
}

func funcGetCrypto(updateBot *UpdateBot) (ans string, keyboard tgbotapi.InlineKeyboardMarkup, err error) {
	var isMode int
	// Вычислим offset по значению из кеша
	offset := 10

	if updateBot.Menu.OffsetNavi > 10 {
		offset = updateBot.Menu.OffsetNavi
	}
	callBackData := updateBot.Data

	// Обработка кнопок Назад/Вперед
	if callBackData[0] == GetCryptoNext { // Сработала кнопка дальше, запись значения в кеш
		callBackData[0] = keyboardBot.buttons.GetParentNode(callBackData[0]).Name
		offset += 10
	} else if callBackData[0] == GetCryptoBack { // Сработала кнопка назад
		callBackData[0] = keyboardBot.buttons.GetParentNode(callBackData[0]).Name
		if offset > 10 {
			offset -= 10
		} else {
			callBackData[0] = keyboardBot.buttons.GetParentNode(callBackData[0]).Name
			pfunc := keyboardBot.GetFunc(callBackData)
			if pfunc != nil {
				ans, keyboard, err = pfunc(updateBot)
				return ans, keyboard, err
			}
		}
	}
	caching.SetCache(MenuCache, updateBot.User.IdUsr, MenuInfo{
		CurrentMenu: GetCrypto,
		OffsetNavi:  offset,
	}, 0)

	listCryptoCur, lastList, _ := caching.GetCacheOffset(caching.CryptoCache, offset)
	if lastList {
		isMode = LastList
	}
	listButtons := make([]buttonInfo, 0, 10)
	for _, v := range listCryptoCur {
		listButtons = append(listButtons, buttonInfo{v.CryptoName, GetCryptoCurr + `_` + strconv.Itoa(v.CryptoId)})
	}
	ans = cChooseGetCrypto
	keyboard = ConvertToButtonInlineKeyboard(listButtons, callBackData[0], 3, isMode)
	return ans, keyboard, err
}

func funcGetCryptoCurr(updateBot *UpdateBot) (ans string, keyboard tgbotapi.InlineKeyboardMarkup, err error) {
	// Случай, когда валюта пришла из другого места
	if updateBot.Menu.IdCrypto > 0 {
		updateBot.Data = append(updateBot.Data, strconv.Itoa(updateBot.Menu.IdCrypto))
	}
	id, err := strconv.Atoi(updateBot.Data[1])
	if err != nil {
		// Возможно передана мнемоника, делаем поиск ее ключа и пишем в кешменю
		if v, ok := caching.GetCacheElementKeyChain(caching.CryptoCache, updateBot.Data[1]).(int); ok {
			id = v
		} else {
			return ans, keyboard, err
		}
	}

	crypto, err := caching.GetCacheByIdxInMap(caching.CryptoCache, id, 0)
	if err != nil {
		err = fmt.Errorf("%s %s %s", crypto.CryptoName, "не найдена в базе", err.Error())
		return ans, keyboard, err
	}

	updateBot.Menu.IdCrypto = crypto.CryptoId
	caching.SetCache(MenuCache, updateBot.User.IdUsr, updateBot.Menu, 0)

	ans = fmt.Sprintf("1 %s = "+FormatFloatToString(crypto.CryptoLastPrice)+" %s",
		crypto.CryptoName,
		crypto.CryptoLastPrice,
		"$")
	keyboard = MenuToInlineFromNode(updateBot.Data[0], 2)
	return ans, keyboard, err
}

func funcGetNotif(updateBot *UpdateBot) (ans string, keyboard tgbotapi.InlineKeyboardMarkup, err error) {
	ans = "Текущие отслеживания"
	// Вывести от новых к старым в формате Валюта - Значение
	trackings, _ := caching.GetCacheRecordsKeyChain(caching.TrackingCache, updateBot.User.IdUsr, true)
	// Создание списка кнопок
	listButtons := make([]buttonInfo, 0, 10)
	for _, v := range trackings {
		infoCurrency, er := caching.GetCacheByIdxInMap(caching.CryptoCache, v.DctCrpId, 0)
		if er != nil {
			err = fmt.Errorf("funcGetNotif: %s", er.Error())
			return ans, keyboard, err
		}
		listButtons = append(listButtons, buttonInfo{infoCurrency.CryptoName + " - " + fmt.Sprintf(FormatFloatToString(v.ValTrkCrp), v.ValTrkCrp) + " $", GetNotifId + "_" + fmt.Sprintf("%d", v.IdTrkCrp)})
	}
	keyboard = ConvertToButtonInlineKeyboard(listButtons, GetNotif, 3)
	return ans, keyboard, err
}

func funcGetNotifId(updateBot *UpdateBot) (ans string, keyboard tgbotapi.InlineKeyboardMarkup, err error) {
	id, _ := strconv.Atoi(updateBot.Data[1])
	infoTracking, _ := caching.GetCacheByIdxInMap(caching.TrackingCache, id, 0)
	infoCurrency, _ := caching.GetCacheByIdxInMap(caching.CryptoCache, infoTracking.DctCrpId, 0)
	// Считывание типа триггера - добавить
	// Запись в кеш инфы для операций
	caching.SetCache(MenuCache, updateBot.User.IdUsr, MenuInfo{
		IdTracking: infoTracking.IdTrkCrp,
	}, 0)
	ans = fmt.Sprintf("Выбрано отслеживание по %s\nТекущая стоимость 1 %s = "+FormatFloatToString(infoTracking.ValTrkCrp)+" $\nТип триггера: %d\nСрабатывание триггера: "+FormatFloatToString(infoTracking.ValTrkCrp)+"\n", infoCurrency.CryptoName, infoCurrency.CryptoName, infoCurrency.CryptoLastPrice, infoTracking.TypTrkCrpId, infoTracking.ValTrkCrp)

	keyboard = MenuToInlineFromNode(GetNotifId, 2)

	return ans, keyboard, err
}

func funcGetNotifIdOn(updateBot *UpdateBot) (ans string, keyboard tgbotapi.InlineKeyboardMarkup, err error) {
	ans = сGetNotifIdOn
	keyboard = MenuToInlineFromNode(GetNotifIdOn, 2)

	return ans, keyboard, err
}

func funcGetNotifIdOff(updateBot *UpdateBot) (ans string, keyboard tgbotapi.InlineKeyboardMarkup, err error) {
	ans = сGetNotifIdOff
	keyboard = MenuToInlineFromNode(GetNotifIdOff, 2)

	return ans, keyboard, err
}

func funcSetNotif(updateBot *UpdateBot) (ans string, keyboard tgbotapi.InlineKeyboardMarkup, err error) {
	offset := 10
	ans = cChooseGetCrypto

	if updateBot.Menu.OffsetNavi > 10 {
		offset = updateBot.Menu.OffsetNavi
	}

	// Обработка кнопок Назад/Вперед
	if updateBot.Data[0] == SetNotifNext { // Сработала кнопка дальше, запись значения в кеш
		updateBot.Data[0] = keyboardBot.buttons.GetParentNode(updateBot.Data[0]).Name
		offset += 10
	} else if updateBot.Data[0] == SetNotifBack { // Сработала кнопка назад
		updateBot.Data[0] = keyboardBot.buttons.GetParentNode(updateBot.Data[0]).Name
		if offset > 10 {
			offset -= 10
		} else {
			updateBot.Data[0] = keyboardBot.buttons.GetParentNode(updateBot.Data[0]).Name
			pfunc := keyboardBot.GetFunc(updateBot.Data)
			if pfunc != nil {
				ans, keyboard, err = pfunc(updateBot)
				return ans, keyboard, err
			}
		}
	}
	caching.SetCache(MenuCache, updateBot.User.IdUsr, MenuInfo{OffsetNavi: offset}, 0)

	listCryptoCur, _, _ := caching.GetCacheOffset(caching.CryptoCache, offset)
	listButtons := make([]buttonInfo, 0, 10)
	for _, v := range listCryptoCur {
		listButtons = append(listButtons, buttonInfo{v.CryptoName, SetNotifPrice + `_` + strconv.Itoa(v.CryptoId)})
	}
	keyboard = ConvertToButtonInlineKeyboard(listButtons, updateBot.Data[0], 3)
	return ans, keyboard, err
}

func funcSetNotifPrice(updateBot *UpdateBot) (ans string, keyboard tgbotapi.InlineKeyboardMarkup, err error) {
	// Возможно переключить на пакет caching
	// Пишем в кеш ИД крипты
	// Случай, когда валюта пришла из другого места
	if updateBot.Menu.IdCrypto > 0 {
		updateBot.Data = append(updateBot.Data, strconv.Itoa(updateBot.Menu.IdCrypto))
	}
	n, err := strconv.Atoi(updateBot.Data[1])
	if err != nil {
		return ans, keyboard, err
	}
	SetNotifCh.SetIdCrypto(int(updateBot.User.ChatIdUsr), n)

	ans = cChooseSum

	infoCurrency, err := caching.GetCacheByIdxInMap(caching.CryptoCache, n, 0)
	if err != nil {
		// Здесь должна быть обработка ошибки
		ans += fmt.Sprintf("%s %s", "Выбранная криптовалюта", "не найдена в базе")
		return ans, keyboard, err
	}

	// Запись в кеш текущей цены
	SetNotifCh.SetCurrentPrice(int(updateBot.User.ChatIdUsr), infoCurrency.CryptoLastPrice)
	// Запись в кеш мнемоники криптовалюты
	SetNotifCh.SetCrypto(int(updateBot.User.ChatIdUsr), infoCurrency.CryptoName)

	ans += "Выбрана криптовалюта: " + infoCurrency.CryptoName + "\n"

	prices := []PriceInfo{}

	prices = append(prices, PriceInfo{1, infoCurrency.CryptoLastPrice * 1.01})
	prices = append(prices, PriceInfo{5, infoCurrency.CryptoLastPrice * 1.05})
	prices = append(prices, PriceInfo{10, infoCurrency.CryptoLastPrice * 1.1})
	prices = append(prices, PriceInfo{-1, infoCurrency.CryptoLastPrice * 0.99})
	prices = append(prices, PriceInfo{-5, infoCurrency.CryptoLastPrice * 0.95})
	prices = append(prices, PriceInfo{-10, infoCurrency.CryptoLastPrice * 0.9})

	listButtons := make([]buttonInfo, 0, 10)
	for _, v := range prices {
		price := fmt.Sprintf("%f", v.Price)
		n := fmt.Sprintf(FormatFloatToString(v.Price)+" (%+d%%)", v.Price, v.Koeff)
		listButtons = append(listButtons, buttonInfo{n, SetNotifPriceEnter + `_` + price})
	}
	keyboard = ConvertToButtonInlineKeyboard(listButtons, updateBot.Data[0], 3)

	return ans, keyboard, err
}

func funcSetNotifPriceEnter(updateBot *UpdateBot) (ans string, keyboard tgbotapi.InlineKeyboardMarkup, err error) {
	// Запись в кеш выбранной цены
	n, err := strconv.ParseFloat(updateBot.Data[1], 32)
	if err != nil {
		return ans, keyboard, err
	}
	SetNotifCh.SetPrice(int(updateBot.User.ChatIdUsr), float32(n))

	// Определяем и записываем критерий (тип триггера)
	idType := 0
	if SetNotifCh.GetCurrentPrice(int(updateBot.User.ChatIdUsr)) <= float32(n) {
		idType = caching.GetCacheElementKeyChain(caching.TrackingTypeCache, "RAISE_V").(int)
	} else {
		idType = caching.GetCacheElementKeyChain(caching.TrackingTypeCache, "FALL_V").(int)
	}
	SetNotifCh.SetCriterion(int(updateBot.User.ChatIdUsr), idType)

	// Считывание из кеша всего объекта
	setNotif := SetNotifCh.GetObject(int(updateBot.User.ChatIdUsr))

	ans = fmt.Sprintf("Создается отсеживание:\nВалюта - %s\nЦена - %.9f", setNotif.Crypto, setNotif.Price)

	keyboard = MenuToInlineFromNode(SetNotifPriceEnter, 2)

	return ans, keyboard, err
}

func funcSetNotifYes(updateBot *UpdateBot) (ans string, keyboard tgbotapi.InlineKeyboardMarkup, err error) {
	// Считывание из кеша всего объекта
	setNotif := SetNotifCh.GetObject(int(updateBot.User.ChatIdUsr))

	// Найти Тип отслеживания
	idType := setNotif.IdCriterion
	// Установка лимита
	limit := database.Limits{}
	if ans == "" {
		// Чтение LMT003
		lmt, _ := caching.GetCacheByIdxInMap(caching.LimitsDictCache, caching.GetCacheElementKeyChain(caching.LimitsDictCache, "LMT003").(int))
		limit = database.Limits{
			ValAvailLmt: lmt.StdValLmt,
			ActiveLmt:   true,
			UserId:      updateBot.User.IdUsr,
			LtmDctId:    lmt.IdLmtDct,
		}
		if _, id, err := caching.WriteCache(caching.LimitsCache, 0, limit); err != nil {
			ans += fmt.Sprintf("tgbot:%s\n", err.Error())
		} else {
			limit.IdLmt = int(id)
		}
	}

	tracking := database.TrackingCrypto{
		DctCrpId:    setNotif.IdCrypto,
		TypTrkCrpId: idType,
		LmtId:       limit.IdLmt,
		UserId:      updateBot.User.IdUsr,
		ValTrkCrp:   SetNotifCh.GetPrice(int(updateBot.User.ChatIdUsr)),
		OnTrkCrp:    true,
	}
	if _, _, err := caching.WriteCache(caching.TrackingCache, 0, tracking); err != nil {
		ans += fmt.Sprintf("tgbot:%s\n", err.Error())
	} else {
		ans += fmt.Sprintf("Отслеживание по криптовалюте %s успешно добавлено\n", SetNotifCh.GetCrypto(int(updateBot.User.ChatIdUsr)))
	}
	// Запись в кеш
	keyboard = MenuToInlineFromNode(updateBot.Data[0], 2)

	return ans, keyboard, err
}

func funcSetNotifNo(updateBot *UpdateBot) (ans string, keyboard tgbotapi.InlineKeyboardMarkup, err error) {
	ans = "Отслеживание не сохранено\n"

	//Надо чистить елемент мапы
	keyboard = MenuToInlineFromNode(updateBot.Data[0], 2)
	return ans, keyboard, err
}

func funcMenuStart(updateBot *UpdateBot) (ans string, keyboard tgbotapi.InlineKeyboardMarkup, err error) {
	ans = "Привет! Я - " + os.Getenv("BOT_NAME") + " помогу тебе знать актуальную информацию по криптовалюте\n" +
		"Используй кнопки ниже, чтобы узнать интересующую информацию.\n"
	keyboard = MenuToInlineFromNode(Start, 2)
	return ans, keyboard, err
}
