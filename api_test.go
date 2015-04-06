package betfair

import (
	"testing"
	"time"
)

func getTestAPI() *API {
	return NewAPI(GetTestSession())
}

func TestEventTypes(t *testing.T) {
	var api = getTestAPI()
	eventTypes, err := api.ListEventTypes(Options{"filter": MarketFilter{EventTypeIDs: []string{"1"}}})

	if err != nil {
		t.Error(err)
		return
	}

	if len(eventTypes) == 0 {
		t.Error("Could not get any event type")
		return
	}
}

func TestCompetitions(t *testing.T) {
	var api = getTestAPI()

	competitions, err := api.ListCompetitions(Options{})

	if err != nil {
		t.Error(err)
		return
	}

	if len(competitions) == 0 {
		t.Error("Could not get any competition")
		return
	}
}

func TestEvents(t *testing.T) {
	var api = getTestAPI()

	from := time.Now()
	to := from.Add(time.Hour * 24)

	events, err := api.ListEvents(Options{"exchange": "au", "filter": MarketFilter{MarketStartTime: &TimeRange{From: from, To: to}}})

	if err != nil {
		t.Error(err)
		return
	}

	if len(events) == 0 {
		t.Error("Could not get any event")
		return
	}
}

func TestCountries(t *testing.T) {
	var api = getTestAPI()

	countries, err := api.ListCountries(Options{})

	if err != nil {
		t.Error(err)
		return
	}

	if len(countries) == 0 {
		t.Error("Could not get any country")
		return
	}
}

func TestVenues(t *testing.T) {
	var api = getTestAPI()

	venues, err := api.ListVenues(Options{})

	if err != nil {
		t.Error(err)
		return
	}

	if len(venues) == 0 {
		t.Error("Could not get any venue")
		return
	}
}

func TestMarketCatalogue(t *testing.T) {
	var api = getTestAPI()

	markets, err := api.ListMarketCatalogue(Options{"maxResults": 1})

	if err != nil {
		t.Error(err)
		return
	}

	if len(markets) == 0 {
		t.Error("Could not get any market")
		return
	}
}

func TestMarketBook(t *testing.T) {
	var api = getTestAPI()

	markets, err := api.ListMarketCatalogue(Options{"maxResults": 1})

	if err != nil {
		t.Error(err)
		return
	}

	var market = markets[0]

	outcomes, err := api.ListMarketBook([]string{market.MarketID}, Options{})

	if err != nil {
		t.Error(err)
		return
	}

	if len(outcomes) == 0 {
		t.Error("Could not get market book")
		return
	}
}

func TestListCurrentOrders(t *testing.T) {
	var api = getTestAPI()

	_, err := api.ListCurrentOrders(Options{"orderProjection": "ALL"})

	if err != nil {
		t.Error(err)
		return
	}
}

func TestFetchNavigation(t *testing.T) {
	var api = getTestAPI()

	_, err := api.FetchNavigation(Options{})

	if err != nil {
		t.Error(err)
		return
	}
}

func BenchmarkEventTypes(t *testing.B) {
	var api = getTestAPI()
	eventTypes, err := api.ListEventTypes(Options{"filter": MarketFilter{}})

	if err != nil {
		t.Error(err)
		return
	}

	if len(eventTypes) == 0 {
		t.Error("Could not get any event type")
		return
	}
}
