package roledao

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/departmentdao"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

type Dao struct {
	sqlkit.Dao[model.Role]
}

const (
	ResultDefault byte = iota
	ResultNone
)

func New(ds ...*sqlkit.DataSource) Dao {
	dao := sqlkit.New[model.Role](ds...)
	dao.Cascade = func(obj *model.Role) {
		switch dao.ResultType {
		case ResultDefault:
			if obj.Department != nil {
				obj.Department = departmentdao.New(dao.DataSource()).SelectOneById(obj.Department.Id)
			}
		case ResultNone:
			obj.Department = nil
		}
	}
	return Dao{dao}
}

func (dao Dao) FindByName(name string) *model.Role {
	builder := dao.Builder().Select().Where("name=?", name).Limit(1)
	return dao.QueryOne(builder)
}
func (dao Dao) ListFromRootDepart(id int32) []*model.Role {
	where := fmt.Sprintf(`id>0 and department in ( with recursive t(id) as( values(%d) union all select d.id from %s d, t where t.id=d.parent) select id from t )`, id, departmentdao.New(dao.DataSource()).Table())
	builder := dao.Builder().Select().Where(where).OrderBy("id")
	return dao.QueryList(builder)
}

type ListParam struct {
	Departments []int32
}

func (dao Dao) List(param ListParam) []*model.Role {
	builder := dao.Builder().Select().Where("id>0 and department>=0").OrderBy("id")
	if len(param.Departments) > 0 {
		builder = builder.WhereUnnestIn("department", param.Departments)
	}
	return dao.QueryList(builder)
}
func (dao Dao) ListByDepartment(did int32) []*model.Role {
	builder := dao.Builder().Select().Where("id>0 and department=?", did).OrderBy("id")
	return dao.QueryList(builder)
}
