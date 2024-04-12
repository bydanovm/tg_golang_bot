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
// 3 - ИД чата
// 4 - ИД КВ
// 5 - Имя КВ
// 6 - Событие
func notificationsCC(bufferForNotif interface{}) (interface{}, error) {
	notifCCStruct := []NotificationsCCStruct{}

	// Определим для каждого отслеживания, было ли событие
	// Выбираем все записи из таблицы с включенным отслеживанием
	expLst := []database.Expressions{
		{Key: "IdTrkCrp", Operator: database.NotEQ, Value: `'0'`},
		{Key: "OnTrkCrp", Operator: database.EQ, Value: `true`},
	}
	rs, find, _, err := database.ReadDataRow(&database.TrackingCrypto{}, expLst, 0)
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
		// Получаем лимит и увеличиваем его
		lmt := database.Limits{}
		if err := lmt.GetLimit("LMT003", subFields.UserId); err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
		}
		avalLmt, err := lmt.IncrLimit(1)
		if err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
		}
		// Получаем инфу о типе отслеживания
		typeInfo, err := subFields.GetTypeInfo()
		if err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
		}
		_, ok = typeInfo.(database.TypeTrackingCrypto)
		if !ok {
			return nil, fmt.Errorf("notificationsCC:error convert interface to struct")
		}
		// Если лимит будет исчерпан, отключаем отслеживание
		if avalLmt == 0 {
			if err := subFields.OffTracking(); err != nil {
				return nil, fmt.Errorf("notificationsCC:" + err.Error())
			}
		}
		// Кешируем пользователя
		if err := database.UsersCache.CheckCache(subFields.UserId); err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
		}
		// Получаем имя юзера
		user, err := database.UsersCache.GetCache(subFields.UserId)
		if err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
		}
		userName, err := user.GetUserName()
		if err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
		}
		// Получаем номер чата с пользователем
		chatIdUsr, err := user.GetChatId()
		if err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
		}
		// Узнаем разность
		diff := dictCryptos.CryptoLastPrice - subFields.ValTrkCrp

		if diff >= 0 && typeInfo.(database.TypeTrackingCrypto).RisingTypTrkCrp { // Поднялась на N под пунктов (пп)
			// Какая-то проверка
		} else if diff < 0 && !typeInfo.(database.TypeTrackingCrypto).RisingTypTrkCrp { // Опустилась на N под пунктов (пп)
			// Какая-то проверка
			diff *= -1
		} else {
			continue
		}

		notifCCStruct = append(notifCCStruct, NotificationsCCStruct{
			subFields.UserId,
			userName,
			chatIdUsr,
			subFields.DctCrpId,
			dictCryptos.CryptoName,
			fmt.Sprintf("Произошло событие над криптовалютой %s:\n"+
				typeInfo.(database.TypeTrackingCrypto).DescTypTrkCrp+
				" на %.3fUSD\nОсталось уведомлений для данного события: %v",
				dictCryptos.CryptoName, subFields.ValTrkCrp, "USD", diff, avalLmt),
		})
		if avalLmt == 0 {
			notifCCStruct = append(notifCCStruct, NotificationsCCStruct{
				subFields.UserId,
				userName,
				chatIdUsr,
				subFields.DctCrpId,
				dictCryptos.CryptoName,
				fmt.Sprint("Вы можете продлить, изменить и продлить данное отслеживание" +
					" или создать новое отслеживание"),
			})
		}
	}

	return notifCCStruct, nil
}
