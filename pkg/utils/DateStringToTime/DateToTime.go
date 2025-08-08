package DateStringToTime

import (
	"errors"
	"time"
)

func DateStringToTime(dateString string) (time.Time, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}, errors.New("formato de data inválido, esperado: YYYY-MM-DD")
	}
	return t, nil
}
