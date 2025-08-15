package dateStringToTime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func DateStringToTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("data vazia")
	}

	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}

	// Epoch (segundos ou milissegundos).
	if isDigits(s) {
		if unix, err := strconv.ParseInt(s, 10, 64); err == nil {
			// Heurística simples: >= 13 dígitos => ms
			if len(s) >= 13 {
				return time.UnixMilli(unix).In(time.Local), nil
			}
			return time.Unix(unix, 0).In(time.Local), nil
		}
	}

	return time.Time{}, fmt.Errorf("formato de data inválido: %q (esperado: DATE/DATETIME/TIMESTAMP/RFC3339/epoch)", s)
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}
