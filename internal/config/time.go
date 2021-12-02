package config

import (
	"fmt"
	"strings"
	"time"
)

// LayoutDateFormat - формат вывода даты по умолчанию для всего приложения
const LayoutDateFormat = "2006-01-02T15:04:05-07:00"

// LayoutDateFormatLogFile - формат даты для ведения файловых логов
const LayoutDateFormatLogFile = "2006-01-02T15:04:05.999"

// DateTime - Забиваем на все и делаем по своему от и до
type DateTime struct {
	time.Time
}

// UnmarshalJSON - десериализация даты
func (dt *DateTime) UnmarshalJSON(b []byte) (err error) {

	str := strings.Trim(string(b), `"`) // remove quotes
	if str == "" {
		return
	}

	dt.Time, err = time.Parse(LayoutDateFormat, str)

	return
}

// MarshalJSON - Сериализация даты
func (dt DateTime) MarshalJSON() ([]byte, error) {

	if dt.IsZero() {
		return []byte(fmt.Sprintf(`"%s"`, "")), nil
	}

	return []byte(fmt.Sprintf(`"%s"`, dt.Time.Format(LayoutDateFormat))), nil
}

// NewCurrentTime - new current time
func NewCurrentTime() DateTime {
	return DateTime{
		Time: time.Now().UTC(),
	}
}
