package main

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"
)

var API_URL string = "https://api.coingecko.com/api/v3/coins/markets?"


func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check the file extension and set the appropriate MIME type
		switch path.Ext(r.URL.Path) {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".html":
			w.Header().Set("Content-Type", "text/html")
		}

		http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
	})

	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.ListenAndServe("localhost:8080", nil)
}


func queryCoinMarket(coinIds []string, currency string) ([]CoinMarketData, error) {
	queryArgs := "vs_currency=" + currency + "&ids=" + strings.Join(coinIds, ",")

	resp, err := http.Get(API_URL + queryArgs)
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