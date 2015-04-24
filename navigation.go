package betfair

import (
	"encoding/json"
	"time"
)

type NavigationChild struct {
	Type            string          `json:"type"`
	Name            string          `json:"name"`
	ID              string          `json:"id"`
	CountryCode     string          `json:"countryCode"`
	ExchangeID      string          `json:"exchangeId"`
	MarketType      string          `json:"marketType"`
	MarketStartTime time.Time       `json:"marketStartTime"`
	NumberOfWinners json.RawMessage `json:"numberOfWinners"`
	Children        []NavigationChild
}

type Navigation struct {
	Children []NavigationChild `json:"children"`
}
