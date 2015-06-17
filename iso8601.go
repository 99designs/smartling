package smartling

import (
	"net/url"
	"time"
)

const Format = "2006-01-02T15:04:05"
const jsonFormat = `"` + Format + `"`

type Iso8601Time time.Time

func (it Iso8601Time) EncodeValues(key string, v *url.Values) error {
	v.Add(key, time.Time(it).Format(Format))
	return nil
}

func (it Iso8601Time) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(it).Format(jsonFormat)), nil
}

func (it *Iso8601Time) UnmarshalJSON(data []byte) error {
	t, err := time.ParseInLocation(jsonFormat, string(data), time.FixedZone("", 0))
	if err == nil {
		*it = Iso8601Time(t)
	}

	return err
}
