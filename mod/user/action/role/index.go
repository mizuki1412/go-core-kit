package role

import (
	"github.com/mizuki1412/go-core-kit/v2/mod/user/model"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/middleware"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/openapi"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/router"
)

func Init(router *router.Router) {
	tag := "role:用户模块-角色部门管理"
	r := router.Group("/role").Use(middleware.AuthJWT())
	r.Get("/privilege/list", listAllPrivileges).Api(openapi.Tag(tag), openapi.Summary("所有权限列表"),
		openapi.Response([]*model.PrivilegeConstant{}))
	r.Post("/list", listRoles).Api(openapi.Tag(tag), openapi.Summary("role列表"),
		openapi.ReqParam(listRolesParam{}), openapi.Response([]*model.Role{}))
	r.Post("/create", create).Api(openapi.Tag(tag), openapi.Summary("role新增"), openapi.ReqBody(createParams{}))
	r.Post("/update", update).Api(openapi.Tag(tag), openapi.Summary("role修改"), openapi.ReqBody(updateParams{}))
	r.Get("/del", del).Api(openapi.Tag(tag), openapi.Summary("role删除"), openapi.ReqParam(delParams{}))
	r.Post("/listRolesWithUser", listRolesWithUser).Api(openapi.Tag(tag),
		openapi.Summary("列出所有角色，附带所属users"),
		openapi.ReqParam(listByRoleParams{}), openapi.Response([]*model.Role{}))
	r.Post("/department/create", departmentCreate).Api(openapi.Tag(tag),
		openapi.Summary("部门新增"), openapi.ReqBody(departmentCreateParams{}))
	r.Post("/department/update", departmentUpdate).Api(openapi.Tag(tag),
		openapi.Summary("部门修改"), openapi.ReqBody(departmentUpdateParams{}))
	r.Post("/department/del", departmentDel).Api(openapi.Tag(tag),
		openapi.Summary("部门删除"), openapi.ReqParam(delParams{}))
	r.Post("/department/list", listDepartment).Api(openapi.Tag(tag),
		openapi.Summary("部门列表"), openapi.Response([]*model.Department{}))
}
