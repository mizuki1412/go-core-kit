package citydao

import (
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

type Dao struct {
	sqlkit.Dao[model.City]
}

func New(tx ...*sqlx.Tx) *Dao {
	dao := &Dao{}
	if len(tx) > 0 {
		dao.TX = tx[0]
	}
	return dao
}
func NewWithSchema(schema string, tx ...*sqlx.Tx) *Dao {
	dao := New(tx...)
	dao.SetSchema(schema)
	return dao
}

func (dao *Dao) FindById(id class.String) *model.City {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("city")).Where("code=?", id).MustSql()
	return dao.ScanOne(sql, args)
}

func (dao *Dao) FindCodeByName(name, pcode string) string {
	sql, args := sqlkit.Builder().Select("code").From(dao.GetTableD("city")).Where("name=?", name).Where("province=?", pcode).MustSql()
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

func (dao *Dao) ListByProvince(id class.String) []*model.City {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("city")).Where("province=?", id).OrderBy("code").MustSql()
	return dao.ScanList(sql, args)
}
