package departmentdao

import (
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

type Dao struct {
	sqlkit.Dao[model.Department]
}

const (
	ResultDefault byte = iota
	ResultChildren
	ResultAll
	ResultNone
)

func New(tx ...*sqlx.Tx) *Dao {
	dao := &Dao{}
	if len(tx) > 0 {
		dao.TX = tx[0]
	}
	dao.Cascade = func(obj *model.Department) {
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
	return dao
}
func NewWithSchema(schema string, tx ...*sqlx.Tx) *Dao {
	dao := New(tx...)
	dao.SetSchema(schema)
	return dao
}

func (dao *Dao) FindById(id int32) *model.Department {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("department")).Where("id=?", id).MustSql()
	return dao.ScanOne(sql, args)
}

func (dao *Dao) ListByParent(id int32) []*model.Department {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("department")).Where("parent=?", id).OrderBy("no", "id").MustSql()
	return dao.ScanList(sql, args)
}

func (dao *Dao) ListAll() []*model.Department {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("department")).Where("id>=0").OrderBy("parent", "no", "id").MustSql()
	return dao.ScanList(sql, args)
}
