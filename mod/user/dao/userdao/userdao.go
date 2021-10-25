package userdao

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/roledao"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
	"github.com/mizuki1412/go-core-kit/service/sqlkit/pghelper"
	"github.com/spf13/cast"
)

/// auto template
type Dao struct {
	sqlkit.Dao
}

const (
	ResultDefault byte = iota
	ResultNone
)

func New(schema string, tx ...*sqlkit.Dao) *Dao {
	dao := &Dao{}
	dao.NewHelper(schema, tx...)
	return dao
}
func (dao *Dao) cascade(obj *model.User) {
	switch dao.ResultType {
	case ResultDefault:
		if obj.Role != nil {
			obj.Role = roledao.New(dao.Schema).FindById(obj.Role.Id)
		}
	case ResultNone:
		obj.Role = nil
	}
}
func (dao *Dao) scan(sql string, args []interface{}) []*model.User {
	rows := dao.Query(sql, args...)
	list := make([]*model.User, 0, 5)
	defer rows.Close()
	for rows.Next() {
		m := &model.User{}
		err := rows.StructScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, m)
	}
	for i := range list {
		dao.cascade(list[i])
	}
	return list
}
func (dao *Dao) scanOne(sql string, args []interface{}) *model.User {
	rows := dao.Query(sql, args...)
	defer rows.Close()
	for rows.Next() {
		m := model.User{}
		err := rows.StructScan(&m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		dao.cascade(&m)
		return &m
	}
	return nil
}

////

func (dao *Dao) Login(pwd, username, phone string) *model.User {
	builder := sqlkit.Builder().Select("*").From(dao.GetTableD("admin_user"))
	if !stringkit.IsNull(username) {
		builder = builder.Where("username=?", username)
	} else {
		builder = builder.Where("phone=?", phone)
	}
	sql, args := builder.Where("pwd=?", pwd).Where("off>=0").Limit(1).MustSql()
	return dao.scanOne(sql, args)
}
func (dao *Dao) FindById(id int32) *model.User {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("admin_user")).Where("id=?", id).MustSql()
	return dao.scanOne(sql, args)
}
func (dao *Dao) FindByPhone(phone string) *model.User {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("admin_user")).Where("phone=?", phone).Where("off>=0").MustSql()
	return dao.scanOne(sql, args)
}

func (dao *Dao) FindByUsername(username string) *model.User {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("admin_user")).Where("username=?", username).Where("off>=0").MustSql()
	return dao.scanOne(sql, args)
}

// FindParam 可以通过extend的值来find
type FindParam struct {
	Extend map[string]interface{}
}

func (dao *Dao) Find(param FindParam) *model.User {
	builder := sqlkit.Builder().Select("*").From(dao.GetTableD("admin_user")).Where("off>=0").Limit(1)
	for k, v := range param.Extend {
		builder = builder.Where(fmt.Sprintf("extend->>'%s'=?", k), cast.ToString(v))
	}
	sql, args := builder.MustSql()
	return dao.scanOne(sql, args)
}

func (dao *Dao) ListFromRootDepart(departId int) []*model.User {
	where := fmt.Sprintf(`
off>-1 and role>0 and role in 
  ( select id from %s where department in 
     (with recursive t(id) as( values(%d) union all select d.id from %s d, t where t.id=d.parent) select id from t )
  )`, dao.GetTable(&model.Role{}), departId, dao.GetTable(&model.Department{}))
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("admin_user")).Where(where).OrderBy("name, id").MustSql()
	return dao.scan(sql, args)
}

type ListParam struct {
	RoleId      int32
	Roles       []int32
	Departments []int32
	IdList      []int32
}

func (dao *Dao) List(param ListParam) []*model.User {
	builder := sqlkit.Builder().Select("*").From(dao.GetTableD("admin_user")).Where("off>-1").OrderBy("name, id")
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
	return dao.scan(sql, args)
}

func (dao *Dao) OffUser(uid int32, off int32) {
	sql, args, err := sqlkit.Builder().Update(dao.GetTableD("admin_user")).Set("off", off).Where("id=?", uid).ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	dao.Exec(sql, args...)
}
func (dao *Dao) SetNull(id int32) {
	sql, args, err := sqlkit.Builder().Update(dao.GetTableD("admin_user")).Set("phone", squirrel.Expr("null")).Where("id=?", id).ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	dao.Exec(sql, args...)
}
