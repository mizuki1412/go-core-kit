package privilegedao

import (
	"github.com/mizuki1412/go-core-kit/v2/mod/user/model"
	"github.com/mizuki1412/go-core-kit/v2/service/sqlkit"
)

type Dao struct {
	sqlkit.Dao[model.PrivilegeConstant]
}

func New(ds ...*sqlkit.DataSource) Dao {
	return Dao{sqlkit.New[model.PrivilegeConstant](ds...)}
}

func (dao Dao) ListPrivileges() []*model.PrivilegeConstant {
	return dao.Select().OrderBy("sort").List()
}
