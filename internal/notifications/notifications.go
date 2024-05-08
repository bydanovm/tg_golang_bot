package notifications

import (
	"fmt"

	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/models"
	"github.com/mitchellh/mapstructure"
)

func RunNotification(
	chanModules chan models.StatusChannel) {
	for {
		// Считывание информации из канала от ретривера
		v, ok := <-chanModules
		if ok {
			// Прием ответа от ретривера
			if v.Module == models.RetrieverCoins && v.Start {
				// Функция работы с уведомлениями по КВ
				if res, err := notificationsCC(database.DCCache); err != nil {
					// Запись в канал об ошибке
					v.Error = err
				} else if res != nil {
					// Здесь будет запись в канал
					v.Start = true
					v.Data = res
				} else {
					// Если ничего не найдено
					v.Start = false
					v.Data = nil
				}
				v.Module = models.Notificator
			}
			// Если ответ не от ретривера то записать инфу обратно в канал
			chanModules <- v
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

	for _, subRs := range database.TCCache {
		subFields := database.TrackingCrypto{}
		mapstructure.Decode(subRs, &subFields)
		// Отсеиваем не активные отслеживания
		if !subFields.OnTrkCrp {
			continue
		}

		dictCryptos := database.DictCrypto{}
		v, ok := bufferForNotif.(database.DictCryptoCache)
		if !ok {
			return nil, fmt.Errorf("notificationsCC:error convert interface to map[int]interface")
		}
		mapstructure.Decode(v[subFields.DctCrpId], &dictCryptos)

		// Получаем информацию о типе отслеживания
		typeInfo, err := database.TypeTCCache.GetCache(subFields.TypTrkCrpId)
		if err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
		}

		// Узнаем разность
		diff := dictCryptos.CryptoLastPrice - subFields.ValTrkCrp
		if diff >= 0 && typeInfo.RisingTypTrkCrp { // Поднялась на N под пунктов (пп)
			// Какая-то проверка
		} else if diff < 0 && !typeInfo.RisingTypTrkCrp { // Опустилась на N под пунктов (пп)
			// Какая-то проверка
			diff *= -1
		} else {
			continue
		}

		// Получаем лимит в соответствии с отслеживанием и увеличиваем его
		lmt := database.LmtCache[subFields.LmtId]
		avalLmt, err := lmt.IncrLimit(1)
		if err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
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

		notifCCStruct = append(notifCCStruct, NotificationsCCStruct{
			subFields.UserId,
			userName,
			chatIdUsr,
			subFields.DctCrpId,
			dictCryptos.CryptoName,
			fmt.Sprintf("Произошло событие над криптовалютой %s:\n"+
				typeInfo.DescTypTrkCrp+
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
