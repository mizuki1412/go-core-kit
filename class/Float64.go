package class

import (
	"database/sql"
	"mizuki/project/core-kit/library/jsonkit"
)

// 同时继承scan和value方法
type Float64 struct {
	sql.NullFloat64
}

func (th Float64) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Float64)
	}
	return jsonkit.JSON().Marshal(nil)
}
func (th *Float64) UnmarshalJSON(data []byte) error {
	var s *float64
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		th.Valid = true
		th.Float64 = *s
	} else {
		th.Valid = false
	}
	return nil
}
