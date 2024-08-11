package tgbot

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/mbydanov/tg_golang_bot/internal/models"
)

const (
	// Comands
	Start                           string = "start"                  // Начало
	NumberOfUsers                   string = "number_of_users"        // Получить количество активных пользователей
	GetCrypto                       string = "GetCrypto"              // Получить актуальную информацию по криптовалюте
	GetCryptoEnter                  string = GetCrypto + Enter        // Ручной ввод КВ
	GetCryptoNext                   string = GetCrypto + "Next"       // Следующая страница
	GetCryptoBack                   string = GetCrypto + "Back"       // Предыдушщая страница / выход
	GetCryptoCurr                   string = GetCrypto + Curr         // Узнать курса валюты
	GetCryptoCurrBack               string = GetCryptoCurr + "Back"   // Выход из курса валюты
	GetCryptoCurrSetNotif           string = GetCryptoCurr + SetNotif // Установить отслеживание из курса валют
	GetCryptoCurrSetNotifPriceEnter string = GetCryptoCurr + SetNotifPriceEnter
	GetCryptoCurrSetNotifNo         string = GetCryptoCurr + SetNotifNo
	GetCryptoCurrSetNotifNoOk       string = GetCryptoCurr + SetNotifNoMyNotif
	GetCryptoCurrSetNotifYes        string = GetCryptoCurr + SetNotifYes
	GetCryptoCurrSetNotifYesOk      string = GetCryptoCurr + SetNotifNoNewNotif
	GetNotif                        string = "GetNotif"                 // Получить свои оповещения
	GetNotifId                      string = GetNotif + "Id"            // Получить оповещение по ИД
	GetNotifIdOn                    string = GetNotifId + "On"          // Включить оповещение по ИД
	GetNotifIdOnOk                  string = GetNotifIdOn + "Ok"        //
	GetNotifIdOff                   string = GetNotifId + "Off"         // Отключить оповещение по ИД
	GetNotifIdOffOk                 string = GetNotifIdOff + "Ok"       //
	GetNotifUp                      string = GetNotif + Up              // Свои оповещения - Назад
	GetNotifBack                    string = GetNotif + Back            // Свои оповещения - Назад
	GetNotifYet                     string = GetNotif + Yet             // Свои оповещения - Вперед
	GetNotifCrypto                  string = GetNotif + Crypto          //
	GetNotifCryptoYet               string = GetNotifCrypto + Yet       //
	SetNotif                        string = "SetNotif"                 // Установить уведомления по изменению цены криптовалюты
	SetNotifBack                    string = SetNotif + Back            //
	SetNotifNext                    string = SetNotif + Next            //
	SetNotifCurr                    string = SetNotif + Curr            // Выбор или ввод КВ для отслежнивания
	SetNotifCryptoEnter             string = SetNotifCurr + "Enter"     // Ввод своей КВ
	SetNotifCriterion               string = SetNotif + "Criterion"     //
	SetNotifCriterionMore           string = SetNotifCriterion + "More" //
	SetNotifCriterionLess           string = SetNotifCriterion + "Less" //
	SetNotifPrice                   string = SetNotif + Price           //
	SetNotifPriceEnter              string = SetNotifPrice + Enter      //
	SetNotifYes                     string = SetNotif + "Yes"           //
	SetNotifYesMyNotif              string = SetNotifYes + "MyNotif"    //
	SetNotifNo                      string = SetNotif + "No"            //
	SetNotifNoMyNotif               string = SetNotifNo + "MyNotif"     //
	SetNotifNoNewNotif              string = SetNotifNo + "NewNotif"    //

	Help   string = "help"
	Crypto string = "Crypto"
	Price  string = "Price"
	Enter  string = "Enter"
	Yet    string = "Yet"
	Back   string = "Back"
	Next   string = "Next"
	Curr   string = "Curr"
	Up     string = "Up"

	FirstList int = 999
	LastList  int = 1000
)

