package models

const (
	RetrieverCoins string = `RetrieverCoins`
	Notificator    string = `Notificator`
	CoinMarketCap  string = `CoinMarketCap`
)

type StatusRetriever struct {
	MsgError error
}

type StatusChannel struct {
	Module string
	Start  bool
	Stop   bool
	Update bool
	Error  error
	Data   interface{}
}
