package user

import (
	"github.com/mizuki1412/go-core-kit/mod/user/action/role"
	"github.com/mizuki1412/go-core-kit/mod/user/action/smscode"
	"github.com/mizuki1412/go-core-kit/mod/user/action/user"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

// All 用户、部门、角色模块
func All() []func(r *router.Router) {
	return []func(r *router.Router){user.Init, role.Init, smscode.Init}
}
