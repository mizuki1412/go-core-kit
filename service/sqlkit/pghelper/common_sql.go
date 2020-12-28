package pghelper

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
	"github.com/spf13/cast"
)

func CheckSchemaExist(schema string) bool {
	dao := &sqlkit.Dao{}
	dao.NewHelper("")
	rows := dao.Query(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_namespace WHERE nspname = '%s')", schema))
	defer rows.Close()
	for rows.Next() {
		ret, err := rows.SliceScan()
		if err != nil {
			panic(exception.New(err.Error()))
		}
		return len(ret) > 0 && cast.ToBool(ret[0])
	}
	return false
}
