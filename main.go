package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"text/template"
)

var API_URL_COIN_MARKET string = "https://api.coingecko.com/api/v3/coins/markets?"
var API_URL_COIN_LIST string = "https://api.coingecko.com/api/v3/coins/list?include_platform=false"

type CryptoCoin struct {
	Id                string  `json:"id"`
	Symbol            string  `json:"symbol"`
	Name              string  `json:"name"`
}

type CoinMarketData struct {
	Id                string  `json:"id"`
	Symbol            string  `json:"symbol"`
	Name              string  `json:"name"`
	Image             string  `json:"image"`
	CurrentPrice      float64 `json:"current_price"`
	MarketCap         float64 `json:"market_cap"`
	MarketCapRank     int     `json:"market_cap_rank"`
	FullyDilutedValue float64 `json:"fully_diluted_value"`
	TotalVolume       float64 `json:"total_volume"`
	High24            float64 `json:"high_24h"`
	Low24             float64 `json:"low_24h"`
	PriceChange24     float64 `json:"price_change_24h"`
	PriceChangePercentage24 float64 `json:"price_change_percentage_24h"`
	MarketCapChange24 float64 `json:"market_cap_change_24h"`
	MarketCapChangePercentage24 float64 `json:"market_cap_change_percentage_24h"`
	CirculatingSupply float64 `json:"circulating_supply"`
	TotalSupply       float64 `json:"total_supply"`
	MaxSupply         float64 `json:"max_supply"`
	ATH               float64 `json:"ath"`
	ATHChangePercentage float64 `json:"ath_change_percentage"`
	ATHDate           string  `json:"ath_date"`
	ATL               float64 `json:"atl"`
	ATLChangePercentage float64 `json:"atl_change_percentage"`
	ATLDate           string  `json:"atl_date"`
	ROI               float64 `json:"roi"`
	LastUpdated       string  `json:"last_updated"`
}



func main() {

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", displayHome)
	http.HandleFunc("/api/search", displayCoinSearchResults)

	http.ListenAndServe("localhost:8080", nil)
}

// *** Handlers ***

func displayHome(w http.ResponseWriter, r *http.Request) {
	// switch path.Ext(r.URL.Path) {
	// case ".css":
	// 	w.Header().Set("Content-Type", "text/css")
	// case ".js":
	// 	w.Header().Set("Content-Type", "application/javascript")
	// case ".html":
	// 	w.Header().Set("Content-Type", "text/html")
	// }

	// coins, err := queryCoinList()
	coins := []CryptoCoin{
		{Id: "bitcoin", Symbol: "btc", Name: "Bitcoin"},
		{Id: "ethereum", Symbol: "eth", Name: "Ethereum"},
	}
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	tmpl, err := template.ParseFiles("./static/index.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		log.Default().Println(err)
		return
	}

	tmpl.Execute(w, coins)
}

func displayCoinSearchResults(w http.ResponseWriter, r *http.Request) {
	var defaultCurrency string = "cad"

	queryParams := r.URL.Query()
	if len(queryParams) == 0 {
		http.Error(w, "Missing query parameters: ids", http.StatusBadRequest)
		return
	}

	ids := queryParams["ids"]
	if len(ids) == 0 {
		http.Error(w, "Missing query parameter 'ids'", http.StatusBadRequest)
		return
	}
	currency := queryParams.Get("currency")
	if len(currency) == 0 {
		currency = defaultCurrency
	}


	coin, err := queryCoinMarket(ids, currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(coin)

	tmpl, err := template.ParseFiles("./static/coins.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, coin)
}

func displayCoinList(w http.ResponseWriter, r *http.Request) {
	coins, err := queryCoinList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(coins)
}


// *** Queries ***

func queryCoinMarket(coinIds []string, currency string) ([]CoinMarketData, error) {
	queryArgs := "vs_currency=" + currency + "&ids=" + strings.Join(coinIds, ",")

	resp, err := http.Get(API_URL_COIN_MARKET + queryArgs)
	if err != nil {
		return []CoinMarketData{}, err
	}

	defer resp.Body.Close()
	
	var coins []CoinMarketData = []CoinMarketData{}

	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		return []CoinMarketData{}, err
	}

	return coins, nil
}

func queryCoinList() ([]CryptoCoin, error) {
	resp, err := http.Get(API_URL_COIN_LIST)
	if err != nil {
		return []CryptoCoin{}, err
	}

	defer resp.Body.Close()

	var coins []CryptoCoin = []CryptoCoin{}

	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		return []CryptoCoin{}, err
	}

	return coins, nil
}


