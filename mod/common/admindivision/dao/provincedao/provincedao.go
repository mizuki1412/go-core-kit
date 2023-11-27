package provincedao

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/dao/citydao"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

type Dao struct {
	sqlkit.Dao[model.Province]
}

const (
	ResultDefault byte = iota
	ResultNone
)

func New(cascadeType byte, ds ...*sqlkit.DataSource) Dao {
	dao := Dao{sqlkit.New[model.Province](ds...)}
	dao.Cascade = func(obj *model.Province) {
		switch cascadeType {
		case ResultDefault:
			obj.Cities = citydao.New(dao.DataSource()).ListByProvince(obj.Code)
		}
	}
	return dao
}

func (dao Dao) FindCodeByName(name string) string {
	rows := dao.Select("code").Where("name=?", name).QueryRows()
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

func (dao Dao) ListAll() []*model.Province {
	return dao.Select().OrderBy("code").List()
}
