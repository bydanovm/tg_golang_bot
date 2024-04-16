package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/database"
	"github.com/mbydanov/tg_golang_bot/internal/models"
	"github.com/mitchellh/mapstructure"
)

func GetConfig(mainConfigIn chan models.StatusChannel) models.StatusChannel {
	var cfg ConfigStruct
	var stCh models.StatusChannel

	expLst := []database.Expressions{
		{Key: database.Name, Operator: database.NotEQ, Value: `'` + database.Empty + `'`},
	}

	for {
		rs, find, _, err := database.ReadDataRow(&database.SettingsProject{}, expLst, 1)
		if err != nil {
			stCh.Error = fmt.Errorf("GetConfig:" + err.Error())
			mainConfigIn <- stCh
		}

		if find {
			for _, subRs := range rs {
				subFields := database.SettingsProject{}
				mapstructure.Decode(subRs, &subFields)
				// Таймер опроса ретривера
				if subFields.Name == TMR_RESP_RTV && subFields.Active {
					cfg.TmrRespRvt, err = strconv.Atoi(subFields.Value)
					if err != nil {
						// cfg.MsgError = fmt.Errorf("GetConfig:Atoi:" + err.Error())
					}
				}
			}
			// cfgChan <- cfg
		}
		time.Sleep(60 * time.Second)
	}
}
