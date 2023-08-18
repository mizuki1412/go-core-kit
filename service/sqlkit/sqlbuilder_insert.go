package sqlkit

import "github.com/Masterminds/squirrel"

type InsertBuilder struct {
	inner     squirrel.InsertBuilder
	ModelMeta *ModelMeta
}

func (b InsertBuilder) Sql() (string, []interface{}) {
	return b.inner.MustSql()
}
