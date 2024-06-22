package tgbot

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/mbydanov/tg_golang_bot/internal/models"
)

const (
	// Comands
	Start         string = "start"           // Начало
	NumberOfUsers string = "number_of_users" // Получить количество активных пользователей
	GetCrypto     string = "getcrypto"       // Получить актуальную информацию по криптовалюте
	SetNotif      string = "setnotif"        // Установить уведомления по изменению цены криптовалюты
	Help          string = "help"
)

type tgBotMenu struct {
	buttons *models.TreeNode
}

// var tgBotMenu = models.InitTree()

func initMenu() *tgBotMenu {
	buttons := models.InitTree()
	buttons.Add(GetCrypto, "Узнать курс", "0")
	buttons.Add(GetCrypto+"_BTC", "Узнать курс BTC", GetCrypto)
	buttons.Add(GetCrypto+"_ETH", "Узнать курс ETH", GetCrypto)
	buttons.Add(GetCrypto+"_TON", "Узнать курс TON", GetCrypto)
	buttons.Add(SetNotif, "Оповещения", "0")
	buttons.Add(Help, "Справка", "0")

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
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(v.Description, v.Name))
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
