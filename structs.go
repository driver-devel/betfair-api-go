package betfair

import (
	"crypto/tls"
	"time"
)

type LoginMethod int

const (
	Interactive LoginMethod = iota
	NoneInteractive
)

type Account struct {
	Username       string
	Password       string
	ApplicationKey string
	Certificate    tls.Certificate
	KeepAlive      bool
	LoginMethod    LoginMethod
}

type TimeRange struct {
	From time.Time `json:"from,omitempty"`
	To   time.Time `json:"to,omitempty"`
}

type MarketFilter struct {
	TextQuery          string     `json:"textQuery,omitempty"`
	ExchangeIDs        []string   `json:"exchangeIds,omitempty"`
	EventTypeIDs       []string   `json:"eventTypeIds,omitempty"`
	EventIDs           []string   `json:"eventIds,omitempty"`
	CompetitionIDs     []string   `json:"competitionIds,omitempty"`
	MarketIDs          []string   `json:"marketIds,omitempty"`
	Venues             []string   `json:"venues,omitempty"`
	BspOnly            *bool      `json:"bspOnly,omitempty"`
	TurnInPlayEnabled  *bool      `json:"turnInPlayEnabled,omitempty"`
	InPlayOnly         *bool      `json:"inPlayOnly,omitempty"`
	MarketBettingTypes []string   `json:"marketBettingTypes,omitempty"`
	MarketCountries    []string   `json:"marketCountries,omitempty"`
	MarketTypeCodes    []string   `json:"marketTypeCodes,omitempty"`
	MarketStartTime    *TimeRange `json:"marketStartTime,omitempty"`
	WithOrders         []string   `json:"withOrders,omitempty"`
}

type InteractiveSessionResponse struct {
	Token   string
	Product string
	Status  string
	Error   string
}

type NoneInteractiveSessionResponse struct {
	SessionToken string `json:"sessionToken"`
	LoginStatus  string `json:"loginStatus"`
}

type apiResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type apiResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	Error   apiResponseError `json:"error"`
	Result  interface{}      `json:"result"`
}

type apiRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type keepAliveResult struct {
	Token   string `json:"token"`
	Product string `json:"product"`
	Status  string `json:"status"`
	Error   string `json:"error"`
}

type EventType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EventTypeResult struct {
	MarketCount int64     `json:"marketCount"`
	EventType   EventType `json:"eventType"`
}

type Competition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CompetitionResult struct {
	MarketCount int64       `json:"marketCount"`
	Competition Competition `json:"competition"`
}

type Event struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CountryCode string    `json:"countryCode"`
	Timezone    string    `json:"timezone"`
	OpenDate    time.Time `json:"openDate"`
}

type EventResult struct {
	MarketCount int64 `json:"marketCount"`
	Event       Event `json:"event"`
}

type CountryResult struct {
	MarketCount int64  `json:"marketCount"`
	CountryCode string `json:"countryCode"`
}

type VenueResult struct {
	MarketCount int64  `json:"marketCount"`
	Venue       string `json:"venue"`
}

type MarketDescription struct {
	PersistenceEnabled bool    `json:"persistenceEnabled"`
	BspMarket          bool    `json:"bspMarket"`
	MarketTime         string  `json:"marketTime"`
	SuspendTime        string  `json:"suspendTime"`
	SettleTime         string  `json:"settleTime"`
	BettingType        string  `json:"bettingType"`
	TurnInPlayEnabled  bool    `json:"turnInPlayEnabled"`
	MarketType         string  `json:"marketType"`
	Regulator          string  `json:"regulator"`
	MarketBaseRate     float64 `json:"marketBaseRate"`
	DiscountAllowed    bool    `json:"discountAllowed"`
	Wallet             string  `json:"wallet"`
	Rules              string  `json:"rules"`
	RulesHasDate       bool    `json:"rulesHasDate"`
	Certifications     string  `json:"certifications"`
}

