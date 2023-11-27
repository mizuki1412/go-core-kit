package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
)

type DeleteDao[T any] struct {
	Dao[T]
	builder squirrel.DeleteBuilder
}

func (dao DeleteDao[T]) Print() {
	sql, args := dao.Sql()
	logkit.Info("sql print", "sql", sql, "args", jsonkit.ToString(args))
}

func (dao DeleteDao[T]) Exec() int64 {
	res := dao.ExecRaw(dao.Sql())
	rn, _ := res.RowsAffected()
	logkit.Debug("sql res", "rows", rn)
	return rn
}

func (dao DeleteDao[T]) Sql() (string, []any) {
	sqls, args, err := dao.ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return sqls, args
}
func (dao DeleteDao[T]) ToSql() (string, []any, error) {
	dao.builder = dao.builder.PlaceholderFormat(placeholder(dao.dataSource.Driver))
	return dao.builder.ToSql()
}

// SQL methods

// Prefix 在 sql 前写入语句
func (dao DeleteDao[T]) Prefix(sql string, args ...any) DeleteDao[T] {
	dao.builder = dao.builder.Prefix(sql, args...)
	return dao
}
func (dao DeleteDao[T]) PrefixExpr(expr squirrel.Sqlizer) DeleteDao[T] {
	dao.builder = dao.builder.PrefixExpr(expr)
	return dao
}
func (dao DeleteDao[T]) Suffix(sql string, args ...any) DeleteDao[T] {
	dao.builder = dao.builder.Suffix(sql, args...)
	return dao
}
func (dao DeleteDao[T]) SuffixExpr(expr squirrel.Sqlizer) DeleteDao[T] {
	dao.builder = dao.builder.SuffixExpr(expr)
	return dao
}

func (dao DeleteDao[T]) Where(pred any, args ...any) DeleteDao[T] {
	dao.builder = dao.builder.Where(pred, args...)
	return dao
}
