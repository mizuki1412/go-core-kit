package sqlkit

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/v2/class/const/sqlconst"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
	"strings"
)

type SelectDao[T any] struct {
	Dao[T]
	builder squirrel.SelectBuilder
	// from用默认的
	fromAs string
	// fromAs无效
	form           string
	ignoreLogicDel bool
}

type SubQueryInterface interface {
	// 默认占位符的，一般用于子查询
	sqlOriginPlaceholder() (string, []any)
}

func (dao SelectDao[T]) Print() {
	sql, args := dao.Sql()
	logkit.Info("sql print", "sql", sql, "args", jsonkit.ToString(args))
}

func (dao SelectDao[T]) sqlOriginPlaceholder() (string, []any) {
	if dao.fromAs == "" {
		dao.builder = dao.builder.From(dao.modelMeta.getTable())
	} else {
		dao.builder = dao.builder.From(dao.modelMeta.getTable(dao.fromAs))
	}
	dao.builder = dao.builder.PlaceholderFormat(squirrel.Question)
	return dao.builder.MustSql()
}

func (dao SelectDao[T]) Sql() (string, []any) {
	sqls, args, err := dao.ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return sqls, args
}
func (dao SelectDao[T]) ToSql() (string, []any, error) {
	if dao.from != "" {
		dao.builder = dao.builder.From(dao.from)
	} else {
		if dao.fromAs == "" {
			dao.builder = dao.builder.From(dao.modelMeta.getTable())
		} else {
			dao.builder = dao.builder.From(dao.modelMeta.getTable(dao.fromAs))
		}
	}
	dao.builder = dao.builder.PlaceholderFormat(placeholder(dao.dataSource.Driver))
	return dao.builder.ToSql()
}

func (dao SelectDao[T]) IgnoreLogicDel() SelectDao[T] {
	dao.ignoreLogicDel = true
	return dao
}

// SQL methods

// Prefix 在 sql 前写入语句
func (dao SelectDao[T]) Prefix(sql string, args ...any) SelectDao[T] {
	dao.builder = dao.builder.Prefix(sql, args...)
	return dao
}
func (dao SelectDao[T]) PrefixExpr(expr squirrel.Sqlizer) SelectDao[T] {
	dao.builder = dao.builder.PrefixExpr(expr)
	return dao
}
func (dao SelectDao[T]) Suffix(sql string, args ...any) SelectDao[T] {
	dao.builder = dao.builder.Suffix(sql, args...)
	return dao
}
func (dao SelectDao[T]) SuffixExpr(expr squirrel.Sqlizer) SelectDao[T] {
	dao.builder = dao.builder.SuffixExpr(expr)
	return dao
}

// Columns select 中额外增加 column
func (dao SelectDao[T]) Columns(cs ...string) SelectDao[T] {
	dao.builder = dao.builder.Columns(cs...)
	return dao
}

func (dao SelectDao[T]) RemoveColumns() SelectDao[T] {
	dao.builder = dao.builder.RemoveColumns()
	return dao
}
func (dao SelectDao[T]) From(from string) SelectDao[T] {
	dao.form = from
	return dao
}
func (dao SelectDao[T]) FromAs(alias string) SelectDao[T] {
	dao.fromAs = alias
	return dao
}
func (dao SelectDao[T]) FromSubQuery(sub SelectDao[T], alias string) SelectDao[T] {
	dao.builder = dao.builder.FromSelect(sub.builder, alias)
	return dao
}

func (dao SelectDao[T]) Distinct() SelectDao[T] {
	dao.builder = dao.builder.Distinct()
	return dao
}

// Options adds select option to the query
func (dao SelectDao[T]) Options(options ...string) SelectDao[T] {
	dao.builder = dao.builder.Options(options...)
	return dao
}

func (dao SelectDao[T]) Join(dm DaoModelMeta, as string, rest ...any) SelectDao[T] {
	dao.builder = dao.builder.Join(dm.getModelMeta().getTable(as), rest...)
	return dao
}
func (dao SelectDao[T]) LeftJoin(dm DaoModelMeta, as string, rest ...any) SelectDao[T] {
	dao.builder = dao.builder.LeftJoin(dm.getModelMeta().getTable(as), rest...)
	return dao
}
func (dao SelectDao[T]) RightJoin(dm DaoModelMeta, as string, rest ...any) SelectDao[T] {
	dao.builder = dao.builder.RightJoin(dm.getModelMeta().getTable(as), rest...)
	return dao
}
func (dao SelectDao[T]) InnerJoin(dm DaoModelMeta, as string, rest ...any) SelectDao[T] {
	dao.builder = dao.builder.InnerJoin(dm.getModelMeta().getTable(as), rest...)
	return dao
}
func (dao SelectDao[T]) CrossJoin(dm DaoModelMeta, as string, rest ...any) SelectDao[T] {
	dao.builder = dao.builder.CrossJoin(dm.getModelMeta().getTable(as), rest...)
	return dao
}

