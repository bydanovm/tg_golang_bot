package tgbot

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/mbydanov/tg_golang_bot/internal/models"
)

const (
	// Comands
	Start                 string = "start"                    // Начало
	NumberOfUsers         string = "number_of_users"          // Получить количество активных пользователей
	GetCrypto             string = "GetCrypto"                // Получить актуальную информацию по криптовалюте
	GetCryptoEnter        string = GetCrypto + Enter          // Ручной ввод КВ
	GetCryptoNext         string = GetCrypto + "Next"         // Следующая страница
	GetCryptoBack         string = GetCrypto + "Back"         // Предыдушщая страница / выход
	GetCryptoCurr         string = GetCrypto + Curr           // Узнать курса валюты
	GetCryptoCurrBack     string = GetCryptoCurr + "Back"     // Выход из курса валюты
	GetCryptoCurrSetNot   string = GetCryptoCurr + SetNotif   // Установить отслеживание из курса валют
	GetNotif              string = "GetNotif"                 // Получить свои оповещения
	GetNotifId            string = GetNotif + "Id"            // Получить оповещение по ИД
	GetNotifIdOn          string = GetNotifId + "On"          // Включить оповещение по ИД
	GetNotifIdOff         string = GetNotifId + "Off"         // Отключить оповещение по ИД
	GetNotifUp            string = GetNotif + Up              // Свои оповещения - Назад
	GetNotifBack          string = GetNotif + Back            // Свои оповещения - Назад
	GetNotifYet           string = GetNotif + Yet             // Свои оповещения - Вперед
	GetNotifCrypto        string = GetNotif + Crypto          //
	GetNotifCryptoYet     string = GetNotifCrypto + Yet       //
	SetNotif              string = "SetNotif"                 // Установить уведомления по изменению цены криптовалюты
	SetNotifCurr          string = SetNotif + Curr            // Выбор или ввод КВ для отслежнивания
	SetNotifCryptoEnter   string = SetNotifCurr + "Enter"     // Ввод своей КВ
	SetNotifCriterion     string = SetNotif + "Criterion"     //
	SetNotifCriterionMore string = SetNotifCriterion + "More" //
	SetNotifCriterionLess string = SetNotifCriterion + "Less" //
	SetNotifPrice         string = SetNotif + Price           //
	SetNotifPriceEnter    string = SetNotifPrice + Enter      //
	SetNotifYes           string = SetNotif + "Yes"           //
	SetNotifNo            string = SetNotif + "No"            //

	Help   string = "help"
	Crypto string = "Crypto"
	Price  string = "Price"
	Enter  string = "Enter"
	Yet    string = "Yet"
	Back   string = "Back"
	Curr   string = "Curr"
	Up     string = "Up"
)

type FuncHandler func(*tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup)

type tgBotMenu struct {
	buttons  *models.TreeNode
	function map[string]FuncHandler
	Init     bool
}

var keyboardBot = initMenu()

type buttonInfo struct {
	text string
	data string
}

func initMenu() *tgBotMenu {
	buttons := models.InitTree()
	// keyboardBot.Add(GetCrypto, "Узнать курс", "0", true, getCrypto)
	// buttons.Add(GetCrypto, "Узнать курс", "0", true)
	// buttons.Add(GetCryptoCurr, "Узнать курс валюты", GetCrypto, false)
	// buttons.Add(GetCryptoYet, "Дальше", GetCrypto, true)
	// buttons.Add(SetNotif, "Оповещения", "0", true)
	// buttons.Add(GetNotif, "Текущие", SetNotif, true)
	// buttons.Add(GetNotifYet, "Дальше", GetNotif, false)
	// buttons.Add(GetNotifId, "Получить отслеживание по ID", GetNotif, false)
	// buttons.Add(GetNotifIdOn, "Отключить", GetNotifId, false)
	// buttons.Add(GetNotifIdOff, "Включить", GetNotifId, false)
	// buttons.Add(SetNotifCrypto, "Новое отслеживание", SetNotif, true)
	// buttons.Add(SetNotifPrice, "Установить цену", SetNotifCrypto, true)
	// buttons.Add(SetNotifPriceYes, "Да", SetNotifPrice, true)
	// buttons.Add(SetNotifPriceNo, "Нет", SetNotifPrice, true)
	// buttons.Add(Help, "Справка", "0", true)

	menu := &tgBotMenu{
		buttons:  buttons,
		function: make(map[string]FuncHandler),
	}

	return menu
}

func (tgm *tgBotMenu) Add(name, desc, parentId string, visible bool, foo ...FuncHandler) {
	tgm.buttons.Add(name, desc, parentId, visible)
	for _, v := range foo {
		tgm.function[name] = v
		break
	}
}

func (tgm *tgBotMenu) GetFunc(name string) FuncHandler {
	if _, ok := tgm.function[name]; !ok {
		return nil
	}
	return tgm.function[name]
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
	nodeParent := tgm.buttons.GetParentNode(node)
	if nodeParent != nil {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("Назад", nodeParent.Name))
	}
	nodesChild := tgm.buttons.GetNodeChild(node)
	for _, v := range nodesChild {
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

// Сделать кнопки для InlineKeyboard по имени узла
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
func MenuToInlineFromNode(node string, column int) tgbotapi.InlineKeyboardMarkup {
	buttons := keyboardBot.GetMainMenuInlineMarkupFromNode(node)
	keyboard := MenuToInlineKeyboard(buttons, column)
	return keyboard
}
