package dateStringToTime

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func DateStringToTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("data vazia")
	}

	// 1) Tenta RFC3339/ISO com fuso
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	// 2) Aceita YYYY-MM-DD (sem hora/TZ) -> meia-noite na timezone local do servidor
	if len(s) == 10 && s[4] == '-' && s[7] == '-' {
		// se quiser cravar o fuso da app:
		// loc, _ := time.LoadLocation("America/Campo_Grande")
		// return time.ParseInLocation("2006-01-02", s, loc)
		return time.ParseInLocation("2006-01-02", s, time.Local)
	}

	// 3) Demais formatos sem TZ já usados no projeto
	for _, layout := range []string{
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05",
	} {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported date format: %q", s)
}
