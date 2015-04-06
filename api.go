package betfair

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Options map[string]interface{}

var (
	defaultOptions = Options{"locale": "en", "exchange": "uk"}

	BettingApiEndpoints = map[string]string{
		"uk": "https://api.betfair.com/exchange/betting/json-rpc/v1",
		"au": "https://api-au.betfair.com/exchange/betting/json-rpc/v1",
	}

	NavigationMenuEndpointFormat = "https://api.betfair.com/exchange/betting/rest/v1/%s/navigation/menu.json"
)

const (
	listEventTypes      = "SportsAPING/v1.0/listEventTypes"
	listCompetitions    = "SportsAPING/v1.0/listCompetitions"
	listEvents          = "SportsAPING/v1.0/listEvents"
	listCountries       = "SportsAPING/v1.0/listCountries"
	listVenues          = "SportsAPING/v1.0/listVenues"
	listMarketCatalogue = "SportsAPING/v1.0/listMarketCatalogue"
	listMarketBook      = "SportsAPING/v1.0/listMarketBook"
	listCurrentOrders   = "SportsAPING/v1.0/listCurrentOrders"
	listClearedOrders   = "SportsAPING/v1.0/listClearedOrders"
)

func (opts1 Options) Merge(opts2 Options) Options {
	var mergedOptions = Options{}

	for k, v := range opts1 {
		mergedOptions[k] = v
	}

	for k, v := range opts2 {
		mergedOptions[k] = v
	}

	return mergedOptions
}

type API struct {
	session *Session
}

func NewAPI(session *Session) *API {
	return &API{session: session}
}

func (api *API) ListEventTypes(options Options) (payload []EventTypeResult, err error) {
	err = api.doRequest(listEventTypes, &payload, extendOptions(Options{"filter": MarketFilter{}}, options))
	return payload, err
}

func (api *API) ListCompetitions(options Options) (result []CompetitionResult, err error) {
	err = api.doRequest(listCompetitions, &result, extendOptions(Options{"filter": MarketFilter{}}, options))
	return result, err
}

func (api *API) ListEvents(options Options) (result []EventResult, err error) {
	err = api.doRequest(listEvents, &result, extendOptions(Options{"filter": MarketFilter{}}, options))
	return result, err
}

func (api *API) ListCountries(options Options) (result []CountryResult, err error) {
	err = api.doRequest(listCountries, &result, extendOptions(Options{"filter": MarketFilter{}}, options))
	return result, err
}

func (api *API) ListVenues(options Options) (result []VenueResult, err error) {
	err = api.doRequest(listVenues, &result, extendOptions(Options{"filter": MarketFilter{}}, options))
	return result, err
}

func (api *API) ListMarketCatalogue(options Options) (result []MarketCatalogue, err error) {
	var catalogueDefaultOptions = Options{
		"filter":           MarketFilter{},
		"marketProjection": []string{"EVENT", "EVENT_TYPE", "COMPETITION"},
		"maxResults":       1000,
	}

	err = api.doRequest(listMarketCatalogue, &result, extendOptions(catalogueDefaultOptions, options))
	return result, err
}

func (api *API) ListMarketBook(marketIds []string, options Options) (result []MarketBook, err error) {
	var marketBookDefaultOptions = Options{
		"marketIds": marketIds,
	}

	err = api.doRequest(listMarketBook, &result, extendOptions(marketBookDefaultOptions, options))
	return result, err
}

func (api *API) ListCurrentOrders(options Options) (result CurrentOrderSummaryReport, err error) {
	var currentOrdersOptions = Options{}
	err = api.doRequest(listCurrentOrders, &result, extendOptions(currentOrdersOptions, options))
	return result, err
}

func (api *API) ListClearedOrders(betStatus string, options Options) (result ClearedOrderSummaryReport, err error) {
	var clearedOrdersOptions = Options{}
	err = api.doRequest(listClearedOrders, &result, extendOptions(clearedOrdersOptions, options))
	return result, err
}

func (api *API) FetchNavigation(options Options) (*Navigation, error) {
	options = extendOptions(Options{}, options)
	locale, _ := options["locale"]
	navigationEndpoint := fmt.Sprintf(NavigationMenuEndpointFormat, locale)
	body, err := api.session.doRawRequest("GET", navigationEndpoint, &strings.Reader{})

	if err != nil {
		return nil, err
	}

	navigation := new(Navigation)
	err = json.Unmarshal(body, navigation)
	return navigation, err
}

func (api *API) PlaceOrders() {
}

func (api *API) CancelOrders() {
}

func (api *API) UpdateOrders() {
}

func (api *API) ReplaceOrders() {
}

func (api *API) buildRequestBody(method string, options Options) ([]byte, error) {
	return json.Marshal(apiRequest{JSONRPC: "2.0", Method: method, Params: options})
}

func (api *API) doRequest(method string, payload interface{}, options Options) error {
	endpoint, err := buildExchangeEndpoint(options)

	if err != nil {
		return err
	}

	body, err := api.buildRequestBody(method, options)

	if err != nil {
		return err
	}

	var response = apiResponse{Result: payload}

	bodyReader := bytes.NewReader(body)

	err = api.session.doRequest(&response, endpoint, bodyReader)

	if err != nil {
		return err
	}

	if response.Error.Code != 0 {
		return fmt.Errorf("API error %d: %s", response.Error.Code, response.Error.Message)
	}

	return nil
}

func buildExchangeEndpoint(options Options) (string, error) {
	exchange := strings.ToLower(fmt.Sprintf("%v", options["exchange"]))
	endpoint, ok := BettingApiEndpoints[exchange]

	if !ok {
		return "", fmt.Errorf("Invalid exchange name `%v`", options["exchange"])
	}

	return endpoint, nil
}

func extendOptions(opts1 Options, opts2 Options) Options {
	var options = Options{}
	var optionsList = []Options{defaultOptions, opts1, opts2}

	for _, opts := range optionsList {
		if opts != nil {
			for k, v := range opts {
				options[k] = v
			}
		}
	}

	return options
}
