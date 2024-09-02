package retrievercoins

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
			cryptoprices, _, err = caching.WriteCache(caching.CryptoPricesCache, cryptoprices.CryptoId, cryptoprices, false)
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

	case int(CoinMarketCap):
		if cryptoCur, err = saveFromCMC(responseBody, cryptoCur); err != nil {
			return fmt.Errorf("getInfoCoins:" + err.Error())
		}
	}

	if len(cryptoCur) != 0 {
		return errors.New(`Криптовалюта ` + strings.Join(cryptoCur, `, `) + ` не найдена`)
	}

	return nil
}

func builderRequest(cryptoCur []string, marketId int, endpointId int) (request *http.Request, err error) {
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
	coinMarketsHand, err := caching.GetCacheRecordsKeyChain(caching.CoinMarketsHandCache, coinMarketEndpoint.IdMrktEnd, false)
	if err != nil {
		return nil, fmt.Errorf("%s:%s", "builderRequest", err.Error())
	}

	if coinMarket.IdMrkt == int(LiveCoinWatch) {
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

		request, err = http.NewRequest(coinMarketEndpoint.Method, coinMarket.Url+coinMarketEndpoint.Endpoint, payload)
		if err != nil {
			return nil, fmt.Errorf("%s:%s", "builderRequest", err.Error())
		}
	} else if coinMarket.IdMrkt == int(CoinMarketCap) {
		query := url.Values{}
		for _, item := range coinMarketsHand {
			switch item.Type {
			case "array":
				if len(cryptoCur) > 0 {
					query.Add(item.Key, strings.Join(cryptoCur, ","))
				} else if item.Value != "" {
					query.Add(item.Key, item.Value)
				}
			default:
				if item.Value != "" {
					query.Add(item.Key, item.Value)
				}
			}
		}
		request, err = http.NewRequest(coinMarketEndpoint.Method, coinMarket.Url+coinMarketEndpoint.Endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("%s:%s", "builderRequest", err.Error())
		}
		request.URL.RawQuery = query.Encode()
	}

	request.Header.Add(func() (out string) {
		if coinMarket.HeadContentType != "" {
			out = coinMarket.HeadContentType
		} else {
			out = coinMarket.HeadAccept
		}
		return out
	}(), "application/json")
	request.Header.Add(coinMarket.HeadApiKey, os.Getenv("API_"+coinMarket.Description))

	return request, nil
}
