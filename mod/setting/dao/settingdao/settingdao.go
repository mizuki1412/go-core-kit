package settingdao

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/mod/setting/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

type Dao struct {
	sqlkit.Dao[model.Setting]
}

func New(ds ...*sqlkit.DataSource) Dao {
	return Dao{sqlkit.New[model.Setting](ds...)}
}

func (dao Dao) Set(data map[string]interface{}) {
	sql, args, err := dao.Builder().Update().Set("data", jsonkit.ToString(data)).Where("id=?", 1).ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	dao.Exec(sql, args...)
}

func (dao Dao) Get() map[string]interface{} {
	sql, args := dao.Builder().Select("data").Where("id=?", 1).Sql()
	rows := dao.Query(sql, args...)
	var data string
	defer rows.Close()
	for rows.Next() {
		_ = rows.Scan(&data)
	}
	return jsonkit.ParseMap(data)
}
