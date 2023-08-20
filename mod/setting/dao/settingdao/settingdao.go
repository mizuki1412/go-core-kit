package settingdao

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

/** more_setting: id, data */

type Dao struct {
	sqlkit.Dao[map[string]any]
}

func New(ds ...*sqlkit.DataSource) Dao {
	dao := Dao{}
	if len(ds) > 0 {
		dao.SetDataSource(ds[0])
	}
	return dao
}

func (dao *Dao) Set(data map[string]interface{}) {
	sql, args, err := dao.Builder().Update(sqlkit.GetSchemaTable(dao.Schema, "more_setting")).Set("data", jsonkit.ToString(data)).Where("id=?", 1).ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	dao.Exec(sql, args...)
}

func (dao *Dao) Get() map[string]interface{} {
	sql, args := dao.Builder().Select("data").From(sqlkit.GetSchemaTable(dao.Schema, "more_setting")).Where("id=?", 1).MustSql()
	rows := dao.Query(sql, args...)
	var data string
	defer rows.Close()
	for rows.Next() {
		_ = rows.Scan(&data)
	}
	return jsonkit.ParseMap(data)
}
