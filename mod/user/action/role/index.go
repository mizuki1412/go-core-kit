package role

import (
	"github.com/mizuki1412/go-core-kit/mod/middleware"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

func Init(router *router.Router) {
	tag := "role:用户模块-角色部门管理"
	r := router.Group("/rest/role")
	r.Use(middleware.AuthUsernameAndPwd())
	{
		r.Post("/privilege/list", listAllPrivileges).Openapi.Tag(tag).Summary("所有权限列表")
		r.Post("/list", listRoles).Openapi.Tag(tag).Summary("role列表").ReqParam(listRolesParam{})
		r.Post("/create", create).Openapi.Tag(tag).Summary("role新增").ReqParam(createParams{})
		r.Post("/update", update).Openapi.Tag(tag).Summary("role修改").ReqParam(updateParams{})
		r.Post("/del", del).Openapi.Tag(tag).Summary("role删除").ReqParam(delParams{})
		r.Post("/listRolesWithUser", listRolesWithUser).Openapi.Tag(tag).Summary("列出所有角色，附带所属users").ReqParam(listByRoleParams{})
		r.Post("/department/create", departmentCreate).Openapi.Tag(tag).Summary("部门新增").ReqParam(departmentCreateParams{})
		r.Post("/department/update", departmentUpdate).Openapi.Tag(tag).Summary("部门修改").ReqParam(departmentUpdateParams{})
		r.Post("/department/del", departmentDel).Openapi.Tag(tag).Summary("部门删除").ReqParam(delParams{})
		r.Post("/department/list", listDepartment).Openapi.Tag(tag).Summary("部门列表")
	}
}