func (dao SelectDao[T]) JoinRaw(join string, rest ...any) SelectDao[T] {
	dao.builder = dao.builder.Join(join, rest...)
	return dao
}
func (dao SelectDao[T]) LeftJoinRaw(join string, rest ...any) SelectDao[T] {
	dao.builder = dao.builder.LeftJoin(join, rest...)
	return dao
}
func (dao SelectDao[T]) RightJoinRaw(join string, rest ...any) SelectDao[T] {
	dao.builder = dao.builder.RightJoin(join, rest...)
	return dao
}
func (dao SelectDao[T]) InnerJoinRaw(join string, rest ...any) SelectDao[T] {
	dao.builder = dao.builder.InnerJoin(join, rest...)
	return dao
}
func (dao SelectDao[T]) CrossJoinRaw(join string, rest ...any) SelectDao[T] {
	dao.builder = dao.builder.CrossJoin(join, rest...)
	return dao
}

func (dao SelectDao[T]) Where(pred any, args ...any) SelectDao[T] {
	dao.builder = dao.builder.Where(pred, args...)
	return dao
}
func (dao SelectDao[T]) Having(pred any, rest ...any) SelectDao[T] {
	dao.builder = dao.builder.Where(pred, rest...)
	return dao
}

func (dao SelectDao[T]) GroupBy(groupBys ...string) SelectDao[T] {
	dao.builder = dao.builder.GroupBy(dao.modelMeta.escapeNames(groupBys)...)
	return dao
}
func (dao SelectDao[T]) GroupByRow(groupBys ...string) SelectDao[T] {
	dao.builder = dao.builder.GroupBy(groupBys...)
	return dao
}

func (dao SelectDao[T]) OrderBy(field string) SelectDao[T] {
	if strings.Contains(field, " ") {
		panic(exception.New("order by 不能包含空格"))
	}
	dao.builder = dao.builder.OrderBy(dao.modelMeta.escapeName(field))
	return dao
}
func (dao SelectDao[T]) OrderByDesc(field string) SelectDao[T] {
	dao.builder = dao.builder.OrderBy(dao.modelMeta.escapeName(field) + " DESC")
	return dao
}

func (dao SelectDao[T]) Limit(limit uint64) SelectDao[T] {
	dao.builder = dao.builder.Limit(limit)
	return dao
}

func (dao SelectDao[T]) Offset(offset uint64) SelectDao[T] {
	dao.builder = dao.builder.Offset(offset)
	return dao
}

// custom

// 重置select选项
func (dao SelectDao[T]) resetColumns(fields ...string) SelectDao[T] {
	dao.builder = dao.builder.RemoveColumns().Columns(fields...)
	return dao
}

func (dao SelectDao[T]) whereNLogicDel() SelectDao[T] {
	if dao.modelMeta.logicDelKey.Key != "" {
		return dao.Where(squirrel.NotEq{dao.modelMeta.logicDelKey.Key: dao.LogicDelVal[0]})
	}
	return dao
}

// 生成sql中: sth in (select unnest(Array[?,?,?])) []any
// 注意使用时 args...
func (dao SelectDao[T]) whereUnnest(arr any, key, flag string) SelectDao[T] {
	switch dao.dataSource.Driver {
	case sqlconst.Postgres:
		s, v := pgArray(arr)
		return dao.Where(fmt.Sprintf("%s %s (select unnest(%s))", dao.modelMeta.escapeName(key), flag, s), v...)
	default:
		panic(exception.New("whereUnnest not supported"))
	}
}
func (dao SelectDao[T]) WhereUnnestIn(key string, arr any) SelectDao[T] {
	return dao.whereUnnest(arr, key, "IN")
}
func (dao SelectDao[T]) WhereUnnestNotIn(key string, arr any) SelectDao[T] {
	return dao.whereUnnest(arr, key, "NOT IN")
}

// WherePGArrayIn 用于PG中array类型数据的包含比较
func (dao SelectDao[T]) WherePGArrayIn(key string, arr any) SelectDao[T] {
	switch dao.dataSource.Driver {
	case sqlconst.Postgres:
		s, v := pgArray(arr)
		return dao.Where(fmt.Sprintf("%s @> %s", dao.modelMeta.escapeName(key), s), v...)
	default:
		panic(exception.New("WherePGArrayIn not supported"))
	}
}
func (dao SelectDao[T]) WherePGArrayNotIn(key string, arr any) SelectDao[T] {
	switch dao.dataSource.Driver {
	case sqlconst.Postgres:
		s, v := pgArray(arr)
		return dao.Where(fmt.Sprintf("not (%s @> %s)", dao.modelMeta.escapeName(key), s), v...)
	default:
		panic(exception.New("WherePGArrayNotIn not supported"))
	}
}

func (dao SelectDao[T]) WhereIn(key string, sub SubQueryInterface) SelectDao[T] {
	sql, args := sub.sqlOriginPlaceholder()
	return dao.Where(squirrel.Expr(key+" IN ("+sql+")", args...))
}
func (dao SelectDao[T]) WhereLike(field string, val string) SelectDao[T] {
	return dao.Where(squirrel.Like{field: "%" + val + "%"})
}
