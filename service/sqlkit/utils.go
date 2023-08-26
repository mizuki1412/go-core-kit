package sqlkit

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
)

func logReqSqlInfo(sql string, args []any) string {
	return fmt.Sprintf(`
==> %s
==> %s`, sql, jsonkit.ToString(args))
}

func logResSqlInfo(rows int64) string {
	return fmt.Sprintf(`
<== rows: %d`, rows)
}

func scanObjList[T any](dao SelectDao[T]) []*T {
	rows := dao.QueryRows()
	list := make([]*T, 0, 5)
	defer rows.Close()
	for rows.Next() {
		m := new(T)
		err := rows.StructScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, m)
	}
	if dao.Cascade != nil {
		for i := range list {
			dao.Cascade(list[i])
		}
	}
	return list
}
