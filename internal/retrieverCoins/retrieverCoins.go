package retrievercoins

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/caching"
	"github.com/mbydanov/tg_golang_bot/internal/coinmarketcup"
	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/exchange"
	"github.com/mbydanov/tg_golang_bot/internal/models"
)

type quotesLatestAnswerExt struct {
	coinmarketcup.QuotesLatestAnswer
}

func RunRetrieverCoins() error {
	timeout := 300
	var modeluInfo exchange.StatusChannel
	for {
		wg.Wait()
		wg.Add(1)
		res, err := retrieverCoins()
		modeluInfo.Start = true
		modeluInfo.Data = res
		modeluInfo.Error = func() error {
			if len(err) > 0 {
				return err[0]
			}
			return nil
		}()
		modeluInfo.Update = true

		if modeluInfo.Update {
			modeluInfo.Module = models.RetrieverCoins
			exchange.Exchange.WriteChannel(exchange.RetrieverNotification, modeluInfo)
		}
		wg.Done()
		time.Sleep(time.Duration(timeout) * time.Second)
	}
}
func retrieverCoins() (res interface{}, errSl []error) {
	// Получить из кеша данные по криптовалюте
	currencies, err := caching.GetCacheAllRecord(caching.CryptoCache)
	if err != nil {
		return nil, []error{fmt.Errorf("retrieverCoins:" + err.Error())}
	}

	// Строим список валют для запроса
	needFind := make(map[Markets][]string)
	for _, currency := range currencies {
		if currency.Active {
			needFind[Markets(currency.CoinMrktId)] =
				append(needFind[Markets(currency.CoinMrktId)], currency.CryptoName)
		}
	}

	if len(needFind) > 0 {
		for key, value := range needFind {
			switch key {
			case LiveCoinWatch:
				if err = getInfoCoins(value[:], int(key), Map); err != nil {
					errSl = append(errSl, err)
				}
			case CoinMarketCap:
				if err = getInfoCoins(value[:], int(key), QuotesLatest); err != nil {
					errSl = append(errSl, err)
				}
			}
		}

	}

	return nil, errSl
}

func saveFromCMC(responseBody []byte, cryptoCur []string) (cryptoCurOut []string, err error) {

	qla := &quotesLatestAnswerExt{}
	if err = json.Unmarshal([]byte(responseBody), qla); err != nil {
		return nil, fmt.Errorf("saveFromCMC:" + err.Error())
	}
	if qla.Error_code != 0 {
		return nil, fmt.Errorf("saveFromCMC:" + err.Error())
	}
	for i := range qla.QuotesLatestAnswerResults {

		cryptoprices := database.Cryptoprices{
			Timestamp:    time.Now(),
			CryptoId:     qla.QuotesLatestAnswerResults[i].Id,
			CryptoPrice:  qla.QuotesLatestAnswerResults[i].Price,
			CryptoUpdate: qla.QuotesLatestAnswerResults[i].Last_updated,
		}
		cryptoprices, _, err = caching.WriteCache(caching.CryptoPricesCache, cryptoprices.CryptoId, cryptoprices, false)
		if err != nil {
			return nil, fmt.Errorf("saveFromCMC:" + err.Error())
		}

		// Берем запись из кеша
		currency, err := caching.GetCacheByIdxInMap(caching.CryptoCache, qla.QuotesLatestAnswerResults[i].Id)
		if err != nil {
			return nil, fmt.Errorf("saveFromCMC:" + err.Error())
		}

		// Обновляем запись в кеше и БД
		currency.CryptoLastPrice = qla.QuotesLatestAnswerResults[i].Price
		currency.CryptoRank = qla.QuotesLatestAnswerResults[i].Cmc_rank
		currency.CryptoUpdate = qla.QuotesLatestAnswerResults[i].Last_updated
		currency, err = caching.UpdateCacheRecord(caching.CryptoCache, currency.CryptoId, currency)
		if err != nil {
			return nil, fmt.Errorf("saveFromCMC:" + err.Error())
		}

		// Поиск индекса найденной валюты и её удаление из массива needFind
		cryptoCur = models.FindCellAndDelete(cryptoCur, qla.QuotesLatestAnswerResults[i].Symbol)
	}
	if len(cryptoCur) != 0 {
		cryptoCurOut = cryptoCur[:]
	}

	return cryptoCurOut, nil
}

func (qla *quotesLatestAnswerExt) UnmarshalJSON(bs []byte) error {
	var quotesLatest coinmarketcup.QuotesLatest
	if err := json.Unmarshal(bs, &quotesLatest); err != nil {
		return err
	}
	qla.Error_code = quotesLatest.Status.ErrorCode
	qla.Error_message = quotesLatest.Status.Error_message
	for _, value0 := range quotesLatest.Data {
		if len(value0) > 0 {
			qla.QuotesLatestAnswerResults = append(qla.QuotesLatestAnswerResults, coinmarketcup.QuotesLatestAnswerResult{
				Id:           value0[0].Id,
				Name:         value0[0].Name,
				Symbol:       value0[0].Symbol,
				Cmc_rank:     value0[0].Cmc_rank,
				Price:        value0[0].Quote["USD"].Price,
				Currency:     "USD",
				Last_updated: value0[0].Quote["USD"].Last_updated,
			})
		}
	}
	return nil
}
