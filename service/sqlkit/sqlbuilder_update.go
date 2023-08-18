package sqlkit

import "github.com/Masterminds/squirrel"

type UpdateBuilder struct {
	inner     squirrel.UpdateBuilder
	ModelMeta *ModelMeta
}

func (b UpdateBuilder) Sql() (string, []interface{}) {
	return b.inner.MustSql()
}
