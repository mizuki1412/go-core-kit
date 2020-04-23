package class

import (
	"database/sql"
	"mizuki/project/core-kit/library/jsonkit"
)

// 同时继承scan和value方法
type Int64 struct {
	sql.NullInt64
}

func (th Int64) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Int64)
	}
	return jsonkit.JSON().Marshal(nil)
}
func (th *Int64) UnmarshalJSON(data []byte) error {
	var s *int64
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		th.Valid = true
		th.Int64 = *s
	} else {
		th.Valid = false
	}
	return nil
}
