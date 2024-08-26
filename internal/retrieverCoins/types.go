package retrievercoins

// Для определения мнемоники маркета, чтобы не трогать кеш
type Markets int

// Markets
const (
	LiveCoinWatch Markets = iota + 1
	CoinMarketCap
)

func (m Markets) String() string {
	switch m {
	case LiveCoinWatch:
		return "LCW"
	case CoinMarketCap:
		return "CMC"
	}
	return "N/A"
}

func MarketDesc(idx int) string {
	m := Markets(idx)
	return m.String()
}

// Endpoints
const (
	List int = iota + 1
	Map
)

// Live Coin Watch
type LCWResponseShort struct {
	Coins []LCWCoinsResponseShort
	Error LCWerror
}
type LCWCoinsResponseFull struct {
	Name              string   // coin's name
	Symbol            string   // coin's symbol
	Rank              int      // coin's rank
	Age               int      // coin's age in days
	Color             string   // hexadecimal color code (#282a2a)
	Png32             string   // 32-pixel png image of coin icon
	Png64             string   // 64-pixel png image of coin icon
	Webp32            string   // 32-pixel webp image of coin icon
	Webp64            string   // 64-pixel webpg image of coin icon
	Exchanges         int      // number of exchange coin is present at
	Markets           int      // number of markets coin is present at
	Pairs             int      // number of unique markets coin is present at
	AllTimeHighUSD    float32  // all-time high in USD
	CirculatingSupply int      // number of coins minted, but not locked
	TotalSupply       int      // number of coins minted, including locked
	MaxSupply         int      // maximum number of coins that can be minted
	Categories        []string // array of category string
	LCWCoinsResponseShort
}

type LCWCoinsResponseShort struct {
	Code   string  // coin's code
	Rate   float32 // coin rate in the specified currency
	Volume int     // 24-hour volume of coin
	Cap    int     // market cap of coin
	Delta  LCWdelta
}

type LCWerror struct {
	Code        int
	Status      string
	Description string
}

type LCWdelta struct {
	Hour    float32 // rate of change in the last hour
	Day     float32 // rate of change in the last 24 hours
	Week    float32 // rate of change in the last 7 days
	Month   float32 // rate of change in the last 30 days
	Quarter float32 // rate of change in the last 90 days
	Year    float32 // rate of change in the last 365 days
}
