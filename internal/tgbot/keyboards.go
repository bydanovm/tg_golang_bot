package tgbot

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/mbydanov/tg_golang_bot/internal/models"
)

const (
	// Comands
	Start                 string = "start"           // Начало
	NumberOfUsers         string = "number_of_users" // Получить количество активных пользователей
	GetCrypto             string = "GetCrypto"       // Получить актуальную информацию по криптовалюте
	GetCryptoEnter        string = GetCrypto + "_" + Enter
	GetCryptoYet          string = GetCrypto + "_" + Yet
	GetCryptoCurr         string = GetCrypto + "_" + Curr
	SetNotif              string = "SetNotif"                // Установить уведомления по изменению цены криптовалюты
	SetNotifCrypto        string = SetNotif + "_" + Crypto   // Выбор или ввод КВ
	SetNotifCryptoEnter   string = SetNotifCrypto + "_Enter" // Ввод своей КВ
	SetNotifCriterion     string = SetNotif + "_Criterion"
	SetNotifCriterionMore string = SetNotifCriterion + "_More"
	SetNotifCriterionLess string = SetNotifCriterion + "_Less"
	SetNotifPrice         string = SetNotif + "_" + Price
	SetNotifPriceEnter    string = SetNotifPrice + "_" + Enter
	SetNotifPriceYes      string = SetNotifPrice + "_Yes"
	SetNotifPriceNo       string = SetNotifPrice + "_No"

	Help   string = "help"
	Crypto string = "Crypto"
	Price  string = "Price"
	Enter  string = "Enter"
	Yet    string = "Yet"
	Curr   string = "Curr"
)

type tgBotMenu struct {
	buttons *models.TreeNode
}

func initMenu() *tgBotMenu {
	buttons := models.InitTree()
	buttons.Add(GetCrypto, "Узнать курс", "0", true)
	// buttons.Add(GetCryptoCurr, "Узнать курс валюты", GetCrypto, false)
	// buttons.Add(GetCrypto, "Назад", GetCryptoCurr, true)
	// buttons.Add(SetNotifCrypto, "Установить отслеживание", GetCryptoCurr, true)
	buttons.Add(SetNotif, "Оповещения", "0", true)
	buttons.Add(SetNotifCrypto, "Выбрать крипту", SetNotif, true)
	buttons.Add(SetNotifCriterion, "Установить критерий", SetNotif, false)
	buttons.Add(SetNotifPrice, "Установить цену", SetNotif, false)
	buttons.Add(SetNotifPriceYes, "Да", SetNotifPrice, true)
	buttons.Add(SetNotifPriceNo, "Нет", SetNotifPrice, true)
	buttons.Add(SetNotifCriterionMore, "Больше >=", SetNotifCriterion, true)
	buttons.Add(SetNotifCriterionLess, "Меньше <=", SetNotifCriterion, true)
	buttons.Add(Start, "Назад", SetNotif, true)
	buttons.Add(Help, "Справка", "0", true)

	menu := &tgBotMenu{
		buttons: buttons,
	}
	return menu
}
func (tgm *tgBotMenu) GetMainMenuReplyMarkup() (buttons []tgbotapi.KeyboardButton) {
	nodes := tgm.buttons.GetNodeChild("0")
	for _, v := range nodes {
		buttons = append(buttons, tgbotapi.KeyboardButton{Text: v.Description})
	}
	return buttons
}

// Получить меню для формата InlineKeyboardMarkup
func (tgm *tgBotMenu) GetMainMenuInlineMarkup() (buttons []tgbotapi.InlineKeyboardButton) {
	nodes := tgm.buttons.GetNodeChild("0")
	for _, v := range nodes {
		if v.Visible {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(v.Description, v.Name))
		}
	}
	return buttons
}

func (tgm *tgBotMenu) GetMainMenuInlineMarkupFromNode(node string) (buttons []tgbotapi.InlineKeyboardButton) {
	nodes := tgm.buttons.GetNodeChild(node)
	for _, v := range nodes {
		if v.Visible {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(v.Description, v.Name))
		}
	}
	return buttons
}

// Получить готовую клавиатуру
func MenuToInlineKeyboard(buttons []tgbotapi.InlineKeyboardButton, columns int) (keyboard tgbotapi.InlineKeyboardMarkup) {
	row := []tgbotapi.InlineKeyboardButton{}
	for k, v := range buttons {
		row = append(row, v)
		if (k+1)%columns == 0 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			row = nil
		} else if (k+1) == len(buttons) && (k+1)%2 == 1 { // Если последний элемент
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}
	return keyboard
}

// Сделать кнопки для InlineKeyboard
func ConvertToButtonInlineKeyboard(in interface{}) (buttons []tgbotapi.InlineKeyboardButton) {

	return buttons
}
