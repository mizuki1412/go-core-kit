package role

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/departmentdao"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/roledao"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/userdao"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"time"
)

type departmentCreateParams struct {
	No          class.String
	Name        string `validate:"required"`
	Description class.String
	ParentId    class.Int32
}

func departmentCreate(ctx *context.Context) {
	params := departmentCreateParams{}
	ctx.BindForm(&params)
	department := &model.Department{}
	dao := departmentdao.New()
	dao.DataSource().Schema = ctx.SessionGetSchema()
	dao.ResultType = departmentdao.ResultNone
	if params.ParentId.Valid {
		parent := dao.SelectOneById(params.ParentId.Int32)
		if parent == nil {
			panic(exception.New("父级部门不存在"))
		}
		department.Parent = parent
	}
	department.Name.Set(params.Name)
	if params.No.Valid {
		department.No.Set(params.No.String)
	}
	if params.Description.Valid {
		department.Descr.Set(params.Description.String)
	}
	department.CreateDt.Set(time.Now())
	dao.InsertObj(department)
	ctx.JsonSuccess(nil)
}

type departmentUpdateParams struct {
	Id          int32 `validate:"required"`
	No          class.String
	Name        class.String
	Description class.String
	ParentId    class.Int32
}

func departmentUpdate(ctx *context.Context) {
	params := departmentUpdateParams{}
	ctx.BindForm(&params)
	dao := departmentdao.New()
	dao.DataSource().Schema = ctx.SessionGetSchema()
	department := dao.SelectOneById(params.Id)
	dao.ResultType = departmentdao.ResultNone
	if department == nil {
		panic(exception.New("部门不存在"))
	}
	if params.No.Valid {
		department.No.Set(params.No.String)
	}
	if params.Name.Valid {
		department.Name.Set(params.Name.String)
	}
	if params.Description.Valid {
		department.Descr.Set(params.Description.String)
	}
	if params.ParentId.Valid && (department.Parent == nil || params.ParentId.Int32 != department.Parent.Id) {
		parent := dao.SelectOneById(params.ParentId.Int32)
		if parent == nil {
			panic(exception.New("父级部门不存在"))
		}
		department.Parent = parent
	}
	dao.UpdateObj(department)
	ctx.JsonSuccess(nil)
}
func departmentDel(ctx *context.Context) {
	params := delParams{}
	ctx.BindForm(&params)
	dao := departmentdao.New()
	dao.DataSource().Schema = ctx.SessionGetSchema()
	dao.ResultType = departmentdao.ResultNone
	department := dao.SelectOneById(params.Id)
	if department == nil {
		panic(exception.New("部门不存在"))
	}
	roleDao := roledao.New()
	roleDao.DataSource().Schema = ctx.SessionGetSchema()
	roleDao.ResultType = userdao.ResultNone
	rs := roleDao.ListByDepartment(params.Id)
	if rs != nil && len(rs) > 0 {
		panic(exception.New("部门下还有角色,不能删除"))
	}
	if val, ok := department.Extend.Map["immutable"]; ok && val.(bool) {
		panic(exception.New("该部门不可删除"))
	}
	dao.DeleteById(department.Id)
	ctx.JsonSuccess(nil)
}

func listDepartment(ctx *context.Context) {
	dao := departmentdao.New()
	dao.DataSource().Schema = ctx.SessionGetSchema()
	dao.ResultType = departmentdao.ResultAll
	ctx.JsonSuccess(dao.ListAll())
}
