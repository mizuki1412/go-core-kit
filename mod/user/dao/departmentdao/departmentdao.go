package departmentdao

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/model"
	"github.com/mizuki1412/go-core-kit/v2/service/sqlkit"
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

func New(cascadeType byte, ds ...*sqlkit.DataSource) Dao {
	d := sqlkit.New[model.Department](ds...)
	dao := Dao{d}
	dao.Cascade = func(obj *model.Department) {
		switch cascadeType {
		case ResultChildren:
			obj.Children = dao.ListByParent(obj.Id)
			obj.Parent = nil
		case ResultDefault:
			if obj.Parent != nil {
				obj.Parent = dao.SelectOneWithDelById(obj.Parent.Id)
			}
		case ResultAll:
			obj.Children = dao.ListByParent(obj.Id)
			if obj.Parent != nil {
				obj.Parent = dao.SelectOneWithDelById(obj.Parent.Id)
			}
		case ResultNone:
			obj.Parent = nil
		}
	}
	return dao
}

func (dao Dao) ListByParent(id int64) model.DeptList {
	return dao.Select().Where(squirrel.Eq{"parent": id}).OrderBy("no").OrderBy("id").List()
}

func (dao Dao) ListAll() model.DeptList {
	return dao.Select().Where(squirrel.Gt{"id": 0}).OrderBy("parent").OrderBy("no").OrderBy("id").List()
}

func (dao Dao) FindByName(name string) *model.Department {
	return dao.Select().Where(squirrel.Eq{"name": name}).Limit(1).One()
}

func (dao Dao) FindByNo(no string) *model.Department {
	return dao.Select().Where(squirrel.Eq{"no": no}).Limit(1).One()
}
