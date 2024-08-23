package main

import (
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/caching"
	"github.com/mbydanov/tg_golang_bot/internal/models"
	"github.com/mbydanov/tg_golang_bot/internal/notifications"
	retriever "github.com/mbydanov/tg_golang_bot/internal/retrieverCoins"
	"github.com/mbydanov/tg_golang_bot/internal/services"
	"github.com/mbydanov/tg_golang_bot/internal/tgbot"
)

func main() {

	time.Sleep(5 * time.Second)
	// Инициализация логирования
	services.InitLogger()
	// Создаем таблицу
	// if os.Getenv("CREATE_TABLE") == "yes" {
	// 	if os.Getenv("DB_SWITCH") == "on" {
	// 		if err := database.CreateTables(); err != nil {
	// 			services.Logging.Panic(err.Error())
	// 		}
	// 	}
	// }

	time.Sleep(2 * time.Second)
	// chConfig := make(chan config.ConfigStruct)
	// cfg := config.ConfigStruct{}
	chanModules := make(chan models.StatusChannel)
	// Получение настроек
	// go config.GetConfig(chConfig)

	// Кешировние
	// Кеширование пользователей
	caching.FillCache(caching.UsersCache, 100)
	// Кешируем словарь криптовалют
	caching.FillCache(caching.CryptoCache, 1000)
	// Кешируем отслеживания
	caching.FillCache(caching.TrackingCache, 10000)
	// Кешируем типы отслеживаний
	caching.FillCache(caching.TrackingTypeCache, 10)
	// Кешируем лимиты
	caching.FillCache(caching.LimitsCache, 100)
	// Кешируем словарь лимитов
	caching.FillCache(caching.LimitsDictCache, 10)
	// Кешируем коин маркеты
	caching.FillCache(caching.CoinMarketsCache, 5)
	caching.FillCache(caching.CoinMarketsEndpointCache, 10)
	caching.FillCache(caching.CoinMarketsHandCache, 100)

	// Функция считывания настроек из канала
	// go func() {
	// 	for {
	// 		// // Отправляем сообщение об ошибке
	// 		val, ok := <-chConfig
	// 		if ok {
	// 			if val.MsgError != nil {
	// 				ch <- models.StatusRetriever{MsgError: val.MsgError}
	// 			} else {
	// 				cfg = val
	// 			}
	// 		}
	// 	}
	// }()
	// for {
	// 	if cfg.TmrRespRvt == 0 {
	// 		continue
	// 	} else {
	// 		break
	// 	}
	// }
	// Вызов функции автоматического обновления КВ
	go retriever.RunRetrieverCoins(
		chanModules)
	go notifications.RunNotification(
		chanModules)
	// Вызываем бота
	tgbot.TelegramBot(
		chanModules)
}
