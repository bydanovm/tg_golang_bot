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

	"github.com/mbydanov/tg_golang_bot/internal/coinmarketcup"
	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/models"
	"github.com/mitchellh/mapstructure"
)

type quotesLatestAnswerExt struct {
	coinmarketcup.QuotesLatestAnswer
}

func RunRetrieverCoins(
	chanModules chan models.StatusChannel) error {
	timeout := 600
	var modeluInfo models.StatusChannel
	for {
		if res, err := retrieverCoins(); err != nil {
			modeluInfo.Error = err
			modeluInfo.Update = true
		} else {
			// Передача данных в нотификатор
			modeluInfo.Start = true
			modeluInfo.Data = res
			modeluInfo.Update = true
		}
		if modeluInfo.Update {
			modeluInfo.Module = models.RetrieverCoins
			chanModules <- modeluInfo
		}
		time.Sleep(time.Duration(timeout) * time.Second)
	}
}
func retrieverCoins() (interface{}, error) {
	// Получаем основные КВ для запроса актуальной информации из справочника КВ
	fields := database.DictCrypto{}
	expLst := []database.Expressions{}

	expLst = append(expLst, database.Expressions{
		Key: database.CryptoName, Operator: database.NotEQ, Value: `'` + database.Empty + `'`,
	})

	rs, find, _, err := database.ReadDataRow(&fields, expLst, 0)
	if err != nil {
		return nil, err
	}
	// Если какие-то записи найдены, то мы строим запрос для обращения к API
	if find {
		var needFind []string
		for _, subRs := range rs {
			subFields := database.DictCrypto{}
			mapstructure.Decode(subRs, &subFields)
			if subFields.Active {
				needFind = append(needFind, subFields.CryptoName)
			}
		}
		if len(needFind) > 0 {
			if res, err := getAndSaveFromAPI(needFind); err != nil {
				return res, err
			} else {
				return res, err
			}
		}
	}
	return nil, nil
}

func getAndSaveFromAPI(cryptoCur []string) (interface{}, error) {
	bufferForNotif := database.DictCrypto{}        // Буфер для посыла в нотификатор
	bufferForNotifMap := make(map[int]interface{}) // Буфер для посыла в нотификатор
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("symbol", strings.Join(cryptoCur, ","))
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", os.Getenv("API_CMC"))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, _ := io.ReadAll(resp.Body)
	qla := &quotesLatestAnswerExt{}
	if err = json.Unmarshal([]byte(respBody), qla); err != nil {
		return nil, err
	}
	if qla.Error_code != 0 {
		return nil, err
	}
	for i := range qla.QuotesLatestAnswerResults {

		dateTime, err := models.ConvertDateTimeToMSK(qla.QuotesLatestAnswerResults[i].Last_updated)
		if err != nil {
			return nil, fmt.Errorf("getAndSaveFromAPI:" + err.Error())
		}
		// dateTimeUTC3, _ := time.ParseInLocation(layout, dateTime, dateTimeLocUTC3)
		// Добавление найденной валюты в таблицы текущих цен и обновление справочника валют
		cryptoprices := map[string]string{
			"CryptoId":     fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Id),
			"CryptoPrice":  fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Price),
			"CryptoUpdate": dateTime,
		}
		if err := database.WriteData("cryptoprices", cryptoprices); err != nil {
			return nil, err
		}

		dictCryptos := map[string]string{
			// "CryptoId":        fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Id),
			// "CryptoName":      fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Symbol),
			"CryptoLastPrice": fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Price),
			"CryptoUpdate":    dateTime,
		}
		expLst := []database.Expressions{}
		expLst = append(expLst, database.Expressions{
			Key: database.CryptoId, Operator: database.EQ, Value: `'` + cryptoprices["CryptoId"] + `'`,
		})
		if err := database.UpdateData("dictcrypto", dictCryptos, expLst); err != nil {
			return nil, err
		} else {
			// Если не было ошибки при обновлении, то кешируем
			d := database.DCCache[qla.QuotesLatestAnswerResults[i].Id]
			d.CryptoLastPrice = qla.QuotesLatestAnswerResults[i].Price
			// d.CryptoUpdate = dateTime
			database.DCCache[qla.QuotesLatestAnswerResults[i].Id] = d
		}

		// Поиск индекса найденной валюты и её удаление из массива needFind
		cryptoCur = models.FindCellAndDelete(cryptoCur, qla.QuotesLatestAnswerResults[i].Symbol)
		// Добавление в буфер
		bufferForNotif = database.DictCrypto{
			Id:              0,
			Timestamp:       time.Now(),
			CryptoId:        qla.QuotesLatestAnswerResults[i].Id,
			CryptoName:      qla.QuotesLatestAnswerResults[i].Symbol,
			CryptoLastPrice: qla.QuotesLatestAnswerResults[i].Price,
			CryptoUpdate:    time.Now(),
			Active:          true,
			CryptoCounter:   0,
		}
		bufferForNotifMap[qla.QuotesLatestAnswerResults[i].Id] = bufferForNotif
	}
	// Есть не найденная криптовалюта
	if len(cryptoCur) != 0 {
		// return bufferForNotif, errors.New(`Криптовалюта ` + strings.Join(cryptoCur, `, `) + ` не найдена`)
		return bufferForNotifMap, errors.New(`Криптовалюта ` + strings.Join(cryptoCur, `, `) + ` не найдена`)
	}
	// return bufferForNotif, nil
	return bufferForNotifMap, nil
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
