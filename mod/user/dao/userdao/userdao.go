package userdao

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/mizuki1412/go-core-kit/v2/library/stringkit"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/dao/departmentdao"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/dao/roledao"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/model"
	"github.com/mizuki1412/go-core-kit/v2/service/sqlkit"
	"github.com/spf13/cast"
)

type Dao struct {
	sqlkit.Dao[model.User]
}

const (
	ResultDefault byte = iota
	ResultNone
)

func New(cascadeType byte, ds ...*sqlkit.DataSource) Dao {
	dao := sqlkit.New[model.User](ds...)
	//dao.LogicDelVal = []any{-1, 0}
	dao.Cascade = func(obj *model.User) {
		switch cascadeType {
		case ResultDefault:
			if obj.Role != nil {
				obj.Role = roledao.New(roledao.ResultDefault, dao.DataSource()).SelectOneWithDelById(obj.Role.Id)
			}
			if obj.Department != nil {
				obj.Department = departmentdao.New(departmentdao.ResultDefault, dao.DataSource()).SelectOneWithDelById(obj.Department.Id)
			}
		case ResultNone:
			obj.Role = nil
			obj.Department = nil
		}
	}
	return Dao{dao}
}

func (dao Dao) Login(pwd, username, phone string) *model.User {
	builder := dao.Select()
	if !stringkit.IsNull(username) {
		builder = builder.Where("username=?", username)
	} else {
		builder = builder.Where("phone=?", phone)
	}
	return builder.Where("pwd=?", pwd).Limit(1).One()
}

func (dao Dao) FindByPhone(phone string) *model.User {
	return dao.Select().Where("phone=?", phone).One()
}

func (dao Dao) FindByUsername(username string) *model.User {
	return dao.Select().Where("username=?", username).One()
}
func (dao Dao) FindByUsernameDeleted(username string) *model.User {
	return dao.Select().Where("username=?", username).One()
}

// FindParam 可以通过extend的值来find
type FindParam struct {
	Extend map[string]any
}

func (dao Dao) Find(param FindParam) *model.User {
	builder := dao.Select().Limit(1)
	for k, v := range param.Extend {
		builder = builder.Where(fmt.Sprintf("extend->>'%s'=?", k), cast.ToString(v))
	}
	return builder.One()
}

func (dao Dao) ListFromRootDepart(departId int64) []*model.User {
	builder := dao.Select()
	//if len(roleIds) > 0 {
	//	builder = builder.WhereUnnestIn("role", roleIds)
	//}
	where := fmt.Sprintf(`department in(with recursive t(id) as( values(%d) union all select d.id from %s d, t where t.id=d.parent) select id from t )`,
		departId,
		departmentdao.New(departmentdao.ResultDefault, dao.DataSource()).Table())
	builder = builder.Where(where)
	return builder.OrderBy("name").OrderBy("id").List()
}

func (dao Dao) CountFromRootDepart(departId int64) int64 {
	builder := dao.Select()
	where := fmt.Sprintf(`department in(with recursive t(id) as( values(%d) union all select d.id from %s d, t where t.id=d.parent) select id from t )`,
		departId,
		departmentdao.New(departmentdao.ResultDefault, dao.DataSource()).Table())
	builder = builder.Where(where)
	return builder.Count()
}

type ListParam struct {
	RoleId      class.Int64
	Roles       []int64
	Departments []int64
	IdList      []int64
}

func (dao Dao) List(param ListParam) []*model.User {
	builder := dao.Select()
	if param.RoleId.IsValid() {
		builder = builder.Where("role=?", param.RoleId)
	}
	if len(param.IdList) > 0 {
		builder = builder.WhereUnnestIn("id", param.IdList)
	}
	if len(param.Roles) > 0 {
		builder = builder.WhereUnnestIn("role", param.Roles)
	}
	if len(param.Departments) > 0 {
		//rb := roledao.New(roledao.ResultDefault).Select("id").WhereUnnestIn("department", param.Departments)
		//builder = builder.WhereIn("role", rb)
		builder = builder.WhereUnnestIn("department", param.Departments)
	}
	return builder.List()
}

func (dao Dao) FreezeUser(uid int64, status int32) {
	dao.Update().Set("status", status).Where("id=?", uid).Exec()
}
func (dao Dao) SetNull(id int64) {
	dao.Update().Set("phone", squirrel.Expr("null")).Set("username", squirrel.Expr("null")).Where("id=?", id).Exec()
}
