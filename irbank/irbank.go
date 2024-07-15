package irbank

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ErrorNotFound           = errors.New("not found")
	ErrorServiceUnavailable = errors.New("service unavailable")
)

func GmtDividend(code string) (float64, error) {
	resp, err := http.Get("https://f.irbank.net/files/" + code + "/fy-stock-dividend.csv")
	if err != nil {
		return 0, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return 0, ErrorNotFound
	} else if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusServiceUnavailable {
		return 0, ErrorServiceUnavailable
	} else if resp.StatusCode != http.StatusOK {
		return 0, errors.New("status code is not 200")
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(buf), "\n")

	// 予想がある場合は最新の予想を返す
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], "（予想）") {
			s := strings.Split(lines[i], ",")
			if len(s) < 2 {
				return 0, errors.New("invalid format")
			}
			dividend, err := strconv.ParseFloat(s[1], 64)
			if err != nil {
				return 0, err
			}
			return dividend, nil
		}
	}

	// 予想がなかった場合は直近の実績を返す
	// 現在の年度を探す
	year := getFiscalYear()
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], fmt.Sprintf("%s/04", year)) {
			s := strings.Split(lines[i], ",")
			if len(s) < 2 {
				return 0, errors.New("invalid format")
			}
			dividend, err := strconv.ParseFloat(s[1], 64)
			if err != nil {
				return 0, err
			}
			return dividend, nil
		}
	}

	return 0, errors.New("not found")
}

// 年度を取得する
func getFiscalYear() string {
	// 現在の年と月を取得
	year := time.Now().Year()
	month := int(time.Now().Month())

	if month >= 4 {
		return strconv.Itoa(year)
	}
	return strconv.Itoa(year - 1)
}
