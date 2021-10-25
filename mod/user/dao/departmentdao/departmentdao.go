package departmentdao

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

/// auto template
type Dao struct {
	sqlkit.Dao
}

const (
	ResultDefault byte = iota
	ResultChildren
	ResultAll
	ResultNone
)

func New(schema string, tx ...*sqlkit.Dao) *Dao {
	dao := &Dao{}
	dao.NewHelper(schema, tx...)
	return dao
}
func (dao *Dao) cascade(obj *model.Department) {
	switch dao.ResultType {
	case ResultChildren:
		obj.Children = dao.ListByParent(obj.Id)
		obj.Parent = nil
	case ResultDefault:
		if obj.Parent != nil {
			obj.Parent = dao.FindById(obj.Parent.Id)
		}
	case ResultAll:
		obj.Children = dao.ListByParent(obj.Id)
		if obj.Parent != nil {
			obj.Parent = dao.FindById(obj.Parent.Id)
		}
	case ResultNone:
		obj.Parent = nil
	}
}
func (dao *Dao) scan(sql string, args []interface{}) []*model.Department {
	rows := dao.Query(sql, args...)
	list := make([]*model.Department, 0, 5)
	defer rows.Close()
	for rows.Next() {
		m := &model.Department{}
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
func (dao *Dao) scanOne(sql string, args []interface{}) *model.Department {
	rows := dao.Query(sql, args...)
	defer rows.Close()
	for rows.Next() {
		m := model.Department{}
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

func (dao *Dao) FindById(id int32) *model.Department {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("department")).Where("id=?", id).MustSql()
	return dao.scanOne(sql, args)
}

func (dao *Dao) ListByParent(id int32) []*model.Department {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("department")).Where("parent=?", id).OrderBy("no", "id").MustSql()
	return dao.scan(sql, args)
}

func (dao *Dao) ListAll() []*model.Department {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("department")).Where("id>=0").OrderBy("parent", "no", "id").MustSql()
	return dao.scan(sql, args)
}
