package role

import (
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/dao/departmentdao"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/dao/privilegedao"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/dao/roledao"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/dao/userdao"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/model"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
	"time"
)

func listAllPrivileges(ctx *context.Context) {
	dao := privilegedao.New()
	dao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	ctx.JsonSuccess(dao.ListPrivileges())
}

type createParams struct {
	Name           string          `validate:"required"`
	PrivilegesJson class.ArrString `validate:"required" default:"[]" comment:"[a,b,c]"`
	DepartmentId   class.Int64
	Extend         class.MapString
}

func create(ctx *context.Context) {
	params := createParams{}
	ctx.BindForm(&params)
	role := &model.Role{}
	if params.DepartmentId.Valid {
		departmentDao := departmentdao.New(departmentdao.ResultNone)
		departmentDao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
		department := departmentDao.SelectOneById(params.DepartmentId)
		if department == nil {
			panic(exception.New("部门不存在"))
		}
		role.Department = department
	}
	role.Name.Set(params.Name)
	role.Privileges = params.PrivilegesJson
	role.CreateDt.Set(time.Now())
	role.Extend.Set(params.Extend)
	rdao := roledao.New(roledao.ResultDefault)
	rdao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	rdao.InsertObj(role)
	ctx.JsonSuccess()
}

type updateParams struct {
	Id             int64 `validate:"required"`
	Name           class.String
	PrivilegesJson class.ArrString `comment:"数组json字符串：[a,b,c]"`
	DepartmentId   class.Int64
	Extend         class.MapString
}

func update(ctx *context.Context) {
	params := updateParams{}
	ctx.BindForm(&params)
	dao := roledao.New(roledao.ResultDefault)
	dao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	role := dao.SelectOneById(params.Id)
	if role == nil {
		panic(exception.New("角色不存在"))
	}
	if params.DepartmentId.Valid && (role.Department == nil || params.DepartmentId.Int64 != role.Department.Id) {
		departmentDao := departmentdao.New(departmentdao.ResultNone)
		departmentDao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
		d := departmentDao.SelectOneById(params.DepartmentId.Int64)
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
	if params.Extend.IsValid() {
		role.Extend.PutAll(params.Extend.Map)
	}
	dao.UpdateObj(role)
	ctx.JsonSuccess()
}

type delParams struct {
	Id int64 `validate:"required"`
}

func del(ctx *context.Context) {
	params := delParams{}
	ctx.BindForm(&params)
	dao := roledao.New(roledao.ResultNone)
	dao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	role := dao.SelectOneById(params.Id)
	if role == nil {
		panic(exception.New("角色不存在"))
	}
	if val, ok := role.Extend.Map["immutable"]; ok && val.(bool) {
		panic(exception.New("该角色不可删除"))
	}
	userDao := userdao.New(userdao.ResultNone)
	userDao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	rid := class.Int64{}
	rid.Set(params.Id)
	us := userDao.List(userdao.ListParam{RoleId: rid})
	if us != nil && len(us) > 0 {
		panic(exception.New("角色下还有用户,不能删除"))
	}
	dao.DeleteById(role.Id)
	ctx.JsonSuccess()
}

type listRolesParam struct {
	Root class.Int64 `comment:"指定根department"`
}

func listRoles(ctx *context.Context) {
	params := listRolesParam{}
	ctx.BindForm(&params)
	var roles []*model.Role
	dao := roledao.New(roledao.ResultDefault)
	dao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	if params.Root.Valid {
		roles = dao.ListFromRootDepart(params.Root.Int64)
	} else {
		roles = dao.List(roledao.ListParam{})
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
	RoleId int64 `validate:"required"`
}

func listRolesWithUser(ctx *context.Context) {
	params := listByRoleParams{}
	ctx.BindForm(&params)
	dao := roledao.New(roledao.ResultDefault)
	dao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	list := dao.List(roledao.ListParam{})
	udao := userdao.New(userdao.ResultDefault)
	udao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	for _, r := range list {
		r.Extend.PutAll(map[string]any{
			"users": udao.List(userdao.ListParam{Roles: []int64{params.RoleId}}),
		})
	}
	ctx.JsonSuccess(list)
}
