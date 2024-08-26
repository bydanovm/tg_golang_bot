package retrievercoins

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/caching"
	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/models"
)

func getInfoCoins(cryptoCur []string, marketId int, endpointId int) error {

	request, err := builderRequest(cryptoCur, marketId, endpointId)
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

	switch marketId {
	case int(LiveCoinWatch):
		body := []LCWCoinsResponseShort{}
		if err := json.Unmarshal([]byte(responseBody), &body); err != nil {
			body := LCWerror{}
			if err := json.Unmarshal([]byte(responseBody), &body); err != nil {
				return err
			}
			return err
		}

		for _, item := range body {
			currencyId := caching.GetCacheElementKeyChain(caching.CryptoCache, item.Code).(int)

			cryptoprices := database.Cryptoprices{
				Timestamp:    time.Now(),
				CryptoId:     currencyId,
				CryptoPrice:  item.Rate,
				CryptoUpdate: time.Now(),
			}
			cryptoprices, _, err = caching.WriteCache(caching.CryptoPricesCache, cryptoprices.CryptoId, cryptoprices)
			if err != nil {
				return fmt.Errorf("getInfoCoins:" + err.Error())
			}

			currency, err := caching.GetCacheByIdxInMap(caching.CryptoCache, currencyId)
			if err != nil {
				return fmt.Errorf("getInfoCoins:" + err.Error())
			}
			currency.CryptoLastPrice = item.Rate
			currency.CryptoUpdate = time.Now()

			// Обновить в кеше и БД
			currency, err = caching.UpdateCacheRecord(caching.CryptoCache, currency.CryptoId, currency)
			if err != nil {
				return fmt.Errorf("getInfoCoins:" + err.Error())
			}

			cryptoCur = models.FindCellAndDelete(cryptoCur, currency.CryptoName)

		}

	}

	if len(cryptoCur) != 0 {
		return errors.New(`Криптовалюта ` + strings.Join(cryptoCur, `, `) + ` не найдена`)
	}

	return nil
}

func builderRequest(cryptoCur []string, marketId int, endpointId int) (*http.Request, error) {
	// Определить маркет
	coinMarket, err := caching.GetCacheByIdxInMap(caching.CoinMarketsCache, marketId)
	if err != nil {
		return nil, fmt.Errorf("%s:%s", "builderRequest", err.Error())
	}
	// Определить Endpoint
	coinMarketEndpoint, err := caching.GetCacheByIdxInMap(caching.CoinMarketsEndpointCache, endpointId)
	if err != nil {
		return nil, fmt.Errorf("%s:%s", "builderRequest", err.Error())
	}
	// Определить данные для запроса
	coinMarketsHand, err := caching.GetCacheRecordsKeyChain(caching.CoinMarketsHandCache, coinMarketEndpoint.IdMrktEnd, true)
	if err != nil {
		return nil, fmt.Errorf("%s:%s", "builderRequest", err.Error())
	}
	// Начать построение запроса
	payload := strings.NewReader(`{ ` +
		func() (out string) {
			for k, item := range coinMarketsHand {
				out += `"` + item.Key + `":`
				switch item.Type {
				case "string":
					out += `"` + item.Value + `"`
				case "number":
					out += item.Value
				case "boolean":
					out += item.Value
				case "array": // ["ETH","BTC","TON"]
					if len(cryptoCur) != 0 {
						out += `["` + strings.Join(cryptoCur, `","`) + `"]`
					} else if item.Value != "" {
						out += item.Value
					}
				}
				if k != len(coinMarketsHand)-1 {
					out += `,`
				}
			}
			return out
		}() +
		`}`)

	request, err := http.NewRequest(coinMarketEndpoint.Method, coinMarket.Url+coinMarketEndpoint.Endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("%s:%s", "builderRequest", err.Error())
	}

	// Для CMC
	// query := url.Values{}
	// for _, item := range coinMarketsHand {
	// 	query.Add(item.Key, func() (out string) {
	// 		switch item.Type {
	// 		case "string":
	// 			out = `"` + item.Value + `"`
	// 		case "number":
	// 			out = item.Value
	// 		case "boolean":
	// 			out = item.Value
	// 		case "array": // ["ETH","BTC","TON"]
	// 			if item.Value == "" {
	// 				out = `["` + strings.Join(cryptoCur, `","`) + `"]`
	// 			} else {
	// 				out = item.Value
	// 			}
	// 		}
	// 		return out
	// 	}())
	// }
	// request.URL.RawQuery = query.Encode()

	request.Header.Add(func() (out string) {
		if coinMarket.HeadContentType != "" {
			out = coinMarket.HeadContentType
		} else {
			out = coinMarket.HeadAccept
		}
		return out
	}(), "application/json")
	request.Header.Add(coinMarket.HeadApiKey, os.Getenv("API_LCW"))

	return request, nil
}
