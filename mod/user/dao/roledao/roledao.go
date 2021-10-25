package roledao

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/departmentdao"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
	"github.com/mizuki1412/go-core-kit/service/sqlkit/pghelper"
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
func (dao *Dao) cascade(obj *model.Role) {
	switch dao.ResultType {
	case ResultDefault:
		if obj.Department != nil {
			obj.Department = departmentdao.New(dao.Schema).FindById(obj.Department.Id)
		}
	case ResultNone:
		obj.Department = nil
	}
}
func (dao *Dao) scanPrivilege(sql string, args []interface{}) []*model.PrivilegeConstant {
	rows := dao.Query(sql, args...)
	list := make([]*model.PrivilegeConstant, 0, 5)
	defer rows.Close()
	for rows.Next() {
		m := &model.PrivilegeConstant{}
		err := rows.StructScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, m)
	}
	return list
}
func (dao *Dao) scan(sql string, args []interface{}) []*model.Role {
	rows := dao.Query(sql, args...)
	list := make([]*model.Role, 0, 5)
	defer rows.Close()
	for rows.Next() {
		m := &model.Role{}
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
func (dao *Dao) scanOne(sql string, args []interface{}) *model.Role {
	rows := dao.Query(sql, args...)
	defer rows.Close()
	for rows.Next() {
		m := model.Role{}
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

func (dao *Dao) FindById(id int32) *model.Role {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("role")).Where("id=?", id).MustSql()
	return dao.scanOne(sql, args)
}
func (dao *Dao) FindByName(name string) *model.Role {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("role")).Where("name=?", name).Limit(1).MustSql()
	return dao.scanOne(sql, args)
}
func (dao *Dao) ListFromRootDepart(id int32) []*model.Role {
	where := fmt.Sprintf(`id>0 and department in ( with recursive t(id) as( values(%d) union all select d.id from %s d, t where t.id=d.parent) select id from t )`, id, dao.GetTable(&model.Department{}))
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("role")).Where(where).OrderBy("id").MustSql()
	return dao.scan(sql, args)
}

type ListParam struct {
	Departments []int32
}

func (dao *Dao) List(param ListParam) []*model.Role {
	builder := sqlkit.Builder().Select("*").From(dao.GetTableD("role")).Where("id>0 and department>=0").OrderBy("id")
	if len(param.Departments) > 0 {
		builder = pghelper.WhereUnnestInt(builder, "department in ", param.Departments)
	}
	sql, args := builder.MustSql()
	return dao.scan(sql, args)
}
func (dao *Dao) ListByDepartment(did int32) []*model.Role {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("role")).Where("id>0 and department=?", did).OrderBy("id").MustSql()
	return dao.scan(sql, args)
}

// privilege
func (dao *Dao) ListPrivileges() []*model.PrivilegeConstant {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("privilege_constant")).Where("type<>'dev'").OrderBy("sort").MustSql()
	return dao.scanPrivilege(sql, args)
}
