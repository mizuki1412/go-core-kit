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

var meta = sqlkit.InitModelMeta(&model.City{})

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

func (dao *Dao) FindById(id class.String) *model.City {
	sql, args := sqlkit.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where("code=?", id).MustSql()
	return dao.ScanOne(sql, args)
}

func (dao *Dao) FindCodeByName(name, pcode string) string {
	sql, args := sqlkit.Builder().Select("code").From(sqlkit.GetSchemaTable(dao.Schema, "city")).Where("name=?", name).Where("province=?", pcode).MustSql()
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
	sql, args := sqlkit.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where("province=?", id).OrderBy("code").MustSql()
	return dao.ScanList(sql, args)
}
