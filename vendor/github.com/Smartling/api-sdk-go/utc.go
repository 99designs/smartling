package smartling

import (
	"encoding/json"
	"time"
)

const utcFormat = "2006-01-02T15:04:05Z"

// UTC represents time in UTC format (zero timezone).
type UTC struct {
	time.Time
}

// MarshalJSON returns JSON representation of UTC.
func (utc UTC) MarshalJSON() ([]byte, error) {
	return json.Marshal(utc.String())
}

// UnmarshalJSON parses JSON representation of UTC.
func (utc *UTC) UnmarshalJSON(data []byte) error {
	var formatted string

	err := json.Unmarshal(data, &formatted)
	if err != nil {
		return err
	}

	location, err := time.LoadLocation("UTC")
	if err != nil {
		return err
	}

	parsed, err := time.ParseInLocation(utcFormat, formatted, location)
	if err != nil {
		return err
	}

	*utc = UTC{parsed}

	return nil
}

// String returns string reprenation of UTC.
func (utc UTC) String() string {
	return utc.Format(utcFormat)
}
