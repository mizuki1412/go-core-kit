package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
)

type UpdateDao[T any] struct {
	Dao[T]
	builder squirrel.UpdateBuilder
}

func (dao UpdateDao[T]) Print() {
	sql, args := dao.Sql()
	logkit.Info("sql print", "sql", sql, "args", jsonkit.ToString(args))
}

func (dao UpdateDao[T]) Exec() int64 {
	res := dao.ExecRaw(dao.Sql())
	rn, _ := res.RowsAffected()
	logkit.Debug("sql res", "rows", rn)
	return rn
}

func (dao UpdateDao[T]) Sql() (string, []any) {
	sqls, args, err := dao.ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return sqls, args
}
func (dao UpdateDao[T]) ToSql() (string, []any, error) {
	dao.builder = dao.builder.PlaceholderFormat(placeholder(dao.dataSource.Driver))
	return dao.builder.ToSql()
}

// SQL methods

// Prefix 在 sql 前写入语句
func (dao UpdateDao[T]) Prefix(sql string, args ...interface{}) UpdateDao[T] {
	dao.builder = dao.builder.Prefix(sql, args...)
	return dao
}
func (dao UpdateDao[T]) PrefixExpr(expr squirrel.Sqlizer) UpdateDao[T] {
	dao.builder = dao.builder.PrefixExpr(expr)
	return dao
}
func (dao UpdateDao[T]) Suffix(sql string, args ...interface{}) UpdateDao[T] {
	dao.builder = dao.builder.Suffix(sql, args...)
	return dao
}
func (dao UpdateDao[T]) SuffixExpr(expr squirrel.Sqlizer) UpdateDao[T] {
	dao.builder = dao.builder.SuffixExpr(expr)
	return dao
}

func (dao UpdateDao[T]) Set(column string, value interface{}) UpdateDao[T] {
	dao.builder = dao.builder.Set(dao.modelMeta.escapeName(column), value)
	return dao
}
func (dao UpdateDao[T]) Where(pred interface{}, args ...interface{}) UpdateDao[T] {
	dao.builder = dao.builder.Where(pred, args...)
	return dao
}
func (dao UpdateDao[T]) FromSelect(from SelectDao[T], alias string) UpdateDao[T] {
	dao.builder = dao.builder.FromSelect(from.builder, alias)
	return dao
}