type FuncHandler func(*tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup)
type keyboardFeature struct {
	KeyboardSettings
	function FuncHandler
}
type KeyboardSettings struct {
	visible   bool
	multipage bool
}
type tgBotMenu struct {
	buttons *models.TreeNode
	feature map[string]keyboardFeature
	Init    bool
}

var keyboardBot = initMenu()

type buttonInfo struct {
	text string
	data string
}

func initMenu() *tgBotMenu {
	buttons := models.InitTree()

	menu := &tgBotMenu{
		buttons: buttons,
		feature: make(map[string]keyboardFeature),
	}

	return menu
}

func (tgm *tgBotMenu) Add(name, desc, parentId string, settings KeyboardSettings, foo ...FuncHandler) {
	tgm.buttons.Add(name, desc, parentId)
	tgm.feature[name] = keyboardFeature{
		KeyboardSettings: KeyboardSettings{
			settings.visible,
			settings.multipage,
		}}

	for _, v := range foo {
		item, ok := tgm.feature[name]
		if ok {
			item.function = v
			tgm.feature[name] = item
		}
		break
	}
}

func (tgm *tgBotMenu) GetFunc(name string) FuncHandler {
	if _, ok := tgm.feature[name]; !ok {
		return nil
	}
	return tgm.feature[name].function
}

func (tgm *tgBotMenu) GetMultiPage(name string) bool {
	if _, ok := tgm.feature[name]; !ok {
		return false
	}
	return tgm.feature[name].multipage
}

// func (tgm *tgBotMenu) EditVisibleButton(name string) {
// 	if item, ok := tgm.feature[name]; ok {
// 		item.visible = !item.visible
// 		tgm.feature[name] = item
// 	}
// }

func (tgm *tgBotMenu) GetMainMenuReplyMarkup() (buttons []tgbotapi.KeyboardButton) {
	nodes := tgm.buttons.GetNodeChild("0")
	for _, v := range nodes {
		buttons = append(buttons, tgbotapi.KeyboardButton{Text: v.Description})
	}
	return buttons
}

// Получить меню для формата InlineKeyboardMarkup
// func (tgm *tgBotMenu) GetMainMenuInlineMarkup() (buttons []tgbotapi.InlineKeyboardButton) {
// 	nodes := tgm.buttons.GetNodeChild("0")
// 	for _, v := range nodes {
// 		if item, ok := tgm.feature[name]; ok {
// 			if item.Visible {
// 				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(v.Description, v.Name))
// 			}
// 		}
// 	}
// 	return buttons
// }

func (tgm *tgBotMenu) GetMainMenuInlineMarkupFromNode(node string, mode ...int) (buttons []tgbotapi.InlineKeyboardButton) {
	var isMode int
	for _, v := range mode {
		isMode = v
		break
	}
	isMultiPage := tgm.GetMultiPage(node)
	if !isMultiPage {
		nodeParent := tgm.buttons.GetParentNode(node)
		if nodeParent != nil {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("Назад", nodeParent.Name))
		}
	}
	nodesChild := tgm.buttons.GetNodeChild(node)
	for _, v := range nodesChild {
		if item, ok := tgm.feature[v.Name]; ok {
			if item.visible {
				if !(isMode == LastList && v.Description == "Дальше") {
					buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(v.Description, v.Name))
				}
			}
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
func ConvertToButtonInlineKeyboard(in []buttonInfo, node string, column int, mode ...int) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var row []tgbotapi.InlineKeyboardButton
	var modeIs int
	for k, v := range in {
		btn := tgbotapi.NewInlineKeyboardButtonData(v.text, v.data)
		row = append(row, btn)
		// Делим на N строк по column элементов
		if (k+1)%column == 0 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			row = nil
		}
	}
	for _, v := range mode {
		modeIs = v
		break
	}
	nodeButtons := keyboardBot.GetMainMenuInlineMarkupFromNode(node, modeIs)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, nodeButtons)

	return keyboard
}
func MenuToInlineFromNode(node string, column int) tgbotapi.InlineKeyboardMarkup {
	buttons := keyboardBot.GetMainMenuInlineMarkupFromNode(node)
	keyboard := MenuToInlineKeyboard(buttons, column)
	return keyboard
}
