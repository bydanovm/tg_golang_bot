package main

import (
	"os"
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/config"
	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/models"
	"github.com/mbydanov/tg_golang_bot/internal/notifications"
	retrievercoins "github.com/mbydanov/tg_golang_bot/internal/retrieverCoins"
	"github.com/mbydanov/tg_golang_bot/internal/tgbot"
)

func main() {

	time.Sleep(5 * time.Second)

	// Создаем таблицу
	if os.Getenv("CREATE_TABLE") == "yes" {
		if os.Getenv("DB_SWITCH") == "on" {
			if err := database.CreateTables(); err != nil {
				panic(err)
			}
		}
	}

	time.Sleep(2 * time.Second)
	ch := make(chan models.StatusRetriever)
	chConfig := make(chan config.ConfigStruct)
	cfg := config.ConfigStruct{}
	chanRetrieverNotif := make(chan models.StatusChannel)
	chanNotifTelegram := make(chan models.StatusChannel)
	// Получение настроек
	go config.GetConfig(chConfig)

	// Функция считывания настроек из канала
	go func() {
		for {
			// // Отправляем сообщение об ошибке
			val, ok := <-chConfig
			if ok {
				if val.MsgError != nil {
					ch <- models.StatusRetriever{MsgError: val.MsgError}
				} else {
					cfg = val
				}
			}
		}
	}()
	for {
		if cfg.TmrRespRvt == 0 {
			continue
		} else {
			break
		}
	}
	// Вызов функции автоматического обновления КВ
	go retrievercoins.RunRetrieverCoins(
		cfg.TmrRespRvt, ch, chanRetrieverNotif)
	go notifications.RunNotification(
		chanRetrieverNotif, chanNotifTelegram)
	// Вызываем бота
	tgbot.TelegramBot(
		ch, chanNotifTelegram)
}
