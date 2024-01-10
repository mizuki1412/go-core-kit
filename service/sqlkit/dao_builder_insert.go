package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
)

type InsertDao[T any] struct {
	Dao[T]
	builder squirrel.InsertBuilder
}

func (dao InsertDao[T]) Print() {
	sql, args := dao.Sql()
	logkit.Info("sql print", "sql", sql, "args", jsonkit.ToString(args))
}

func (dao InsertDao[T]) Exec() int64 {
	res := dao.ExecRaw(dao.Sql())
	rn, _ := res.RowsAffected()
	logkit.Debug("sql res", "rows", rn)
	return rn
}

func (dao InsertDao[T]) ReturnOne(dest *T) {
	rows := dao.QueryRaw(dao.Sql())
	defer rows.Close()
	for rows.Next() {
		// return 赋值
		err := rows.StructScan(dest)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		break
	}
}

func (dao InsertDao[T]) Sql() (string, []any) {
	sqls, args, err := dao.ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return sqls, args
}
func (dao InsertDao[T]) ToSql() (string, []any, error) {
	dao.builder = dao.builder.PlaceholderFormat(placeholder(dao.dataSource.Driver))
	return dao.builder.ToSql()
}

// SQL methods

// Prefix 在 sql 前写入语句
func (dao InsertDao[T]) Prefix(sql string, args ...any) InsertDao[T] {
	dao.builder = dao.builder.Prefix(sql, args...)
	return dao
}
func (dao InsertDao[T]) PrefixExpr(expr squirrel.Sqlizer) InsertDao[T] {
	dao.builder = dao.builder.PrefixExpr(expr)
	return dao
}
func (dao InsertDao[T]) Suffix(sql string, args ...any) InsertDao[T] {
	dao.builder = dao.builder.Suffix(sql, args...)
	return dao
}
func (dao InsertDao[T]) SuffixExpr(expr squirrel.Sqlizer) InsertDao[T] {
	dao.builder = dao.builder.SuffixExpr(expr)
	return dao
}

func (dao InsertDao[T]) Options(options ...string) InsertDao[T] {
	dao.builder = dao.builder.Options(options...)
	return dao
}

// Columns adds insert columns to the query.
func (dao InsertDao[T]) Columns(columns ...string) InsertDao[T] {
	dao.builder = dao.builder.Columns(dao.modelMeta.escapeNames(columns)...)
	return dao
}

// Values adds a single row's values to the query.
func (dao InsertDao[T]) Values(values ...any) InsertDao[T] {
	dao.builder = dao.builder.Values(values...)
	return dao
}

func (dao InsertDao[T]) Select(sb SelectDao[T]) InsertDao[T] {
	dao.builder = dao.builder.Values(sb.builder)
	return dao
}

//func (b InsertDao[T) SetMap(clauses map[string]any) InsertDao[T {
//	b.builder = b.builder.SetMap(clauses)
//	return b
//}
