package sqlkit

import "github.com/Masterminds/squirrel"

type DeleteBuilder struct {
	inner     squirrel.DeleteBuilder
	ModelMeta *ModelMeta
}

func (b DeleteBuilder) Sql() (string, []interface{}) {
	return b.inner.MustSql()
}
