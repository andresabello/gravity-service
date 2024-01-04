package renderutil

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"pi-search/pkg/tracer"
)

type Rendered struct {
	Rendered sql.NullString `json:"rendered"`
}

func (r *Rendered) UnmarshalJSON(data []byte) error {
	var temp struct {
		Rendered string `json:"rendered"`
	}

	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	if temp.Rendered == "" {
		r.Rendered = sql.NullString{String: "", Valid: false}
		return nil
	}

	r.Rendered = sql.NullString{String: temp.Rendered, Valid: true}

	return nil
}

// Scan assigns a value from a database driver. Scanner interface.
func (r *Rendered) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		r.Rendered = sql.NullString{}
	case string:
		r.Rendered = sql.NullString{String: v, Valid: true}
	default:
		return tracer.TraceError(
			fmt.Errorf("unsupported type for Rendered: %T", value),
		)
	}
	return nil
}

// Value return json value, implement driver.Valuer interface.
func (r *Rendered) Value() (driver.Value, error) {
	// if len(r) == 0 {
	return nil, nil
	// }
	// return json.RawMessage(r).MarshalJSON()
}
