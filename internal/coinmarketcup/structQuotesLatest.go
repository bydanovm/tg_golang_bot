package coinmarketcup

import (
	"time"

	"github.com/mbydanov/tg_golang_bot/internal/database"
)

type QuotesLatest struct {
	Status status
	Data   map[string][]cryptocurrencObject
}

type status struct {
	Timestamp     string
	ErrorCode     int
	Error_message string
	Elapsed       int
	CreditCount   int
	Notice        string
}

type cryptocurrencObject struct {
	Id                 int
	Name               string
	Symbol             string
	Slug               string
	Num_market_pairs   int
	Date_added         string
	Max_supply         float32
	Circulating_supply float32
	Total_supply       float32
	Is_active          int
	Infinite_supply    bool
	Cmc_rank           int
	Is_fiat            int
	Last_updated       string
	Quote              map[string]currency
}

type currency struct {
	Price                    float32
	Volume_24h               float32
	Volume_change_24h        float32
	Percent_change_1h        float32
	Percent_change_24h       float32
	Percent_change_7d        float32
	Percent_change_30d       float32
	Percent_change_60d       float32
	Percent_change_90d       float32
	Market_cap               float32
	Market_cap_dominance     float32
	Fully_diluted_market_cap float32
	Last_updated             time.Time
}

type QuotesLatestAnswer struct {
	Error_code                int
	Error_message             string
	QuotesLatestAnswerResults []QuotesLatestAnswerResult
}

type QuotesLatestAnswerResult struct {
	Id           int
	Name         string
	Symbol       string
	Cmc_rank     int
	Price        float32
	Currency     string
	Last_updated time.Time
}
type GetLatestObject struct {
	Crypto database.DictCrypto
	Find   bool
}
