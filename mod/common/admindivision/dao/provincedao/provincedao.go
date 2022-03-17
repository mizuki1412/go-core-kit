package provincedao

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/dao/citydao"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

/// auto template
type Dao struct {
	sqlkit.Dao
}

const (
	ResultDefault byte = iota
	ResultNone
)

func New(schema string, tx ...*sqlkit.Dao) *Dao {
	dao := &Dao{}
	dao.NewHelper(schema, tx...)
	return dao
}
func (dao *Dao) cascade(obj *model.Province) {
	switch dao.ResultType {
	case ResultDefault:
		obj.Cities = citydao.New(dao.Schema).ListByProvince(obj.Code)
	}
}
func (dao *Dao) scan(sql string, args []any) []*model.Province {
	rows := dao.Query(sql, args...)
	list := make([]*model.Province, 0, 5)
	defer rows.Close()
	for rows.Next() {
		m := &model.Province{}
		err := rows.StructScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, m)
	}
	for i := range list {
		dao.cascade(list[i])
	}
	return list
}
func (dao *Dao) scanOne(sql string, args []any) *model.Province {
	rows := dao.Query(sql, args...)
	defer rows.Close()
	for rows.Next() {
		m := model.Province{}
		err := rows.StructScan(&m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		dao.cascade(&m)
		return &m
	}
	return nil
}

////

func (dao *Dao) FindById(id class.String) *model.Province {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("province")).Where("code=?", id).MustSql()
	return dao.scanOne(sql, args)
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
	return dao.scan(sql, args)
}
