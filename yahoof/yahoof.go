package yahoof

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (
	ErrorNotFound = errors.New("not found")
)

type QuoteTypeResponse struct {
	QuoteType struct {
		Result []struct {
			Symbol                string `json:"symbol"`
			QuoteType             string `json:"quoteType"`
			Exchange              string `json:"exchange"`
			ShortName             string `json:"shortName"`
			LongName              string `json:"longName"`
			MessageBoardId        string `json:"messageBoardId"`
			ExchangeTimezoneName  string `json:"exchangeTimezoneName"`
			ExchangeTimezoneShort string `json:"exchangeTimezoneShortName"`
			GmtOffSetMilliseconds string `json:"gmtOffSetMilliseconds"`
			Market                string `json:"market"`
			IsEsgPopulated        bool   `json:"isEsgPopulated"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"quoteType"`
}

func GetQuoteType(id string) (QuoteTypeResponse, error) {
	// https://query2.finance.yahoo.com/v1/finance/quoteType/?symbol=${id}.T&lang=ja-JP&region=JP
	resp, err := http.Get("https://query2.finance.yahoo.com/v1/finance/quoteType/?symbol=" + id + ".T&lang=ja-JP&region=JP")
	if err != nil {
		return QuoteTypeResponse{}, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return QuoteTypeResponse{}, ErrorNotFound
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return QuoteTypeResponse{}, err
	}
	quoteTypeResponse := QuoteTypeResponse{}
	if err := json.Unmarshal(buf, &quoteTypeResponse); err != nil {
		return QuoteTypeResponse{}, err
	}
	return quoteTypeResponse, nil
}

type ChartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency             string  `json:"currency"`
				Symbol               string  `json:"symbol"`
				ExchangeName         string  `json:"exchangeName"`
				FullExchangeName     string  `json:"fullExchangeName"`
				InstrumentType       string  `json:"instrumentType"`
				FirstTradeDate       int     `json:"firstTradeDate"`
				RegularMarketTime    int     `json:"regularMarketTime"`
				HasPrePostMarketData bool    `json:"hasPrePostMarketData"`
				Gmtoffset            int     `json:"gmtoffset"`
				Timezone             string  `json:"timezone"`
				ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
				RegularMarketPrice   float64 `json:"regularMarketPrice"`
				FiftyTwoWeekHigh     float64 `json:"fiftyTwoWeekHigh"`
				FiftyTwoWeekLow      float64 `json:"fiftyTwoWeekLow"`
				RegularMarketDayHigh float64 `json:"regularMarketDayHigh"`
				RegularMarketDayLow  float64 `json:"regularMarketDayLow"`
				RegularMarketVolume  float64 `json:"regularMarketVolume"`
				ChartPreviousClose   float64 `json:"chartPreviousClose"`
				PriceHint            float64 `json:"priceHint"`
				CurrentTradingPeriod struct {
					Pre struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"per"`
					Regular struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"regular"`
					Post struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"post"`
				} `json:"currentTradingPeriod"`
				DataGranularity string   `json:"dataGranularity"`
				Range           string   `json:"range"`
				ValidRanges     []string `json:"validRanges"`
			} `json:"meta"`
			Timestamp  []int `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Volume []float64 `json:"volume"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
				} `json:"quote"`
				Adjclose []struct {
					Adjclose []float64 `json:"adjclose"`
				} `json:"adjclose"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

func GetChart(id string) (ChartResponse, error) {
	// https://query1.finance.yahoo.com/v8/finance/chart/${id}.T?interval=1d
	resp, err := http.Get("https://query1.finance.yahoo.com/v8/finance/chart/" + id + ".T?interval=1d")
	if err != nil {
		return ChartResponse{}, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return ChartResponse{}, ErrorNotFound
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChartResponse{}, err
	}
	chartResponse := ChartResponse{}
	if err := json.Unmarshal(buf, &chartResponse); err != nil {
		return ChartResponse{}, err
	}
	return chartResponse, nil
}
