package user

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/cryptokit"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/roledao"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/userdao"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"time"
)

var AdditionUserExAdminFunc func(ctx *context.Context, u *model.User)

type listUsersParams struct {
	DepartmentId int
}

func listUsers(ctx *context.Context) {
	params := listUsersParams{}
	ctx.BindForm(&params)
	list := userdao.New(ctx.SessionGetSchema()).ListFromRootDepart(params.DepartmentId)
	if AdditionUserExAdminFunc != nil {
		for _, u := range list {
			AdditionUserExAdminFunc(ctx, u)
		}
	}
	ctx.JsonSuccess(list)
}

type listByRoleParams struct {
	RoleId int32 `validate:"required"`
}

func listByRole(ctx *context.Context) {
	params := listByRoleParams{}
	ctx.BindForm(&params)
	list := userdao.New(ctx.SessionGetSchema()).List(userdao.ListParam{RoleId: params.RoleId})
	if AdditionUserExAdminFunc != nil {
		for _, u := range list {
			AdditionUserExAdminFunc(ctx, u)
		}
	}
	ctx.JsonSuccess(list)
}

type AddUserParams struct {
	Username   class.String `validate:"required"`
	Pwd        class.String `validate:"required"`
	Role       int32        `validate:"required"`
	Name       class.String
	Phone      class.String
	Gender     int8
	Image      class.String
	Address    class.String
	ExtendJson class.MapString
}

func AddUser(ctx *context.Context) {
	params := AddUserParams{}
	ctx.BindForm(&params)
	u := AddUserHandle(ctx, params)
	ctx.JsonSuccess(u)
}

func AddUserHandle(ctx *context.Context, params AddUserParams) model.User {
	dao := userdao.New(ctx.SessionGetSchema())
	dao.SetResultType(userdao.ResultNone)
	if dao.FindByUsername(params.Username.String) != nil {
		panic(exception.New("用户名已经存在"))
	}
	if params.Phone.Valid && dao.FindByPhone(params.Phone.String) != nil {
		panic(exception.New("手机号已经存在"))
	}
	roleDao := roledao.New(ctx.SessionGetSchema())
	roleDao.SetResultType(roledao.ResultNone)
	r := roleDao.FindById(params.Role)
	if r == nil {
		panic(exception.New("角色不存在"))
	}
	u := model.User{}
	if params.Username.Valid {
		u.Username.Set(params.Username)
	}
	if params.Pwd.Valid {
		u.Pwd.Set(cryptokit.MD5(params.Pwd.String))
	}
	u.Role = r
	if params.Name.Valid {
		u.Name.Set(params.Name)
	}
	if params.Phone.Valid {
		u.Phone.Set(params.Phone)
	}
	if params.Image.Valid {
		u.Image.Set(params.Image)
	}
	if params.Gender != 0 {
		u.Gender.Set(params.Gender)
	}
	if params.Address.Valid {
		u.Address.Set(params.Address)
	}
	if params.ExtendJson.Valid {
		u.Extend.PutAll(params.ExtendJson.Map)
	}
	u.CreateDt.Set(time.Now())
	dao.Insert(&u)
	return u
}

type UpdateParams struct {
	Id         int32 `validate:"required"`
	Username   class.String
	Name       class.String
	Phone      class.String
	Gender     int8
	Image      class.String
	Address    class.String
	Pwd        class.String
	Role       int32
	ExtendJson class.MapString
}

func UpdateUser(ctx *context.Context) {
	params := UpdateParams{}
	ctx.BindForm(&params)
	UpdateUserHandle(ctx, params)
	ctx.JsonSuccess(nil)
}

func UpdateUserHandle(ctx *context.Context, params UpdateParams) {
	dao := userdao.New(ctx.SessionGetSchema())
	u := dao.FindById(params.Id)
	if u == nil || u.Off.Int32 == model.UserOffDelete {
		panic(exception.New("用户不存在"))
	}
	if u.Role != nil && u.Role.Id == 0 {
		panic(exception.New("该用户不能设置"))
	}
	if params.Phone.Valid && params.Phone.String != "" && params.Phone.String != u.Phone.String && dao.FindByPhone(params.Phone.String) != nil {
		panic(exception.New("手机号已存在"))
	}
	if params.Username.Valid && params.Username.String != u.Username.String {
		if dao.FindByUsername(params.Username.String) != nil {
			panic(exception.New("该用户名已被使用"))
		} else {
			u.Username.Set(params.Username.String)
		}
	}
	if params.Role > 0 && (u.Role == nil || params.Role != u.Role.Id) {
		r := roledao.New(ctx.SessionGetSchema()).FindById(params.Role)
		if r == nil {
			panic(exception.New("role不存在"))
		}
		u.Role = r
	}
	if params.Name.Valid {
		u.Name.Set(params.Name.String)
	}
	if params.Phone.Valid {
		u.Phone.Set(params.Phone.String)
	}
	if params.Image.Valid {
		u.Image.Set(params.Image)
	}
	if params.Pwd.Valid && params.Pwd.String != "" {
		u.Pwd.Set(cryptokit.MD5(params.Pwd.String))
	}
	if params.Gender != 0 {
		u.Gender.Set(params.Gender)
	}
	if params.Address.Valid {
		u.Address.Set(params.Address.String)
	}
	if params.ExtendJson.Valid {
		u.Extend.PutAll(params.ExtendJson.Map)
	}
	dao.Update(u)
}

type infoAdminParams struct {
	Uid int32 `validate:"required"`
}

func infoAdmin(ctx *context.Context) {
	params := infoAdminParams{}
	ctx.BindForm(&params)
	user := userdao.New(ctx.SessionGetSchema()).FindById(params.Uid)
	if user == nil {
		panic(exception.New("无此用户"))
	}
	ctx.JsonSuccess(user)
}

type DelParams struct {
	Id  int32       `validate:"required"`
	Off class.Int32 `validate:"required" description:"0-删除，1-冻结，2-解冻"`
}

func DelUser(ctx *context.Context) {
	params := DelParams{}
	ctx.BindForm(&params)
	mine := ctx.SessionGetUser()
	if mine == nil {
		panic(exception.New("登录的用户错误"))
	}
	if mine.Id == params.Id {
		panic(exception.New("不能操作自己"))
	}
	dao := userdao.New(ctx.SessionGetSchema(), ctx.DBTx())
	target := userdao.New(ctx.SessionGetSchema()).FindById(params.Id)
	if target == nil {
		panic(exception.New("用户不存在"))
	}
	if target.Role != nil && target.Role.Id == 0 {
		panic(exception.New("该用户不能设置"))
	}
	//
	if params.Off.Int32 == 0 {
		dao.OffUser(params.Id, model.UserOffDelete)
		dao.SetNull(params.Id)
		//target.setOff(User.OFF_DEL);
		// todo
		//userCenter.add(target);
	} else if params.Off.Int32 == 1 {
		dao.OffUser(params.Id, model.UserOffFreeze)
		//target.setOff(User.OFF_FREEZE);
		//userCenter.add(target);
	} else if params.Off.Int32 == 2 {
		dao.OffUser(params.Id, model.UserOffOK)
		//target.setOff(User.OFF_OK);
		//userCenter.add(target);
	}
	ctx.JsonSuccess(nil)
}
