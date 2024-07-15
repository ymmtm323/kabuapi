package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ymmtm323/kabuapi/irbank"
	"github.com/ymmtm323/kabuapi/yahoof"
)

func main() {
	// echo API
	e := echo.New()

	// CORSの許可
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// ルーティング
	// /search/:id で証券番号の検索
	e.GET("/search/:id", search)
	e.GET("/price/:id", search)

	e.Logger.Fatal(e.Start(":1323"))
}

type SearchResponse struct {
	Number   int     `json:"number"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Dividend float64 `json:"dividend"`
}

func search(c echo.Context) error {
	id := c.Param("id")
	number, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	quoteTypeResponse, err := yahoof.GetQuoteType(id)
	if errors.Is(err, yahoof.ErrorNotFound) {
		return c.JSON(http.StatusNotFound, err)
	} else if errors.Is(err, yahoof.ErrorServiceUnavailable) {
		return c.JSON(http.StatusServiceUnavailable, err)
	} else if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	yahooFinanceResponse, err := yahoof.GetChart(id)
	if errors.Is(err, yahoof.ErrorNotFound) {
		return c.JSON(http.StatusNotFound, err)
	} else if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	dividend, err := irbank.GmtDividend(id)
	if errors.Is(err, irbank.ErrorNotFound) {
		return c.JSON(http.StatusNotFound, err)
	} else if errors.Is(err, irbank.ErrorServiceUnavailable) {
		return c.JSON(http.StatusServiceUnavailable, err)
	} else if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if len(quoteTypeResponse.QuoteType.Result) == 0 || len(yahooFinanceResponse.Chart.Result) == 0 {
		return c.JSON(http.StatusNotFound, yahoof.ErrorNotFound)

	}

	resp := SearchResponse{
		Number:   number,
		Name:     quoteTypeResponse.QuoteType.Result[0].LongName,               // 銘柄名
		Price:    yahooFinanceResponse.Chart.Result[0].Meta.RegularMarketPrice, // 現在の株価
		Dividend: dividend,                                                     // 配当金
	}

	return c.JSON(http.StatusOK, resp)
}

type PriceResponse struct {
	Number int     `json:"number"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
}

func price(c echo.Context) error {
	id := c.Param("id")
	number, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	quoteTypeResponse, err := yahoof.GetQuoteType(id)
	if errors.Is(err, yahoof.ErrorNotFound) {
		return c.JSON(http.StatusNotFound, err)
	} else if errors.Is(err, yahoof.ErrorServiceUnavailable) {
		return c.JSON(http.StatusServiceUnavailable, err)
	} else if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	yahooFinanceResponse, err := yahoof.GetChart(id)
	if errors.Is(err, yahoof.ErrorNotFound) {
		return c.JSON(http.StatusNotFound, err)
	} else if errors.Is(err, yahoof.ErrorServiceUnavailable) {
		return c.JSON(http.StatusServiceUnavailable, err)
	} else if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if len(quoteTypeResponse.QuoteType.Result) == 0 || len(yahooFinanceResponse.Chart.Result) == 0 {
		return c.JSON(http.StatusNotFound, yahoof.ErrorNotFound)
	}

	resp := PriceResponse{
		Number: number,
		Name:   quoteTypeResponse.QuoteType.Result[0].LongName,               // 銘柄名
		Price:  yahooFinanceResponse.Chart.Result[0].Meta.RegularMarketPrice, // 現在の株価
	}

	return c.JSON(http.StatusOK, resp)
}
