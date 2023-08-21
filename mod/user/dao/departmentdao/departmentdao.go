package departmentdao

import (
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

const (
	ResultDefault byte = iota
	ResultChildren
	ResultAll
	ResultNone
)

type Dao struct {
	sqlkit.Dao[model.Department]
}

func New(ds ...*sqlkit.DataSource) Dao {
	d := sqlkit.New[model.Department](ds...)
	dao := Dao{d}
	dao.Cascade = func(obj *model.Department) {
		switch dao.ResultType {
		case ResultChildren:
			obj.Children = dao.ListByParent(obj.Id)
			obj.Parent = nil
		case ResultDefault:
			if obj.Parent != nil {
				obj.Parent = dao.SelectOneById(obj.Parent.Id)
			}
		case ResultAll:
			obj.Children = dao.ListByParent(obj.Id)
			if obj.Parent != nil {
				obj.Parent = dao.SelectOneById(obj.Parent.Id)
			}
		case ResultNone:
			obj.Parent = nil
		}
	}
	return dao
}

func (dao Dao) ListByParent(id int32) []*model.Department {
	sql, args := dao.Builder().Select().Where("parent=?", id).OrderBy("no").OrderBy("id").Sql()
	return dao.ScanList(sql, args)
}

func (dao Dao) ListAll() []*model.Department {
	sql, args := dao.Builder().Select().Where("id>=0").OrderBy("parent").OrderBy("no").OrderBy("id").Sql()
	return dao.ScanList(sql, args)
}
