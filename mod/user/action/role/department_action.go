package role

import (
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/dao/departmentdao"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/dao/roledao"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/dao/userdao"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/model"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
	"time"
)

type departmentCreateParams struct {
	No          class.String `validate:"required"`
	Name        string       `validate:"required"`
	Description class.String
	ParentId    class.Int64
	Extend      class.MapString
}

func departmentCreate(ctx *context.Context) {
	params := departmentCreateParams{}
	ctx.BindForm(&params)
	department := &model.Department{}
	dao := departmentdao.New(departmentdao.ResultNone)
	dao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	if params.ParentId.Valid {
		parent := dao.SelectOneById(params.ParentId.Int64)
		if parent == nil {
			panic(exception.New("父级部门不存在"))
		}
		department.Parent = parent
	}
	department.Name.Set(params.Name)
	if params.No.Valid {
		if dao.FindByNo(params.No.String) != nil {
			panic(exception.New("当前编号已被占用"))
		}
		department.No.Set(params.No.String)
	}
	if params.Description.Valid {
		department.Descr.Set(params.Description.String)
	}
	department.CreateDt.Set(time.Now())
	department.Extend.Set(params.Extend)
	dao.InsertObj(department)
	ctx.JsonSuccess()
}

type departmentUpdateParams struct {
	Id          int64 `validate:"required"`
	No          class.String
	Name        class.String
	Description class.String
	ParentId    class.Int64
	Extend      class.MapString
}

func departmentUpdate(ctx *context.Context) {
	params := departmentUpdateParams{}
	ctx.BindForm(&params)
	dao := departmentdao.New(departmentdao.ResultNone)
	dao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	department := dao.SelectOneById(params.Id)
	if department == nil {
		panic(exception.New("部门不存在"))
	}
	if params.No.Valid && params.No.String != department.No.String {
		if dao.FindByNo(params.No.String) != nil {
			panic(exception.New("当前编号已被占用"))
		}
		department.No.Set(params.No.String)
	}
	if params.Name.Valid {
		department.Name.Set(params.Name.String)
	}
	if params.Description.Valid {
		department.Descr.Set(params.Description.String)
	}
	if params.ParentId.Valid && (department.Parent == nil || params.ParentId.Int64 != department.Parent.Id) {
		parent := dao.SelectOneById(params.ParentId.Int64)
		if parent == nil {
			panic(exception.New("父级部门不存在"))
		}
		department.Parent = parent
	}
	if params.Extend.Valid {
		department.Extend.PutAll(params.Extend.Map)
	}
	dao.UpdateObj(department)
	ctx.JsonSuccess()
}
func departmentDel(ctx *context.Context) {
	params := delParams{}
	ctx.BindForm(&params)
	dao := departmentdao.New(departmentdao.ResultNone)
	dao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	department := dao.SelectOneById(params.Id)
	if department == nil {
		panic(exception.New("部门不存在"))
	}
	if val, ok := department.Extend.Map["immutable"]; ok && val.(bool) {
		panic(exception.New("该部门不可删除"))
	}
	// 判断是否有角色
	roleDao := roledao.New(userdao.ResultNone)
	roleDao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	rNum := roleDao.CountFromRootDepart(department.Id)
	if rNum > 0 {
		panic(exception.New("部门下还有角色,不能删除"))
	}
	// 判断是否有用户
	userDao := userdao.New(userdao.ResultNone)
	userDao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	uNum := userDao.CountFromRootDepart(department.Id)
	if uNum > 0 {
		panic(exception.New("部门下还有用户,不能删除"))
	}

	dao.DeleteById(department.Id)
	ctx.JsonSuccess()
}

func listDepartment(ctx *context.Context) {
	dao := departmentdao.New(departmentdao.ResultAll)
	dao.DataSource().Schema = ctx.GetJwt().Ext.GetString("schema")
	ctx.JsonSuccess(dao.ListAll())
}
