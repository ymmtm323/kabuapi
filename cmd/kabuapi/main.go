package main

import (
	"errors"
	"fmt"
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
	e.GET("/price/:id", price)

	e.Logger.Fatal(e.Start(":1323"))
}

type SearchResponse struct {
	Number   int     `json:"number"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Dividend float64 `json:"dividend"`
	Error    string  `json:"error"`
}

func search(c echo.Context) error {
	var resp SearchResponse
	id := c.Param("id")
	number, err := strconv.Atoi(id)
	if err != nil {
		resp.Error = fmt.Sprintf("invalid id: %s", id)
		return c.JSON(http.StatusBadRequest, resp)
	}
	resp.Number = number

	quoteTypeResponse, err := yahoof.GetQuoteType(id)
	if errors.Is(err, yahoof.ErrorNotFound) {
		resp.Error = fmt.Sprintf("not found: %s", id)
		return c.JSON(http.StatusNotFound, resp)
	} else if errors.Is(err, yahoof.ErrorServiceUnavailable) {
		resp.Error = fmt.Sprintf("service unavailable: %s", id)
		return c.JSON(http.StatusServiceUnavailable, resp)
	} else if err != nil {
		resp.Error = fmt.Sprintf("internal server error: %s", id)
		return c.JSON(http.StatusInternalServerError, resp)
	}
	if len(quoteTypeResponse.QuoteType.Result) == 0 {
		resp.Error = fmt.Sprintf("not found: %s", id)
		return c.JSON(http.StatusNotFound, resp)
	}
	resp.Name = quoteTypeResponse.QuoteType.Result[0].LongName

	yahooFinanceResponse, err := yahoof.GetChart(id)
	if errors.Is(err, yahoof.ErrorNotFound) {
		resp.Error = fmt.Sprintf("not found: %s", id)
		return c.JSON(http.StatusNotFound, resp)
	} else if err != nil {
		resp.Error = fmt.Sprintf("internal server error: %s", id)
		return c.JSON(http.StatusInternalServerError, resp)
	}
	if len(yahooFinanceResponse.Chart.Result) == 0 {
		resp.Error = fmt.Sprintf("not found: %s", id)
		return c.JSON(http.StatusNotFound, resp)
	}
	resp.Price = yahooFinanceResponse.Chart.Result[0].Meta.RegularMarketPrice

	dividend, err := irbank.GmtDividend(id)
	if errors.Is(err, irbank.ErrorNotFound) {
		resp.Error = fmt.Sprintf("not found: %s", id)
		return c.JSON(http.StatusNotFound, resp)
	} else if errors.Is(err, irbank.ErrorServiceUnavailable) {
		resp.Error = fmt.Sprintf("service unavailable: %s", id)
		return c.JSON(http.StatusServiceUnavailable, resp)
	} else if err != nil {
		resp.Error = fmt.Sprintf("internal server error: %s", id)
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Dividend = dividend

	return c.JSON(http.StatusOK, resp)
}

type PriceResponse struct {
	Number int     `json:"number"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Error  string  `json:"error"`
}

func price(c echo.Context) error {
	var resp PriceResponse
	id := c.Param("id")
	number, err := strconv.Atoi(id)
	if err != nil {
		resp.Error = fmt.Sprintf("invalid id: %s", id)
		return c.JSON(http.StatusBadRequest, resp)
	}
	resp.Number = number

	quoteTypeResponse, err := yahoof.GetQuoteType(id)
	if errors.Is(err, yahoof.ErrorNotFound) {
		resp.Error = fmt.Sprintf("not found: %s", id)
		return c.JSON(http.StatusNotFound, resp)
	} else if errors.Is(err, yahoof.ErrorServiceUnavailable) {
		resp.Error = fmt.Sprintf("service unavailable: %s", id)
		return c.JSON(http.StatusServiceUnavailable, resp)
	} else if err != nil {
		resp.Error = fmt.Sprintf("internal server error: %s", id)
		return c.JSON(http.StatusInternalServerError, resp)
	}
	if len(quoteTypeResponse.QuoteType.Result) == 0 {
		resp.Error = fmt.Sprintf("not found: %s", id)
		return c.JSON(http.StatusNotFound, resp)
	}
	resp.Name = quoteTypeResponse.QuoteType.Result[0].LongName

	yahooFinanceResponse, err := yahoof.GetChart(id)
	if errors.Is(err, yahoof.ErrorNotFound) {
		resp.Error = fmt.Sprintf("not found: %s", id)
		return c.JSON(http.StatusNotFound, resp)
	} else if errors.Is(err, yahoof.ErrorServiceUnavailable) {
		resp.Error = fmt.Sprintf("service unavailable: %s", id)
		return c.JSON(http.StatusServiceUnavailable, resp)
	} else if err != nil {
		resp.Error = fmt.Sprintf("internal server error: %s", id)
		return c.JSON(http.StatusInternalServerError, resp)
	}
	if len(yahooFinanceResponse.Chart.Result) == 0 {
		resp.Error = fmt.Sprintf("not found: %s", id)
		return c.JSON(http.StatusNotFound, resp)
	}
	resp.Price = yahooFinanceResponse.Chart.Result[0].Meta.RegularMarketPrice

	return c.JSON(http.StatusOK, resp)
}
