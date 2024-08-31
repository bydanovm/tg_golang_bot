package retrievercoins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/caching"
	"github.com/mbydanov/tg_golang_bot/internal/services"
	"github.com/sirupsen/logrus"
)

func RunUpdaterRank() {
	for {
		wg.Wait()
		wg.Add(1)
		if err := updaterRank(); err != nil {
			services.Logging.WithFields(logrus.Fields{
				"module": "updaterRank",
				"error": func(err error) interface{} {
					// if err
					return err
				}(err),
			}).Error()
		}
		wg.Done()
		<-time.After(time.Minute * 60)
	}
}

func updaterRank() error {
	request, err := builderRequest([]string{}, int(CoinMarketCap), 3)
	if err != nil {
		return fmt.Errorf("%s:%s", "getInfoCoins", err.Error())
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("%s:%s", "getInfoCoins", err.Error())
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("%s:%s", "getInfoCoins", err.Error())
	}

	body := CryptoCyrrencyMap{}
	if err := json.Unmarshal([]byte(responseBody), &body); err != nil {
		return err
	}
	if body.Status.Error_code != 0 {
		return fmt.Errorf("%s:%d:%s", "getInfoCoins", body.Status.Error_code, body.Status.Error_message)
	}
	// Если нет ошибок опроса, то стираем ранги в БД без изменения кеша
	currency, err := caching.GetCacheAllRecord(caching.CryptoCache)
	if err != nil {

	}
	for _, item := range currency {
		item.CryptoRank = 0
		item.Active = false
		_, err := caching.UpdateCacheRecord(caching.CryptoCache, item.CryptoId, item, false)
		if err != nil {
			// Если ошибка обновления, то логировать и продолжить дальше или бить ошибку?
			// return fmt.Errorf("getAndSaveFromAPI:" + err.Error())
		}
	}
	for key, item := range body.Data {
		currency, err := caching.GetCacheByIdxInMap(caching.CryptoCache, item.Id)
		if err != nil {
			err = nil
			currency.CryptoId = item.Id
			currency.CoinMrktId = int(LiveCoinWatch)
			// Если валюта не найдена в базе, нужно ее добавить
			// return fmt.Errorf("getAndSaveFromAPI:" + err.Error())
		}
		currency.Active = func() (active bool) {
			if item.Is_active > 0 {
				active = true
			}
			return active
		}()
		currency.CryptoRank = key + 1
		currency.CryptoName = item.Symbol
		currency.CryptoUpdate = body.Status.Timestamp
		currency.Slug = item.Slug
		currency.CryptoFullName = item.Name
		currency, err = caching.UpdateCacheRecord(caching.CryptoCache, currency.CryptoId, currency)
		if err != nil {
			// Если ошибка обновления, то логировать и продолжить дальше
			// return fmt.Errorf("getAndSaveFromAPI:" + err.Error())
		}
	}

	return nil
}
