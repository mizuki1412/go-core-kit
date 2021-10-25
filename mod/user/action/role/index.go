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
		r.Post("/privilege/list", listAllPrivileges).Swagger.Tag(tag).Summary("所有权限列表")
		r.Post("/list", listRoles).Swagger.Tag(tag).Summary("role列表").Param(listRolesParam{})
		r.Post("/create", create).Swagger.Tag(tag).Summary("role新增").Param(createParams{})
		r.Post("/update", update).Swagger.Tag(tag).Summary("role修改").Param(updateParams{})
		r.Post("/del", del).Swagger.Tag(tag).Summary("role删除").Param(delParams{})
		r.Post("/listRolesWithUser", listRolesWithUser).Swagger.Tag(tag).Summary("列出所有角色，附带所属users").Param(listByRoleParams{})
		r.Post("/department/create", departmentCreate).Swagger.Tag(tag).Summary("部门新增").Param(departmentCreateParams{})
		r.Post("/department/update", departmentUpdate).Swagger.Tag(tag).Summary("部门修改").Param(departmentUpdateParams{})
		r.Post("/department/del", departmentDel).Swagger.Tag(tag).Summary("部门删除").Param(delParams{})
		r.Post("/department/list", listDepartment).Swagger.Tag(tag).Summary("部门列表")
	}
}
