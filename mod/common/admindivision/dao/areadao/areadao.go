package areadao

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

type Dao struct {
	sqlkit.Dao[model.Area]
}

func New(ds ...*sqlkit.DataSource) Dao {
	return Dao{sqlkit.New[model.Area](ds...)}
}

func (dao Dao) FindCodeByName(name, ccode, pcode string) string {
	builder := dao.Builder().Select("code").Where("name=?", name).Where("city=?", ccode).Where("province=?", pcode)
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

func (dao Dao) ListByCity(id class.String) []*model.Area {
	builder := dao.Builder().Select().Where("city=?", id).OrderBy("code")
	return dao.QueryList(builder)
}
