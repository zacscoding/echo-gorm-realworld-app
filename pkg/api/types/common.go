package types

import (
	"fmt"
	"time"
)

const (
	timestampFmt = "2006-01-02T15:04:05.000Z"
)

type JSONTime time.Time

func (j JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(j).Format(timestampFmt))
	return []byte(stamp), nil
}
