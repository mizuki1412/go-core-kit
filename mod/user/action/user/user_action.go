package user

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/cryptokit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/userdao"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/sqlkit/pghelper"
	"strings"
)

type loginByUsernameParam struct {
	Username string `description:"用户名" validate:"required"`
	Pwd      string `validate:"required"`
	Schema   string `default:"public"`
}

func loginByUsername(ctx *context.Context) {
	//session := ctx.RenewSession()
	params := loginByUsernameParam{}
	ctx.BindForm(&params)
	if !pghelper.CheckSchemaExist(params.Schema) {
		panic(exception.New("schema不存在"))
	}
	params.Username = strings.TrimSpace(params.Username)
	params.Pwd = cryptokit.MD5(params.Pwd)
	user := userdao.New(params.Schema).Login(params.Pwd, params.Username, "")
	if user == nil {
		panic(exception.New("用户名或密码错误"))
	}
	if user.Off.Int32 == model.UserOffFreeze {
		panic(exception.New("账户被冻结"))
	}
	ctx.SessionSetUser(user)
	ctx.SessionSetSchema(params.Schema)
	ctx.SessionSetToken(ctx.SessionID())
	ctx.SessionSave()
	ret := map[string]interface{}{
		"user":  user,
		"token": ctx.SessionID(),
	}
	if AdditionLoginFunc != nil {
		AdditionLoginFunc(ctx, ret)
	}
	// todo usercenter?
	ctx.JsonSuccess(ret)
}

type loginParam struct {
	Username string `description:"用户名"`
	Phone    string `description:"手机号"`
	Pwd      string `validate:"required"`
	Schema   string `default:"public"`
}

/// 通用登录
func login(ctx *context.Context) {
	//session := ctx.RenewSession()
	params := loginParam{}
	ctx.BindForm(&params)
	if stringkit.IsNull(params.Username) && stringkit.IsNull(params.Phone) {
		panic(exception.New("用户名或手机号缺失"))
	}
	if !pghelper.CheckSchemaExist(params.Schema) {
		panic(exception.New("schema不存在"))
	}
	params.Username = strings.TrimSpace(params.Username)
	params.Phone = strings.TrimSpace(params.Phone)
	params.Pwd = cryptokit.MD5(params.Pwd)
	user := userdao.New(params.Schema).Login(params.Pwd, params.Username, params.Phone)
	if user == nil {
		panic(exception.New("账号和密码不匹配"))
	}
	if user.Off.Int32 == model.UserOffFreeze {
		panic(exception.New("账户被冻结"))
	}
	ctx.SessionSetUser(user)
	ctx.SessionSetSchema(params.Schema)
	ctx.SessionSetToken(ctx.SessionID())
	ctx.SessionSave()
	ret := map[string]interface{}{
		"user":  user,
		"token": ctx.SessionID(),
	}
	if AdditionLoginFunc != nil {
		AdditionLoginFunc(ctx, ret)
	}
	// todo usercenter?
	ctx.JsonSuccess(ret)
}

var AdditionLoginFunc func(ctx *context.Context, ret map[string]interface{})

var AdditionUserExFunc func(ctx *context.Context, u *model.User)

var AdditionUserInfoWithIdFunc = func(ctx *context.Context, u *model.User) {
	// 默认不支持普通用户获取其他用户信息
	panic(exception.New("无权限获取用户信息"))
}

type infoParam struct {
	Id class.Int32 `description:"不填获取自己，并且返回的是user和token；否则只返回user"`
}

func info(ctx *context.Context) {
	params := infoParam{}
	ctx.BindForm(&params)
	if !params.Id.Valid {
		// 获取自己的
		// todo 先走数据库
		user := ctx.SessionGetUser()
		user = userdao.New(ctx.SessionGetSchema()).FindById(user.Id)
		// todo user不存在时
		if AdditionUserExFunc != nil {
			AdditionUserExFunc(ctx, user)
		}
		ctx.JsonSuccess(map[string]interface{}{
			"user":  user,
			"token": ctx.SessionGetToken(),
		})
	} else {
		user := userdao.New(ctx.SessionGetSchema()).FindById(params.Id.Int32)
		// todo user不存在时
		if AdditionUserExFunc != nil {
			AdditionUserExFunc(ctx, user)
		}
		AdditionUserInfoWithIdFunc(ctx, user)
		ctx.JsonSuccess(user)
	}
}

func logout(ctx *context.Context) {
	ctx.SessionClear()
	ctx.JsonSuccess(nil)
}

type updatePwdParam struct {
	OldPwd string `validate:"required"`
	NewPwd string `validate:"required"`
}

func updatePwd(ctx *context.Context) {
	params := updatePwdParam{}
	ctx.BindForm(&params)
	u := ctx.SessionGetUser()
	usermapper := userdao.New(ctx.SessionGetSchema())
	user := usermapper.FindById(u.Id)
	if user == nil {
		panic(exception.New("用户不存在"))
	}
	if user.Pwd.String != cryptokit.MD5(params.OldPwd) {
		panic(exception.New("原密码错误"))
	}
	user.Pwd.Set(cryptokit.MD5(params.NewPwd))
	usermapper.Update(user)
	ctx.SessionSetUser(user)
	ctx.SessionSave()
	ctx.JsonSuccess(nil)
}

type updateUserInfoParam struct {
	Username   class.String
	Name       class.String
	Phone      class.String
	Gender     int8
	Image      class.String
	Address    class.String
	OldPwd     class.String
	NewPwd     class.String
	ExtendJson class.MapString
}

func updateUserInfo(ctx *context.Context) {
	u := ctx.SessionGetUser()
	params := updateUserInfoParam{}
	ctx.BindForm(&params)
	dao := userdao.New(ctx.SessionGetSchema())
	dao.SetResultType(userdao.ResultNone)
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
	if params.Image.Valid {
		u.Image.Set(params.Image)
	}
	if params.Name.Valid {
		u.Name.Set(params.Name.String)
	}
	if params.Phone.Valid {
		u.Phone.Set(params.Phone.String)
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
	if params.OldPwd.Valid && params.NewPwd.Valid && params.OldPwd.String != "" && params.NewPwd.String != "" {
		user := dao.FindById(u.Id)
		if user.Pwd.String != cryptokit.MD5(params.OldPwd.String) {
			panic(exception.New("原密码错误"))
		}
		user.Pwd.Set(cryptokit.MD5(params.NewPwd.String))
	}
	// todo usercenter
	dao.Update(u)
	ctx.JsonSuccess(nil)
}
