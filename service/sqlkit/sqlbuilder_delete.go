package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
)

type DeleteBuilder struct {
	inner     squirrel.DeleteBuilder
	modelMeta ModelMeta
	driver    string
}

func (b DeleteBuilder) Sql() (string, []interface{}) {
	sql, args, err := b.ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return sql, args
}
func (b DeleteBuilder) ToSql() (string, []interface{}, error) {
	b.inner = b.inner.PlaceholderFormat(placeholder(b.driver))
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
