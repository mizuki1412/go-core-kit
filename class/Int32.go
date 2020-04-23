package class

import (
	"database/sql"
	"mizuki/project/core-kit/library/jsonkit"
)

// 同时继承scan和value方法
type Int32 struct {
	sql.NullInt32
}

func (th Int32) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Int32)
	}
	return jsonkit.JSON().Marshal(nil)
}
func (th *Int32) UnmarshalJSON(data []byte) error {
	var s *int32
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		th.Valid = true
		th.Int32 = *s
	} else {
		th.Valid = false
	}
	return nil
}
