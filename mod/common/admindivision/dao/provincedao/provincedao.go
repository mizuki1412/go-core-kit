package provincedao

import (
	"github.com/mizuki1412/go-core-kit/class"
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

func New(ds ...*sqlkit.DataSource) Dao {
	dao := Dao{sqlkit.New[model.Province](ds...)}
	dao.Cascade = func(obj *model.Province) {
		switch dao.ResultType {
		case ResultDefault:
			obj.Cities = citydao.New(dao.DataSource()).ListByProvince(obj.Code)
		}
	}
	return dao
}

func (dao Dao) FindById(id class.String) *model.Province {
	sql, args := dao.Builder().Select().Where("code=?", id).Sql()
	return dao.ScanOne(sql, args)
}

func (dao Dao) FindCodeByName(name string) string {
	sql, args := dao.Builder().Select("code").Where("name=?", name).Sql()
	rows := dao.Query(sql, args)
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
	sql, args := dao.Builder().Select().OrderBy("code").Sql()
	return dao.ScanList(sql, args)
}
