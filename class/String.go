package class

import (
	"database/sql"
	"mizuki/project/core-kit/library/jsonkit"
)

// 同时继承scan和value方法
type String struct {
	sql.NullString
}

func (ns String) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return jsonkit.JSON().Marshal(ns.String)
	}
	// 返回json中的null
	return jsonkit.JSON().Marshal(nil)
	//return nil,nil
}
func (ns *String) UnmarshalJSON(data []byte) error {
	var s *string
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}
