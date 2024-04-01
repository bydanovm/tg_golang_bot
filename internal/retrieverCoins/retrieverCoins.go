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

func RunRetrieverCoins(timeout int, errorMsg chan models.StatusRetriever) error {
	var chanSrv models.StatusRetriever
	for {
		if err := retrieverCoins(); err != nil {
			chanSrv.MsgError = err
			errorMsg <- chanSrv
		}
		// if timeout == 0 {
		timeout = 300
		// }
		time.Sleep(time.Duration(timeout) * time.Second)
	}
}
func retrieverCoins() error {
	// Получаем основные КВ для запроса актуальной информации из справочника КВ
	fields := database.DictCrypto{}
	expLst := []database.Expressions{}

	expLst = append(expLst, database.Expressions{
		Key: database.CryptoName, Operator: database.NotEQ, Value: `'` + database.Empty + `'`,
	})

	rs, find, _, err := database.ReadDataRow(&fields, expLst, 0)
	if err != nil {
		return err
	}
	// Если какие-то записи найдены, то мы строим запрос для обращения к API
	if find {
		var needFind []string
		for _, subRs := range rs {
			subFields := database.DictCrypto{}
			mapstructure.Decode(subRs, &subFields)
			needFind = append(needFind, subFields.CryptoName)
		}
		if len(needFind) > 0 {
			if err := getAndSaveFromAPI(needFind); err != nil {
				return err
			}
		}
	}
	return nil
}

func getAndSaveFromAPI(cryptoCur []string) error {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		return err
	}

	q := url.Values{}
	q.Add("symbol", strings.Join(cryptoCur, ","))
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", os.Getenv("API_CMC"))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	respBody, _ := io.ReadAll(resp.Body)
	qla := &quotesLatestAnswerExt{}
	if err = json.Unmarshal([]byte(respBody), qla); err != nil {
		return err
	}
	if qla.Error_code != 0 {
		return err
	}
	for i := range qla.QuotesLatestAnswerResults {

		dateTime, err := models.ConvertDateTimeToMSK(qla.QuotesLatestAnswerResults[i].Last_updated)
		if err != nil {
			return fmt.Errorf("getAndSaveFromAPI:" + err.Error())
		}
		// dateTimeUTC3, _ := time.ParseInLocation(layout, dateTime, dateTimeLocUTC3)
		// Добавление найденной валюты в таблицы текущих цен и обновление справочника валют
		cryptoprices := map[string]string{
			"CryptoId":     fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Id),
			"CryptoPrice":  fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Price),
			"CryptoUpdate": dateTime,
		}
		if err := database.WriteData("cryptoprices", cryptoprices); err != nil {
			return err
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
			return err
		}

		// Поиск индекса найденной валюты и её удаление из массива needFind
		cryptoCur = models.FindCellAndDelete(cryptoCur, qla.QuotesLatestAnswerResults[i].Symbol)

	}
	// Есть не найденная криптовалюта
	if len(cryptoCur) != 0 {
		return errors.New(`Криптовалюта ` + strings.Join(cryptoCur, `, `) + ` не найдена`)
	}
	return nil
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
