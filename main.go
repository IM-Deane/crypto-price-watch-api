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
	http.HandleFunc("/coins", displayCoinsList)

	http.ListenAndServe("localhost:8080", nil)
}

// *** Handlers ***

func displayHome(w http.ResponseWriter, r *http.Request) {

	// TODO: uncomment this once the API is working
	// coins, err := queryCoinList()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	coins := []CryptoCoin{
		{Id: "bitcoin", Symbol: "btc", Name: "Bitcoin"},
		{Id: "ethereum", Symbol: "eth", Name: "Ethereum"},
	}
	

	tmpl, err := template.ParseFiles("./static/index.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		log.Default().Println(err)
		return
	}

	tmpl.Execute(w, coins)
}

func displayCoinsList(w http.ResponseWriter, r *http.Request) {
    coinIdsStr := r.URL.Query().Get("ids")
    if coinIdsStr == "" {
        http.Error(w, "No coins specified", http.StatusBadRequest)
        return
    }

    coinIds := strings.Split(coinIdsStr, ",")

    // Fetch coin details
    coins, err := queryCoinMarket(coinIds, "cad") 
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if len(coins) == 0 {
        http.Error(w, "No coins found", http.StatusNotFound)
        return
    }

    // Render the coins.html template with the coins data
    tmpl, err := template.New("coins.html").Funcs(template.FuncMap{
		"colorClass": colorClass,
	}).ParseFiles("./static/coins.html")
    if err != nil {
        http.Error(w, "Failed to load template", http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, coins)
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


// Helpers
func colorClass(v float64) string {
	if v > 0 {
		return "text-green-500"
	} else if (v == 0) {
		return "text-gray-500"
	} else {
		return "text-red-500"
	}
}