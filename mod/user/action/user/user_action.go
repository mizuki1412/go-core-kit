package user

import (
	context2 "context"
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/cryptokit"
	"github.com/mizuki1412/go-core-kit/v2/library/stringkit"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/dao/userdao"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/model"
	"github.com/mizuki1412/go-core-kit/v2/service/jwtkit"
	"github.com/mizuki1412/go-core-kit/v2/service/rediskit"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/v2/service/sqlkit/pghelper"
	"strings"
)

type loginByUsernameParam struct {
	Username string `comment:"用户名" validate:"required"`
	Pwd      string `validate:"required"`
	Schema   string `default:"public"`
}

type ResLogin struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

func loginByUsername(ctx *context.Context) {
	//session := ctx.RenewSession()
	params := loginByUsernameParam{}
	ctx.BindForm(&params)
	//if !pghelper.CheckSchemaExist(params.Schema) {
	//	panic(exception.New("schema不存在"))
	//}
	params.Username = strings.TrimSpace(params.Username)
	params.Pwd = cryptokit.MD5(params.Pwd)
	dao := userdao.New(userdao.ResultDefault)
	dao.DataSource().Schema = params.Schema
	user := dao.Login(params.Pwd, params.Username, "")
	if user == nil {
		panic(exception.New("用户名或密码错误"))
	}
	if user.Off.Int32 == model.UserOffFreeze {
		panic(exception.New("账户被冻结"))
	}
	claim := jwtkit.New(user.Id)
	claim.Ext.Put("schema", params.Schema)
	token := claim.Token()
	ctx.SetJwtCookie(claim, token)
	ret := ResLogin{
		User:  user,
		Token: token,
	}
	if AdditionLoginFunc != nil {
		AdditionLoginFunc(ctx, ret)
	}
	ctx.JsonSuccess(ret)
}

type loginParam struct {
	Username string `comment:"用户名"`
	Phone    string `comment:"手机号"`
	Pwd      string `validate:"required"`
	Schema   string `default:"public"`
}

// / 通用登录
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
	dao := userdao.New(userdao.ResultDefault)
	dao.DataSource().Schema = params.Schema
	user := dao.Login(params.Pwd, params.Username, params.Phone)
	if user == nil {
		panic(exception.New("账号和密码不匹配"))
	}
	if user.Off.Int32 == model.UserOffFreeze {
		panic(exception.New("账户被冻结"))
	}
	claim := jwtkit.New(user.Id)
	claim.Ext.Put("schema", params.Schema)
	token := claim.Token()
	ret := ResLogin{
		User:  user,
		Token: token,
	}
	ctx.SetJwtCookie(claim, token)
	if AdditionLoginFunc != nil {
		AdditionLoginFunc(ctx, ret)
	}
	ctx.JsonSuccess(ret)
}

var AdditionLoginFunc func(ctx *context.Context, ret ResLogin)

var AdditionUserExFunc func(ctx *context.Context, u *model.User)

var AdditionUserInfoWithIdFunc = func(ctx *context.Context, u *model.User) {
	// 默认不支持普通用户获取其他用户信息
	panic(exception.New("无权限获取用户信息"))
}

type infoParam struct {
	Id     class.Int32  `comment:"不填获取自己，并且返回的是user和token；否则只返回user"`
	Schema class.String `comment:"用于校验当前登录的和需要的是不是一个schema"`
}

func info(ctx *context.Context) {
	params := infoParam{}
	ctx.BindForm(&params)
	dao := userdao.New(userdao.ResultDefault)
	schema0 := ctx.GetJwt().Ext.GetString("schema")
	dao.DataSource().Schema = schema0
	if !params.Id.Valid {
		if params.Schema.String != "" && params.Schema.String != schema0 {
			ctx.Json(context.RestRet{
				Result:  context.ResultAuthErr,
				Message: "schema不匹配",
			})
			return
		}
		// 获取自己的
		// todo 先走数据库
		uid := ctx.GetJwt().IdInt32()
		user := dao.SelectOneById(uid)
		claim := jwtkit.New(user.Id)
		claim.Ext.Put("schema", schema0)
		token := claim.Token()
		ret := ResLogin{
			User:  user,
			Token: token,
		}
		ctx.SetJwtCookie(claim, token)
		// todo user不存在时
		if AdditionUserExFunc != nil {
			AdditionUserExFunc(ctx, user)
		}
		ctx.JsonSuccess(ret)
	} else {
		user := dao.SelectOneById(params.Id.Int32)
		// todo user不存在时
		if AdditionUserExFunc != nil {
			AdditionUserExFunc(ctx, user)
		}
		AdditionUserInfoWithIdFunc(ctx, user)
		ctx.JsonSuccess(user)
	}
}

func logout(ctx *context.Context) {
	ctx.DestroyJwt()
	ctx.JsonSuccess()
}

type updatePwdParam struct {
	OldPwd string `validate:"required"`
	NewPwd string `validate:"required"`
}

func updatePwd(ctx *context.Context) {
	params := updatePwdParam{}
	ctx.BindForm(&params)
	uid := ctx.GetJwt().IdInt32()
	dao := userdao.New(userdao.ResultDefault)
	dao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	user := dao.SelectOneById(uid)
	if user == nil {
		panic(exception.New("用户不存在"))
	}
	if user.Pwd.String != cryptokit.MD5(params.OldPwd) {
		panic(exception.New("原密码错误"))
	}
	user.Pwd.Set(cryptokit.MD5(params.NewPwd))
	dao.UpdateObj(user)
	ctx.JsonSuccess()
}

type updateUserInfoParam struct {
	Username   class.String
	Name       class.String
	Phone      class.String
	Sms        class.String
	Gender     int8
	Image      class.String
	Address    class.String
	OldPwd     class.String
	NewPwd     class.String
	ExtendJson class.MapString
}

func updateUserInfo(ctx *context.Context) {
	uid := ctx.GetJwt().IdInt32()
	params := updateUserInfoParam{}
	ctx.BindForm(&params)
	dao := userdao.New(userdao.ResultNone)
	dao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	u := dao.SelectOneWithDelById(uid)
	if params.Phone.Valid && params.Phone.String != "" && params.Phone.String != u.Phone.String {
		if dao.FindByPhone(params.Phone.String) != nil {
			panic(exception.New("手机号已被注册"))
		}
		if !params.Sms.Valid || rediskit.Get(context2.Background(), rediskit.GetKeyWithPrefix("sms:"+params.Phone.String), "") != params.Sms.String {
			panic(exception.New("验证码错误"))
		}
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
		user := dao.SelectOneById(u.Id)
		if user.Pwd.String != cryptokit.MD5(params.OldPwd.String) {
			panic(exception.New("原密码错误"))
		}
		user.Pwd.Set(cryptokit.MD5(params.NewPwd.String))
	}
	//todo usercenter
	dao.UpdateObj(u)
	ctx.JsonSuccess()
}
