package provincedao

import (
	"github.com/jmoiron/sqlx"
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

func New(tx ...*sqlx.Tx) *Dao {
	return NewWithSchema("", tx...)
}
func NewWithSchema(schema string, tx ...*sqlx.Tx) *Dao {
	dao := &Dao{}
	dao.SetSchema(schema)
	if len(tx) > 0 {
		dao.TX = tx[0]
	}
	dao.Cascade = func(obj *model.Province) {
		switch dao.ResultType {
		case ResultDefault:
			obj.Cities = citydao.NewWithSchema(dao.Schema).ListByProvince(obj.Code)
		}
	}
	return dao
}

func (dao *Dao) FindById(id class.String) *model.Province {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("province")).Where("code=?", id).MustSql()
	return dao.ScanOne(sql, args)
}

func (dao *Dao) FindCodeByName(name string) string {
	sql, args := sqlkit.Builder().Select("code").From(dao.GetTableD("province")).Where("name=?", name).MustSql()
	rows := dao.Query(sql, args...)
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

func (dao *Dao) ListAll() []*model.Province {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("province")).OrderBy("code").MustSql()
	return dao.ScanList(sql, args)
}
