package areadao

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

/// auto template
type Dao struct {
	sqlkit.Dao
}

func New(schema string, tx ...*sqlkit.Dao) *Dao {
	dao := &Dao{}
	dao.NewHelper(schema, tx...)
	return dao
}
func (dao *Dao) scan(sql string, args []any) []*model.Area {
	rows := dao.Query(sql, args...)
	list := make([]*model.Area, 0, 5)
	defer rows.Close()
	for rows.Next() {
		m := &model.Area{}
		err := rows.StructScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, m)
	}
	return list
}
func (dao *Dao) scanOne(sql string, args []any) *model.Area {
	rows := dao.Query(sql, args...)
	defer rows.Close()
	for rows.Next() {
		m := model.Area{}
		err := rows.StructScan(&m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		return &m
	}
	return nil
}

////

func (dao *Dao) FindById(id class.String) *model.Area {
	sql, args := sqlkit.Builder().Select("code,name").From(dao.GetTableD("area")).Where("code=?", id).MustSql()
	return dao.scanOne(sql, args)
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
	return dao.scan(sql, args)
}
