package sqlkit

import (
	"github.com/Masterminds/squirrel"
)

func Builder() squirrel.StatementBuilderType {
	if driverName() == "postgres" {
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	} else {
		// todo 未处理oracle
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	}
}
