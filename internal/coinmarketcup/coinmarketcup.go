package coinmarketcup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/models"
	"github.com/mitchellh/mapstructure"
)

func GetLatest(cryptocurrencies string) (answer []string) {
	s := make([]string, 0)
	var needFind []string
	// Обрабатываем входную строку, преобразовываем в массив
	// cryptoCur := strings.Split(cryptocurrencies, ", ")
	cryptoCur := strings.FieldsFunc(cryptocurrencies, func(r rune) bool {
		return r == ',' || r == ' '
	})
	for i := 0; i < len(cryptoCur); i++ {
		cryptoCur[i] = strings.ToUpper(strings.Trim(cryptoCur[i], ` !&.,@#$%^*()-_=+/\?<>{}АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя`))
	}
	// Проверка на пустой массив, если пустой, то удаляем
	cryptoCur = models.ChkArrayBySpace(cryptoCur)
	// Проверяем наличие криптовалюты в БД
	fields := database.DictCrypto{}
	expLst := []database.Expressions{}

	expLst = append(expLst, database.Expressions{
		Key: database.CryptoName, Operator: "IN", Value: "('" + strings.Join(cryptoCur, "','") + "')",
	})

	rs, find, countFind, err := database.ReadDataRow(&fields, expLst, len(cryptoCur))
	if err != nil {
		s = append(s, "Возвращена ошибка при поиске в БД: "+err.Error())
		return s
	}
	// Если запись найдена, возвращаем из БД
	if find {
		var findCryptoCur []string
		for _, subRs := range rs {
			subFields := database.DictCrypto{}
			mapstructure.Decode(subRs, &subFields)
			str := fmt.Sprintf("Криптовалюта: %s\nЦена: %.9f %s\nОбновлено: %s",
				subFields.CryptoName,
				subFields.CryptoLastPrice,
				"USD",
				subFields.CryptoUpdate.Format("2006-01-02 15:04:05"),
			)
			findCryptoCur = append(findCryptoCur, subFields.CryptoName)
			s = append(s, str)
			// Увеличиваем счетчик запроса КВ
			dictCryptos := map[string]string{
				"CryptoCounter": fmt.Sprintf("%v", subFields.CryptoCounter+1),
			}
			expLst := []database.Expressions{}
			expLst = append(expLst, database.Expressions{
				Key: database.CryptoId, Operator: database.EQ, Value: `'` + fmt.Sprintf("%v", subFields.CryptoId) + `'`,
			})
			if err := database.UpdateData("dictcrypto", dictCryptos, expLst); err != nil {
				s = append(s, "Возвращена ошибка при обновлении в БД: "+err.Error())
				return s
			} else {
				d := database.DCCache[subFields.CryptoId].(database.DictCrypto)
				d.CryptoCounter = subFields.CryptoCounter + 1
				database.DCCache[subFields.CryptoId] = d
			}
		}
		// Если нашли все валюты, то возвращаем их
		if countFind == len(cryptoCur) {
			return s
		}
		// Если не все найдены, то определяем какие валюты мы не нашли

		for _, v1 := range cryptoCur {
			for i2 := 0; i2 < len(findCryptoCur); i2++ {
				// for i2, v2 := range findCryptoCur {
				if v1 == findCryptoCur[i2] {
					break
				}
				if v1 != findCryptoCur[i2] && i2 == len(findCryptoCur)-1 {
					needFind = append(needFind, v1)
				}
			}
		}
	} else {
		needFind = cryptoCur
	}
	// Если валюту в БД не нашли, производим поиск не найденной валюты посредством API
	// И добавление найденной в словарь
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		s = append(s, "Возвращена ошибка:\n"+err.Error())
	}

	q := url.Values{}
	q.Add("symbol", strings.Join(needFind, ","))
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", os.Getenv("API_CMC"))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		s = append(s, "Возвращена ошибка:\n"+err.Error())
	}
	respBody, _ := io.ReadAll(resp.Body)
	qla := &QuotesLatestAnswer{}
	if err = json.Unmarshal([]byte(respBody), qla); err != nil {
		s = append(s, "Возвращена ошибка:\n"+err.Error())
	}
	if qla.Error_code != 0 {
		s = append(s, "Возвращена ошибка:\n"+qla.Error_message)
	}
	for i := range qla.QuotesLatestAnswerResults {
		dateTime, err := models.ConvertDateTimeToMSK(qla.QuotesLatestAnswerResults[i].Last_updated)
		if err != nil {
			s = append(s, fmt.Sprintf("getAndSaveFromAPI:"+err.Error()))
		}

		str := fmt.Sprintf("Криптовалюта: %s\nЦена: %.9f %s\nОбновлено: %s",
			qla.QuotesLatestAnswerResults[i].Symbol,
			qla.QuotesLatestAnswerResults[i].Price,
			qla.QuotesLatestAnswerResults[i].Currency,
			dateTime,
		)
		s = append(s, str)

		// Добавление найденной валюты в БД текущих цен и справочник валют
		cryptoprices := map[string]string{
			"CryptoId":     fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Id),
			"CryptoPrice":  fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Price),
			"CryptoUpdate": dateTime,
		}
		dictCryptos := map[string]string{
			"CryptoId":        fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Id),
			"CryptoName":      fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Symbol),
			"CryptoLastPrice": fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Price),
			"CryptoUpdate":    dateTime,
		}
		if err := database.WriteData("dictcrypto", dictCryptos); err != nil {
			s = append(s, fmt.Sprintf("GetLatest:"+err.Error()))
		} else {
			d := database.DictCrypto{
				CryptoId:        qla.QuotesLatestAnswerResults[i].Id,
				CryptoName:      qla.QuotesLatestAnswerResults[i].Symbol,
				CryptoLastPrice: qla.QuotesLatestAnswerResults[i].Price,
				CryptoCounter:   1,
			}
			database.DCCache[d.CryptoId] = d
		}
		if err := database.WriteData("cryptoprices", cryptoprices); err != nil {
			s = append(s, "Возвращена ошибка:\n"+err.Error())
		}
		// Поиск индекса найденной валюты и её удаление из массива needFind
		needFind = models.FindCellAndDelete(needFind, qla.QuotesLatestAnswerResults[i].Symbol)

	}
	// Есть не найденная криптовалюта
	if len(needFind) != 0 {
		s = append(s, "Криптовалюта "+strings.Join(needFind, `, `)+" не найдена")
	}
	return s
}

func (qla *QuotesLatestAnswer) UnmarshalJSON(bs []byte) error {
	var quotesLatest QuotesLatest
	if err := json.Unmarshal(bs, &quotesLatest); err != nil {
		return err
	}
	qla.Error_code = quotesLatest.Status.ErrorCode
	qla.Error_message = quotesLatest.Status.Error_message
	for _, value0 := range quotesLatest.Data {
		if len(value0) > 0 {
			qla.QuotesLatestAnswerResults = append(qla.QuotesLatestAnswerResults, QuotesLatestAnswerResult{
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
