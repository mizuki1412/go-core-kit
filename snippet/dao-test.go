package snippet

import (
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/dao/userdao"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/model"
	"github.com/mizuki1412/go-core-kit/v2/service/sqlkit"
	"log"
)

type Bean struct {
	Id    int32
	Data  class.ArrInt
	Data2 *class.ArrInt
}

func testSQL() {
	dao := userdao.New(userdao.ResultDefault)
	u := &model.User{}
	u.Name.Set("test33")
	u.Role = &model.Role{Id: 1}
	dao.InsertObj(u)
	u.Extend.Put("abc", 1)
	dao.UpdateObj(u)
	log.Println(jsonkit.ToString(dao.SelectOneById(u.Id)))
	log.Println(jsonkit.ToString(dao.Select().Where("name=?", u.Name).List()))
}

func testArr() {
	dao := New()
	t := &Test2{}
	t.Ns.Set([]int64{2, 6, 7})
	//dao.InsertObj(t)
	log.Println(t)
	log.Println(jsonkit.ToString(dao.SelectOneById(1)))
	t.Ns.Set([]int64{4, 5, 3})
	//dao.UpdateObj(t)
	log.Println(jsonkit.ToString(dao.DataSource().DBPool.Stats()))
}

type Test2 struct {
	Id   int32        `json:"id,omitempty" db:"id" pk:"true" table:"test2" auto:"true"`
	Name class.String `db:"name"`
	Ns   class.ArrInt `db:"ns"`
}

type Dao struct {
	sqlkit.Dao[Test2]
}

func New(ds ...*sqlkit.DataSource) *Dao {
	return &Dao{sqlkit.New[Test2](ds...)}
}
