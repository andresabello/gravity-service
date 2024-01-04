package timeutil

import (
	"fmt"
	"pi-search/pkg/tracer"
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		ct.Time = v
	case []byte:
		t, err := time.Parse(time.RFC3339, string(v))
		if err != nil {
			return err
		}
		ct.Time = t
	default:
		return tracer.TraceError(
			fmt.Errorf("unsupported type for CustomTime: %T", value),
		)
	}

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	// Remove surrounding quotes from the date string
	s = s[1 : len(s)-1]

	// Define the layout for your date format
	layout := "2006-01-02T15:04:05"
	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}
