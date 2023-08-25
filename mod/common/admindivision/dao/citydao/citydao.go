package citydao

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

type Dao struct {
	sqlkit.Dao[model.City]
}

func New(ds ...*sqlkit.DataSource) Dao {
	return Dao{sqlkit.New[model.City](ds...)}
}

func (dao Dao) FindCodeByName(name, pcode string) string {
	builder := dao.Builder().Select("code").Where("name=?", name).Where("province=?", pcode)
	rows := dao.Query(builder)
	defer rows.Close()
	for rows.Next() {
		ret, err := rows.SliceScan()
		if err != nil {
			panic(exception.New(err.Error()))
		}
		return ret[0].(string)
	}
	return ""
}

func (dao Dao) ListByProvince(id class.String) []*model.City {
	builder := dao.Builder().Select().Where("province=?", id).OrderBy("code")
	return dao.QueryList(builder)
}
