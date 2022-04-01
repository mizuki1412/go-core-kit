package areadao

import (
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

type Dao struct {
	sqlkit.Dao[model.Area]
}

func New(tx ...*sqlx.Tx) *Dao {
	return NewWithSchema("", tx...)
}
func NewWithSchema(schema string, tx ...*sqlx.Tx) *Dao {
	dao := &Dao{}
	dao.SetSchema(schema)
	if len(tx) > 0 {
		dao.TX = tx[0]
	}
	return dao
}

func (dao *Dao) FindById(id class.String) *model.Area {
	sql, args := sqlkit.Builder().Select("code,name").From(dao.GetTableD("area")).Where("code=?", id).MustSql()
	return dao.ScanOne(sql, args)
}

func (dao *Dao) FindCodeByName(name, ccode, pcode string) string {
	sql, args := sqlkit.Builder().Select("code").From(dao.GetTableD("area")).Where("name=?", name).Where("city=?", ccode).Where("province=?", pcode).MustSql()
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

func (dao *Dao) ListByCity(id class.String) []*model.Area {
	sql, args := sqlkit.Builder().Select("code,name").From(dao.GetTableD("area")).Where("city=?", id).OrderBy("code").MustSql()
	return dao.ScanList(sql, args)
}
