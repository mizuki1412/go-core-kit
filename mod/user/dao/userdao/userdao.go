package userdao

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/departmentdao"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/roledao"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
	"github.com/mizuki1412/go-core-kit/service/sqlkit/pghelper"
	"github.com/spf13/cast"
)

type Dao struct {
	sqlkit.Dao[model.User]
}

var meta = sqlkit.InitModelMeta(&model.User{})

const (
	ResultDefault byte = iota
	ResultNone
)

func New(ds ...*sqlkit.DataSource) *Dao {
	dao := &Dao{}
	if len(ds) > 0 {
		dao.DataSource = ds[0]
	}
	dao.Cascade = func(obj *model.User) {
		switch dao.ResultType {
		case ResultDefault:
			if obj.Role != nil {
				obj.Role = roledao.New(dao.DataSource).FindById(obj.Role.Id)
			}
			if obj.Department != nil {
				obj.Department = departmentdao.New(dao.DataSource).FindById(obj.Department.Id)
			}
		case ResultNone:
			obj.Role = nil
			obj.Department = nil
		}
	}
	return dao
}

func (dao *Dao) Login(pwd, username, phone string) *model.User {
	builder := dao.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema))
	if !stringkit.IsNull(username) {
		builder = builder.Where("username=?", username)
	} else {
		builder = builder.Where("phone=?", phone)
	}
	sql, args := builder.Where("pwd=?", pwd).Where("off>=0").Limit(1).MustSql()
	return dao.ScanOne(sql, args)
}
func (dao *Dao) FindById(id int32) *model.User {
	sql, args := dao.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where("id=?", id).MustSql()
	return dao.ScanOne(sql, args)
}
func (dao *Dao) FindByPhone(phone string) *model.User {
	sql, args := dao.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where("phone=?", phone).Where("off>=0").MustSql()
	return dao.ScanOne(sql, args)
}

func (dao *Dao) FindByUsername(username string) *model.User {
	sql, args := dao.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where("username=?", username).Where("off>=0").MustSql()
	return dao.ScanOne(sql, args)
}
func (dao *Dao) FindByUsernameDeleted(username string) *model.User {
	sql, args := dao.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where("username=?", username).Where("off=-1").MustSql()
	return dao.ScanOne(sql, args)
}

// FindParam 可以通过extend的值来find
type FindParam struct {
	Extend map[string]any
}

func (dao *Dao) Find(param FindParam) *model.User {
	builder := dao.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where("off>=0").Limit(1)
	for k, v := range param.Extend {
		builder = builder.Where(fmt.Sprintf("extend->>'%s'=?", k), cast.ToString(v))
	}
	sql, args := builder.MustSql()
	return dao.ScanOne(sql, args)
}

func (dao *Dao) ListFromRootDepart(departId int) []*model.User {
	where := fmt.Sprintf(`
off>-1 and role>0 and role in 
  ( select id from %s where department in 
     (with recursive t(id) as( values(%d) union all select d.id from %s d, t where t.id=d.parent) select id from t )
  )`, dao.GetTable(&model.Role{}), departId, dao.GetTable(&model.Department{}))
	sql, args := dao.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where(where).OrderBy("name, id").MustSql()
	return dao.ScanList(sql, args)
}

type ListParam struct {
	RoleId      int32
	Roles       []int32
	Departments []int32
	IdList      []int32
}

func (dao *Dao) List(param ListParam) []*model.User {
	builder := dao.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where("off>-1").OrderBy("id")
	if param.RoleId != 0 {
		builder = builder.Where("role=?", param.RoleId)
	}
	if len(param.IdList) > 0 {
		builder = pghelper.WhereUnnestInt(builder, "id in ", param.IdList)
	}
	if len(param.Roles) > 0 {
		builder = pghelper.WhereUnnestInt(builder, "role in ", param.Roles)
	}
	// 根据role组筛选
	if len(param.Departments) > 0 {
		flag, arg := pghelper.GenUnnestInt(param.Departments)
		builder = builder.Where(fmt.Sprintf("role in (select id from role where department in %s)", flag), arg...)
	}
	sql, args := builder.MustSql()
	return dao.ScanList(sql, args)
}

func (dao *Dao) OffUser(uid int32, off int32) {
	sql, args, err := dao.Builder().Update(meta.GetTableName(dao.Schema)).Set("off", off).Where("id=?", uid).ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	dao.Exec(sql, args...)
}
func (dao *Dao) SetNull(id int32) {
	sql, args, err := dao.Builder().Update(meta.GetTableName(dao.Schema)).Set("phone", squirrel.Expr("null")).Set("username", squirrel.Expr("null")).Where("id=?", id).ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	dao.Exec(sql, args...)
}
