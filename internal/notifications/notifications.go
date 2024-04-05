package notifications

import (
	"fmt"

	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/models"
	"github.com/mitchellh/mapstructure"
)

func RunNotification(
	retrieverNotifIn chan models.StatusChannel,
	notifTelegramOut chan models.StatusChannel) {
	// var chanRetrieverNotifIn models.StatusChannel
	for {
		// Считывание информации из канала от ретривера
		val, ok := <-retrieverNotifIn
		if ok {
			if val.Start {
				var charNotifTelegramOut models.StatusChannel
				// Функция работы с уведомлениями по КВ
				if res, err := notificationsCC(val.Data); err != nil {
					// Запись в канал об ошибке
					charNotifTelegramOut.Error = err
				} else {
					// Здесь будет запись в канал
					charNotifTelegramOut.Data = res
				}
				notifTelegramOut <- charNotifTelegramOut
			}
		}
	}
}

// Выходой интерфейс со структурой:
// 1 - ИД Пользователя
// 2 - Имя пользователя
// 2 - ИД чата
// 3 - ИД КВ
// 3 - Имя КВ
// 4 - Событие
func notificationsCC(bufferForNotif interface{}) (interface{}, error) {
	notifCCStruct := []NotificationsCCStruct{}

	// Определим для каждого отслеживания, было ли событие
	fields := database.TrackingCrypto{}
	expLst := []database.Expressions{}

	// Выбираем все записи из таблицы
	expLst = append(expLst, database.Expressions{
		Key: "IdTrkCrp", Operator: database.NotEQ, Value: `'0'`,
	})

	rs, find, _, err := database.ReadDataRow(&fields, expLst, 0)
	if err != nil {
		return nil, fmt.Errorf("notificationsCC:" + err.Error())
	}
	if !find {
		return nil, fmt.Errorf("notificationsCC:not find record")
	}
	for _, subRs := range rs {
		subFields := database.TrackingCrypto{}
		mapstructure.Decode(subRs, &subFields)

		dictCryptos := database.DictCrypto{}
		v, ok := bufferForNotif.(map[int]interface{})
		if !ok {
			return nil, fmt.Errorf("notificationsCC:error convert interface to map[int]interface")
		}
		mapstructure.Decode(v[subFields.DctCrpId], &dictCryptos)

		diff := dictCryptos.CryptoLastPrice - subFields.ValTrkCrp
		var event string
		// Проведем вычисления над КВ
		if diff < 0 && subFields.TypTrkCrpId == 1 { // Поднялась на N под пунктов (пп)
			event = "превышение"
		} else if diff > 0 && subFields.TypTrkCrpId == 2 { // Опустилась на N под пунктов (пп)
			event = "понижение"
		} else {
			continue
		}

		notifCCStruct = append(notifCCStruct, NotificationsCCStruct{
			subFields.UserId,
			"NameUsr",
			0,
			subFields.DctCrpId,
			dictCryptos.CryptoName,
			fmt.Sprintf("Произошло событие на криптовалютой %s: %s %v USD",
				dictCryptos.CryptoName, event, subFields.ValTrkCrp),
		})
	}

	return notifCCStruct, nil
}
