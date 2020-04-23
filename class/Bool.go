package class

import (
	"database/sql"
	"mizuki/project/core-kit/library/jsonkit"
)

// 同时继承scan和value方法
type Bool struct {
	sql.NullBool
}

func (th Bool) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Bool)
	}
	return jsonkit.JSON().Marshal(nil)
}
func (th *Bool) UnmarshalJSON(data []byte) error {
	var s *bool
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		th.Valid = true
		th.Bool = *s
	} else {
		th.Valid = false
	}
	return nil
}