type RunnerCatalogue struct {
	SelectionID  int64                  `json:"selectionId"`
	RunnerName   string                 `json:"runnerName"`
	Handicap     float64                `json:"handicap"`
	SortPriority int64                  `json:"sortPriority"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type MarketCatalogue struct {
	MarketID        string             `json:"marketId"`
	MarketName      string             `json:"marketName"`
	MarketStartTime string             `json:"marketStartTime"`
	Description     *MarketDescription `json:"description"`
	Runners         []RunnerCatalogue  `json:"runners"`
	TotalMatched    float64            `json:"totalMatched"`
	EventType       *EventType         `json:"eventType"`
	Competition     *Competition       `json:"competition"`
	Event           *Event             `json:"event"`
}

type PriceSize struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

type StartingPrices struct {
	NearPrice         string      `json:"nearPrice"`
	FarPrice          string      `json:"farPrice"`
	BackStakeTaken    []PriceSize `json:"backStakeTaken"`
	LayLiabilityTaken []PriceSize `json:"layLiabilityTaken"`
	ActualSP          float64     `json:"actualSP"`
}

type ExchangePrices struct {
	AvailableToBack []PriceSize `json:"availableToBack"`
	AvailableToLay  []PriceSize `json:"availableToLay"`
	TradedVolume    []PriceSize `json:"tradedVolume"`
}

type RunnerOrder struct {
	BetID     string `json:"betId"`
	OrderType string `json:"orderType"`
	Status    string `json:"status"`
}

type RunnerMatch struct {
	BetID     string  `json:"betId"`
	MatchID   string  `json:"matchId"`
	Side      string  `json:"side"`
	Price     float64 `json:"price"`
	Size      float64 `json:"size"`
	MatchDate string  `json:"matchDate"`
}

type Runner struct {
	SelectionID      int64           `json:"selectionId"`
	Handicap         float64         `json:"handicap"`
	Status           string          `json:"status"`
	AdjustmentFactor float64         `json:"adjustmentFactor"`
	LastPriceTraded  float64         `json:"lastPriceTraded"`
	TotalMatched     float64         `json:"totalMatched"`
	RemovalDate      string          `json:"removalDate"`
	SP               *StartingPrices `json:"sp"`
	EX               *ExchangePrices `json:"ex"`
	Orders           []RunnerOrder   `json:"orders"`
	Matches          []RunnerMatch   `json:"matches"`
}

type MarketBook struct {
	MarketID              string     `json:"marketId"`
	IsMarketDataDelayed   bool       `json:"isMarketDataDelayed"`
	Status                string     `json:"status"`
	BetDelay              int64      `json:"betDelay"`
	Runners               []Runner   `json:"runners"`
	BspReconciled         bool       `json:"bspReconciled"`
	Complete              bool       `json:"complete"`
	Inplay                bool       `json:"inplay"`
	NumberOfWinners       int64      `json:"numberOfWinners"`
	NumberOfRunners       int64      `json:"numberOfRunners"`
	NumberOfActiveRunners int64      `json:"numberOfActiveRunners"`
	LastMatchTime         *time.Time `json:"lastMatchTime"`
	TotalMatched          float64    `json:"totalMatched"`
	TotalAvailable        float64    `json:"totalAvailable"`
	CrossMatching         bool       `json:"crossMatching"`
	RunnersVoidable       bool       `json:"runnersVoidable"`
	Version               int64      `json:"version"`
}

type ExBestOffersOverrides struct {
	BestPricesDepth          int64   `json:"bestPricesDepth,omitempty"`
	RollupModel              string  `json:"rollupModel,omitempty"`
	RollupLimit              int64   `json:"rollupLimit,omitempty"`
	RollupLiabilityThreshold float64 `json:"rollupLiabilityThreshold,omitempty"`
	RollupLiabilityFactor    int64   `json:"rollupLiabilityFactor,omitempty"`
}

type PriceProjection struct {
	PriceData             []string               `json:"priceData"`
	ExBestOffersOverrides *ExBestOffersOverrides `json:"exBestOffersOverrides"`
	Virtualise            bool                   `json:"virtualise,omitempty"`
	RolloverStakes        bool                   `json:"rolloverStakes,omitempty"`
}

type CurrentOrderSumary struct {
	BetID               string    `json:"betId"`
	MarketID            string    `json:"marketId"`
	SelectionID         int64     `json:"selectionId"`
	Handicap            float64   `json:"handicap"`
	PriceSize           PriceSize `json:"priceSize"`
	BSPLiability        float64   `json:"bspLiability"`
	Side                string    `json:"side"`
	Status              string    `json:"status"`
	PersistenceType     string    `json:"persistenceType"`
	OrderType           string    `json:"orderType"`
	PlacedDate          string    `json:"placedDate"`
	MatchedDate         string    `json:"matchedDate"`
	AveragePriceMatched float64   `json:"averagePriceMatched"`
	SizeMatched         float64   `json:"sizeMatched"`
	SizeRemaining       float64   `json:"sizeRemaining"`
	SizeLapsed          float64   `json:"sizeLapsed"`
	SizeCancelled       float64   `json:"sizeCancelled"`
	SizeVoided          float64   `json:"sizeVoided"`
	RegulatorAuthCode   string    `json:"regulatorAuthCode"`
	RegulatorCode       string    `json:"regulatorCode"`
}

type CurrentOrderSummaryReport struct {
	CurrentOrders []CurrentOrderSumary `json:"currentOrders"`
	MoreAvailable bool                 `json:"moreAvailable"`
}

type ClearedOrderSummary struct {
	BetID           string    `json:"betId"`
	MarketID        string    `json:"marketId"`
	SelectionID     int64     `json:"selectionId"`
	Handicap        float64   `json:"handicap"`
	EventID         int64     `json:"eventId"`
	EventTypeID     int64     `json:"eventTypeId"`
	PlacedDate      time.Time `json:"placedDate"`
	PersistenceType string    `json:"persistenceType"`
	OrderType       string    `json:"orderType"`
	Side            string    `json:"side"`
	ItemDescription string    `json:"itemDescription"` //
	PriceRequested  float64   `json:"priceRequested"`
	SettledDate     time.Time `json:"settledDate"`
	BetCount        int64     `json:"betCount"`
	Commission      float64   `json:"commission"`
	PriceMatched    float64   `json:"priceMatched"`
	PriceReduced    float64   `json:"priceReduced"`
	SizeSettled     float64   `json:"sizeSettled"`
	Profit          float64   `json:"profit"`
	SizeCancelled   float64   `json:"sizeCancelled"`
}

type ClearedOrderSummaryReport struct {
	ClearedOrers  []ClearedOrderSummary `json:"clearedOrders"`
	MoreAvailable bool                  `json:"moreAvailable"`
}
