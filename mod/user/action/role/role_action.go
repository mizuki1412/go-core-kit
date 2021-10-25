package role

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/departmentdao"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/roledao"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/userdao"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
)

func listAllPrivileges(ctx *context.Context) {
	ctx.JsonSuccess(roledao.New(ctx.SessionGetSchema()).ListPrivileges())
}

type createParams struct {
	Name           string          `validate:"required"`
	PrivilegesJson class.ArrString `validate:"required" default:"[]" description:"数组json字符串：[a,b,c]"`
	DepartmentId   int32
}

func create(ctx *context.Context) {
	params := createParams{}
	ctx.BindForm(&params)
	departmentDao := departmentdao.New(ctx.SessionGetSchema())
	departmentDao.SetResultType(departmentdao.ResultNone)
	department := departmentDao.FindById(params.DepartmentId)
	if department == nil {
		panic(exception.New("部门不存在"))
	}
	role := &model.Role{}
	role.Name.Set(params.Name)
	role.Privileges = params.PrivilegesJson
	role.Department = department
	roledao.New(ctx.SessionGetSchema()).Insert(role)
	ctx.JsonSuccess(nil)
}

type updateParams struct {
	Id             int32 `validate:"required"`
	Name           class.String
	PrivilegesJson class.ArrString `description:"数组json字符串：[a,b,c]"`
	DepartmentId   class.Int32
}

func update(ctx *context.Context) {
	params := updateParams{}
	ctx.BindForm(&params)
	dao := roledao.New(ctx.SessionGetSchema())
	role := dao.FindById(params.Id)
	if role == nil {
		panic(exception.New("角色不存在"))
	}
	if params.DepartmentId.Valid && (role.Department == nil || params.DepartmentId.Int32 != role.Department.Id) {
		departmentDao := departmentdao.New(ctx.SessionGetSchema())
		departmentDao.SetResultType(departmentdao.ResultNone)
		d := departmentDao.FindById(params.DepartmentId.Int32)
		if d == nil {
			panic(exception.New("部门不存在"))
		}
		role.Department = d
	}
	if params.Name.Valid {
		role.Name.Set(params.Name.String)
	}
	if params.PrivilegesJson.Valid {
		role.Privileges = params.PrivilegesJson
	}
	dao.Update(role)
	ctx.JsonSuccess(nil)
}

type delParams struct {
	Id int32 `validate:"required"`
}

func del(ctx *context.Context) {
	params := delParams{}
	ctx.BindForm(&params)
	dao := roledao.New(ctx.SessionGetSchema())
	dao.SetResultType(roledao.ResultNone)
	role := dao.FindById(params.Id)
	if role == nil {
		panic(exception.New("角色不存在"))
	}
	if val, ok := role.Extend.Map["immutable"]; ok && val.(bool) {
		panic(exception.New("该角色不可删除"))
	}
	userDao := userdao.New(ctx.SessionGetSchema())
	userDao.SetResultType(userdao.ResultNone)
	us := userDao.List(userdao.ListParam{RoleId: params.Id})
	if us != nil && len(us) > 0 {
		panic(exception.New("角色下还有用户,不能删除"))
	}
	dao.Delete(role)
	ctx.JsonSuccess(nil)
}

type listRolesParam struct {
	Root class.Int32 `description:"指定根department"`
}

func listRoles(ctx *context.Context) {
	params := listRolesParam{}
	ctx.BindForm(&params)
	var roles []*model.Role
	if params.Root.Valid {
		roles = roledao.New(ctx.SessionGetSchema()).ListFromRootDepart(params.Root.Int32)
	} else {
		roles = roledao.New(ctx.SessionGetSchema()).List(roledao.ListParam{})
	}
	for _, r := range roles {
		if !r.Privileges.Valid {
			r.Privileges.Valid = true
			r.Privileges.Array = []string{}
		}
	}
	ctx.JsonSuccess(roles)
}

type listByRoleParams struct {
	RoleId int32 `validate:"required"`
}

func listRolesWithUser(ctx *context.Context) {
	params := listByRoleParams{}
	ctx.BindForm(&params)
	list := roledao.New(ctx.SessionGetSchema()).List(roledao.ListParam{})
	for _, r := range list {
		r.Extend.PutAll(map[string]interface{}{
			"users": userdao.New(ctx.SessionGetSchema()).List(userdao.ListParam{RoleId: r.Id}),
		})
	}
	ctx.JsonSuccess(list)
}
