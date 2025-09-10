package sqlkit

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
)

// LogicDelVal 全局逻辑删除的 value
var LogicDelVal = []any{true, false}

type Dao[T any] struct {
	meta T
	// 逻辑删除的字段，可替代全局的LogicDelVal
	LogicDelVal []any
	// 返回级联的类型
	//ResultType byte
	// 级联实现的函数
	Cascade func(*T)
	// 数据源
	dataSource *DataSource
	// 目标表结构
	modelMeta ModelMeta
}

type DaoModelMeta interface {
	getModelMeta() ModelMeta
}

// New 必须从初始化函数生成 dao
func New[T any](ds ...*DataSource) Dao[T] {
	dao := Dao[T]{}
	if len(ds) > 0 {
		dao.dataSource = ds[0]
	} else {
		dao.dataSource = DefaultDataSource()
	}
	dao.modelMeta.dateSource = dao.dataSource
	dao.modelMeta = dao.modelMeta.init(dao.meta)
	return dao
}

func (dao Dao[T]) getModelMeta() ModelMeta {
	return dao.modelMeta
}

func (dao Dao[T]) DataSource() *DataSource {
	return dao.dataSource
}

func (dao Dao[T]) QueryRaw(sql string, args []any) *sqlx.Rows {
	logkit.Debug("sql req", "sql", sql, "args", jsonkit.ToString(args))
	return dao.dataSource.Query(sql, args)
}

// QueryRawRows 默认返回 T list，可用于自由sql时，自定义返回值
func (dao Dao[T]) QueryRawRows(sql string, args []any) []*T {
	rows := dao.QueryRaw(sql, args)
	list := make([]*T, 0, 5)
	defer rows.Close()
	for rows.Next() {
		list = append(list, scanStruct[T](rows, dao.dataSource.Driver))
	}
	if dao.Cascade != nil {
		for i := range list {
			dao.Cascade(list[i])
		}
	}
	return list
}

func (dao Dao[T]) ExecRaw(sql string, args []any) sql.Result {
	logkit.Debug("sql exec", "sql", sql, "args", jsonkit.ToString(args))
	return dao.dataSource.Exec(sql, args)
}

/// 小功能

func (dao Dao[T]) Table(alias ...string) string {
	return dao.modelMeta.getTable(alias...)
}
func (dao Dao[T]) EscapeNames(name ...string) []string {
	return dao.modelMeta.escapeNames(name)
}
