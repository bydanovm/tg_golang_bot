package notifications

import (
	"fmt"

	"github.com/mbydanov/tg_golang_bot/internal/caching"
	"github.com/mbydanov/tg_golang_bot/internal/models"
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
				if res, err := notificationsCC(); err != nil {
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
func notificationsCC() (interface{}, error) {
	notifCCStruct := []NotificationsCCStruct{}

	// Считываем по 100 элементов из кеша отслеживаний
	// Может быть считать по пользователям (меньше суемся в кеш пользователей)?
	trackings, err := caching.GetCacheAllRecord(caching.TrackingCache)
	if err != nil {
		return nil, fmt.Errorf("notificationsCC:" + err.Error())
	}

	for _, tracking := range trackings {
		// Отсеиваем не активные отслеживания
		if !tracking.OnTrkCrp {
			continue
		}

		// Получить из кеша данные по криптовалюте
		currency, err := caching.GetCacheByIdxInMap(caching.CryptoCache, tracking.DctCrpId)
		if err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
		}

		// Получаем информацию о типе отслеживания
		typeInfo, err := caching.GetCacheByIdxInMap(caching.TrackingTypeCache, tracking.TypTrkCrpId)
		if err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
		}

		// Получаем лимиты
		lmt, err := caching.GetCacheByIdxInMap(caching.Limits, tracking.LmtId)
		if err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
		}

		// Кешируем пользователя заранее, чтобы не было очистки кеша
		user, err := caching.GetCacheByIdxInMap(caching.UsersCache, tracking.UserId)
		if err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
		}

		// Узнаем разность
		diff := currency.CryptoLastPrice - tracking.ValTrkCrp
		if diff >= 0 && typeInfo.RisingTypTrkCrp { // Поднялась на N под пунктов (пп)
			// Какая-то проверка
		} else if diff < 0 && !typeInfo.RisingTypTrkCrp { // Опустилась на N под пунктов (пп)
			// Какая-то проверка
			diff *= -1
		} else {
			continue
		}

		// Увеличиваем использованные лимиты и обновляем в кеше/бд
		lmt.ValUsedLmt++
		lmt, err = caching.UpdateCacheRecord(caching.Limits, lmt.IdLmt, lmt)
		if err != nil {
			return nil, fmt.Errorf("notificationsCC:" + err.Error())
		}

		// Если лимит будет исчерпан, отключаем отслеживание
		avalLmt := lmt.ValAvailLmt - lmt.ValUsedLmt
		if avalLmt == 0 {
			tracking.OnTrkCrp = false
			tracking, err = caching.UpdateCacheRecord(caching.TrackingCache, tracking.IdTrkCrp, tracking)
			if err != nil {
				return nil, fmt.Errorf("notificationsCC:" + err.Error())
			}
		}

		notifCCStruct = append(notifCCStruct, NotificationsCCStruct{
			tracking.UserId,
			user.NameUsr,
			user.ChatIdUsr,
			tracking.DctCrpId,
			currency.CryptoName,
			fmt.Sprintf("Произошло событие над криптовалютой %s:\n"+
				typeInfo.DescTypTrkCrp+
				" на %.3fUSD\nОсталось уведомлений для данного события: %v",
				currency.CryptoName, tracking.ValTrkCrp, "USD", diff, avalLmt),
		})
		if avalLmt == 0 {
			notifCCStruct = append(notifCCStruct, NotificationsCCStruct{
				tracking.UserId,
				user.NameUsr,
				user.ChatIdUsr,
				tracking.DctCrpId,
				currency.CryptoName,
				fmt.Sprint("Вы можете продлить, изменить и продлить данное отслеживание" +
					" или создать новое отслеживание"),
			})
		}
	}

	return notifCCStruct, nil
}
