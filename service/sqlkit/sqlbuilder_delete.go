package sqlkit

import (
	"github.com/Masterminds/squirrel"
)

type DeleteBuilder struct {
	inner     squirrel.DeleteBuilder
	modelMeta ModelMeta
}

func (b DeleteBuilder) Sql() (string, []interface{}) {
	return b.inner.MustSql()
}
func (b DeleteBuilder) ToSql() (string, []interface{}, error) {
	return b.inner.ToSql()
}

// SQL methods

// Prefix 在 sql 前写入语句
func (b DeleteBuilder) Prefix(sql string, args ...interface{}) DeleteBuilder {
	b.inner = b.inner.Prefix(sql, args...)
	return b
}
func (b DeleteBuilder) PrefixExpr(expr DeleteBuilder) DeleteBuilder {
	b.inner = b.inner.PrefixExpr(expr)
	return b
}
func (b DeleteBuilder) Suffix(sql string, args ...interface{}) DeleteBuilder {
	b.inner = b.inner.Suffix(sql, args...)
	return b
}
func (b DeleteBuilder) SuffixExpr(expr DeleteBuilder) DeleteBuilder {
	b.inner = b.inner.SuffixExpr(expr)
	return b
}

func (b DeleteBuilder) Where(pred interface{}, args ...interface{}) DeleteBuilder {
	b.inner = b.inner.Where(pred, args...)
	return b
}
