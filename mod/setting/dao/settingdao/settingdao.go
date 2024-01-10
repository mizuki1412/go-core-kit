package settingdao

import (
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/v2/mod/setting/model"
	"github.com/mizuki1412/go-core-kit/v2/service/sqlkit"
)

type Dao struct {
	sqlkit.Dao[model.Setting]
}

func New(ds ...*sqlkit.DataSource) Dao {
	return Dao{sqlkit.New[model.Setting](ds...)}
}

func (dao Dao) Set(data map[string]any) {
	dao.Update().Set("data", jsonkit.ToString(data)).Where("id=?", 1).Exec()
}

func (dao Dao) Get() map[string]any {
	builder := dao.Select("data").Where("id=?", 1)
	rows := dao.QueryRaw(builder.Sql())
	var data string
	defer rows.Close()
	for rows.Next() {
		_ = rows.Scan(&data)
	}
	return jsonkit.ParseMap(data)
}
