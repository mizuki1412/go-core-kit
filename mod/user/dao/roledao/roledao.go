package roledao

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/dao/departmentdao"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/model"
	"github.com/mizuki1412/go-core-kit/v2/service/sqlkit"
)

type Dao struct {
	sqlkit.Dao[model.Role]
}

const (
	ResultDefault byte = iota
	ResultNone
)

func New(cascadeType byte, ds ...*sqlkit.DataSource) Dao {
	dao := sqlkit.New[model.Role](ds...)
	dao.Cascade = func(obj *model.Role) {
		switch cascadeType {
		case ResultDefault:
			if obj.Department != nil {
				obj.Department = departmentdao.New(departmentdao.ResultDefault, dao.DataSource()).SelectOneWithDelById(obj.Department.Id)
			}
		case ResultNone:
			obj.Department = nil
		}
	}
	return Dao{dao}
}

func (dao Dao) FindByName(name string) *model.Role {
	return dao.Select().Where(squirrel.Eq{"name": name}).Limit(1).One()
}
func (dao Dao) ListFromRootDepart(id int64) []*model.Role {
	where := fmt.Sprintf(`id>0 and department in ( with recursive t(id) as( values(%d::bigint) union all select d.id from %s d, t where t.id=d.parent) select id from t )`, id, departmentdao.New(departmentdao.ResultDefault, dao.DataSource()).Table())
	return dao.Select().Where(where).OrderBy("id").List()
}

func (dao Dao) CountFromRootDepart(id int64) int64 {
	where := fmt.Sprintf(`department in ( with recursive t(id) as( values(%d::bigint) union all select d.id from %s d, t where t.id=d.parent) select id from t )`, id, departmentdao.New(departmentdao.ResultDefault, dao.DataSource()).Table())
	return dao.Select().Where(where).Count()
}

type ListParam struct {
	Departments []int64
}

func (dao Dao) List(param ListParam) []*model.Role {
	builder := dao.Select().Where("id>0 and department>=0").OrderBy("id")
	if len(param.Departments) > 0 {
		builder = builder.WhereUnnestIn("department", param.Departments)
	}
	return builder.List()
}

func (dao Dao) Count(param ListParam) int64 {
	builder := dao.Select()
	if len(param.Departments) > 0 {
		builder = builder.WhereUnnestIn("department", param.Departments)
	}
	return builder.Count()
}
