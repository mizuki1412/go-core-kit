package sqlkit

import (
	"github.com/Masterminds/squirrel"
)

type SelectBuilder struct {
	inner     squirrel.SelectBuilder
	ModelMeta *ModelMeta
}

func (b SelectBuilder) Sql() (string, []interface{}) {
	b.inner.
	return b.inner.MustSql()
}


