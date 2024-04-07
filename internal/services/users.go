package services

import (
	"fmt"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/mbydanov/tg_golang_bot/internal/database"
)

// Проверка, что пользователь есть в базе
func CheckUser(messageIn *tgbotapi.Message) error {
	user := database.Users{
		IdUsr:     messageIn.From.ID,
		TsUsr:     time.Now(),
		NameUsr:   messageIn.From.UserName,
		FirstName: messageIn.From.FirstName,
		LastName:  messageIn.From.LastName,
		LangCode:  messageIn.From.LanguageCode,
		IsBot:     messageIn.From.IsBot,
		IsBanned:  false,
		ChatIdUsr: messageIn.Chat.ID,
		IdLvlSec:  5}
	idUser, err := user.Find()
	if err != nil {
		return fmt.Errorf("CheckUser:" + err.Error())
	}
	if idUser < 0 {
		// return fmt.Errorf("CheckUser:User not found:%v", idUser)
		// Если пользователя нет в базе - добавляем
		_, err = user.Add()
		if err != nil {
			return fmt.Errorf("CheckUser:" + err.Error())
		}
		return nil
	}
	return nil
}
