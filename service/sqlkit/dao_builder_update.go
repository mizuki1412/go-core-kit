package sqlkit

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/v2/class/const/sqlconst"
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

// custom 参考dao_builder_select

func (dao UpdateDao[T]) whereUnnest(arr any, key, flag string) UpdateDao[T] {
	switch dao.dataSource.Driver {
	case sqlconst.Postgres:
		s, v := pgArray(arr)
		return dao.Where(fmt.Sprintf("%s %s (select unnest(%s))", dao.modelMeta.escapeName(key), flag, s), v...)
	default:
		panic(exception.New("whereUnnest not supported"))
	}
}
func (dao UpdateDao[T]) WhereUnnestIn(key string, arr any) UpdateDao[T] {
	return dao.whereUnnest(arr, key, "IN")
}
func (dao UpdateDao[T]) WhereUnnestNotIn(key string, arr any) UpdateDao[T] {
	return dao.whereUnnest(arr, key, "NOT IN")
}

// WherePGArrayIn 用于PG中array类型数据的包含比较
func (dao UpdateDao[T]) WherePGArrayIn(key string, arr any) UpdateDao[T] {
	switch dao.dataSource.Driver {
	case sqlconst.Postgres:
		s, v := pgArray(arr)
		return dao.Where(fmt.Sprintf("%s @> %s", dao.modelMeta.escapeName(key), s), v...)
	default:
		panic(exception.New("WherePGArrayIn not supported"))
	}
}
func (dao UpdateDao[T]) WherePGArrayNotIn(key string, arr any) UpdateDao[T] {
	switch dao.dataSource.Driver {
	case sqlconst.Postgres:
		s, v := pgArray(arr)
		return dao.Where(fmt.Sprintf("not (%s @> %s)", dao.modelMeta.escapeName(key), s), v...)
	default:
		panic(exception.New("WherePGArrayNotIn not supported"))
	}
}

func (dao UpdateDao[T]) WhereIn(key string, sub SubQueryInterface) UpdateDao[T] {
	sql, args := sub.sqlOriginPlaceholder()
	return dao.Where(squirrel.Expr(key+" IN ("+sql+")", args...))
}
func (dao UpdateDao[T]) WhereLike(field string, val string) UpdateDao[T] {
	return dao.Where(squirrel.Like{field: "%" + val + "%"})
}
