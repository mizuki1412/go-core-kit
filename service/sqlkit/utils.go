package sqlkit

import (
	"fmt"
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
