package irbank

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrorNotFound = errors.New("not found")
)

func GmtDividend(code string) (float64, error) {
	resp, err := http.Get("https://f.irbank.net/files/" + code + "/fy-stock-dividend.csv")
	if err != nil {
		return 0, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return 0, ErrorNotFound
	} else if resp.StatusCode != http.StatusOK {
		return 0, errors.New("status code is not 200")
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(buf), "\n")
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

	return 0, errors.New("not found")
}
