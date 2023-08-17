package roledao

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/departmentdao"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
	"github.com/mizuki1412/go-core-kit/service/sqlkit/pghelper"
)

type Dao struct {
	sqlkit.Dao[model.Role]
}

var meta = sqlkit.InitModelMeta(&model.Role{})

const (
	ResultDefault byte = iota
	ResultNone
)

func New(tx ...*sqlx.Tx) *Dao {
	return NewWithSchema("", tx...)
}
func NewWithSchema(schema string, tx ...*sqlx.Tx) *Dao {
	dao := &Dao{}
	dao.SetSchema(schema)
	if len(tx) > 0 {
		dao.TX = tx[0]
	}
	dao.Cascade = func(obj *model.Role) {
		switch dao.ResultType {
		case ResultDefault:
			if obj.Department != nil {
				obj.Department = departmentdao.NewWithSchema(dao.Schema).FindById(obj.Department.Id)
			}
		case ResultNone:
			obj.Department = nil
		}
	}
	return dao
}
func (dao *Dao) scanPrivilege(sql string, args []any) []*model.PrivilegeConstant {
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

func (dao *Dao) FindById(id int32) *model.Role {
	sql, args := dao.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where("id=?", id).MustSql()
	return dao.ScanOne(sql, args)
}
func (dao *Dao) FindByName(name string) *model.Role {
	sql, args := dao.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where("name=?", name).Limit(1).MustSql()
	return dao.ScanOne(sql, args)
}
func (dao *Dao) ListFromRootDepart(id int32) []*model.Role {
	where := fmt.Sprintf(`id>0 and department in ( with recursive t(id) as( values(%d) union all select d.id from %s d, t where t.id=d.parent) select id from t )`, id, dao.GetTable(&model.Department{}))
	sql, args := dao.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where(where).OrderBy("id").MustSql()
	return dao.ScanList(sql, args)
}

type ListParam struct {
	Departments []int32
}

func (dao *Dao) List(param ListParam) []*model.Role {
	builder := dao.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where("id>0 and department>=0").OrderBy("id")
	if len(param.Departments) > 0 {
		builder = pghelper.WhereUnnestInt(builder, "department in ", param.Departments)
	}
	sql, args := builder.MustSql()
	return dao.ScanList(sql, args)
}
func (dao *Dao) ListByDepartment(did int32) []*model.Role {
	sql, args := dao.Builder().Select(meta.GetColumns()).From(meta.GetTableName(dao.Schema)).Where("id>0 and department=?", did).OrderBy("id").MustSql()
	return dao.ScanList(sql, args)
}

// ListPrivileges privilege
func (dao *Dao) ListPrivileges() []*model.PrivilegeConstant {
	sql, args := dao.Builder().Select("*").From(sqlkit.GetSchemaTable(dao.Schema, "privilege_constant")).OrderBy("sort").MustSql()
	return dao.scanPrivilege(sql, args)
}
