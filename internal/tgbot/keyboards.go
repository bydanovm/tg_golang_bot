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
	GetNotif              string = "GetNotif"            // Получить свои оповещения
	GetNotifUp            string = GetNotif + "_" + Up   // Свои оповещения - Назад
	GetNotifBack          string = GetNotif + "_" + Back // Свои оповещения - Назад
	GetNotifYet           string = GetNotif + "_" + Yet  // Свои оповещения - Вперед
	GetNotifCrypto        string = GetNotif + "_" + Crypto
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
	Back   string = "Back"
	Curr   string = "Curr"
	Up     string = "Up"
)

type tgBotMenu struct {
	buttons *models.TreeNode
}

var keyboardBot = initMenu()

type buttonInfo struct {
	text string
	data string
}

func initMenu() *tgBotMenu {
	buttons := models.InitTree()
	buttons.Add(GetCrypto, "Узнать курс", "0", true)
	// buttons.Add(GetCryptoCurr, "Узнать курс валюты", GetCrypto, false)
	// buttons.Add(GetCrypto, "Назад", GetCryptoCurr, true)
	// buttons.Add(SetNotifCrypto, "Установить отслеживание", GetCryptoCurr, true)
	buttons.Add(SetNotif, "Оповещения", "0", true)
	buttons.Add(Start, "Назад", SetNotif, true)
	buttons.Add(GetNotif, "Текущие", SetNotif, true)
	buttons.Add(SetNotifCrypto, "Новое", SetNotif, true)
	// buttons.Add(SetNotifCriterion, "Установить критерий", SetNotif, false)
	buttons.Add(SetNotifPrice, "Установить цену", SetNotif, false)
	buttons.Add(GetNotifBack, "Назад", GetNotif, true)
	buttons.Add(GetNotifUp, "Оповещения", GetNotif, true)
	buttons.Add(GetNotifYet, "Вперед", GetNotif, true)
	buttons.Add(SetNotifPriceYes, "Да", SetNotifPrice, true)
	buttons.Add(SetNotifPriceNo, "Нет", SetNotifPrice, true)
	// buttons.Add(SetNotifCriterionMore, "Больше >=", SetNotifCriterion, true)
	// buttons.Add(SetNotifCriterionLess, "Меньше <=", SetNotifCriterion, true)
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
func ConvertToButtonInlineKeyboard(in []buttonInfo, node string, column int) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var row []tgbotapi.InlineKeyboardButton
	for k, v := range in {
		btn := tgbotapi.NewInlineKeyboardButtonData(v.text, v.data)
		row = append(row, btn)
		// Делим на N строк по column элементов
		if (k+1)%column == 0 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			row = nil
		}
	}
	nodeButtons := keyboardBot.GetMainMenuInlineMarkupFromNode(node)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, nodeButtons)

	return keyboard
}
