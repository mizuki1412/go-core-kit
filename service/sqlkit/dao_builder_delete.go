package sqlkit

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/v2/class/const/sqlconst"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
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
	sqls, args, err := dao.builder.ToSql()
	return sqls, argsWrap(dao.dataSource.Driver, args), err
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

// custom 参考dao_builder_select

func (dao DeleteDao[T]) whereUnnest(arr any, key, flag string) DeleteDao[T] {
	switch dao.dataSource.Driver {
	case sqlconst.Postgres:
		s, v := pgArray(arr)
		return dao.Where(fmt.Sprintf("%s %s (select unnest(%s))", dao.modelMeta.escapeName(key), flag, s), v...)
	default:
		panic(exception.New("whereUnnest not supported"))
	}
}
func (dao DeleteDao[T]) WhereUnnestIn(key string, arr any) DeleteDao[T] {
	return dao.whereUnnest(arr, key, "IN")
}
func (dao DeleteDao[T]) WhereUnnestNotIn(key string, arr any) DeleteDao[T] {
	return dao.whereUnnest(arr, key, "NOT IN")
}

// WhereArrayIn 用于PG中array类型数据的包含比较
func (dao DeleteDao[T]) WhereArrayIn(key string, arr any) DeleteDao[T] {
	switch dao.dataSource.Driver {
	case sqlconst.Postgres:
		s, v := pgArray(arr)
		return dao.Where(fmt.Sprintf("%s @> %s", dao.modelMeta.escapeName(key), s), v...)
	default:
		panic(exception.New("WhereArrayIn not supported"))
	}
}
func (dao DeleteDao[T]) WhereArrayNotIn(key string, arr any) DeleteDao[T] {
	switch dao.dataSource.Driver {
	case sqlconst.Postgres:
		s, v := pgArray(arr)
		return dao.Where(fmt.Sprintf("not (%s @> %s)", dao.modelMeta.escapeName(key), s), v...)
	default:
		panic(exception.New("WhereArrayNotIn not supported"))
	}
}

func (dao DeleteDao[T]) WhereIn(key string, sub SubQueryInterface) DeleteDao[T] {
	sql, args := sub.sqlOriginPlaceholder()
	return dao.Where(squirrel.Expr(key+" IN ("+sql+")", args...))
}
func (dao DeleteDao[T]) WhereLike(field string, val string) DeleteDao[T] {
	return dao.Where(squirrel.Like{field: "%" + val + "%"})
}
