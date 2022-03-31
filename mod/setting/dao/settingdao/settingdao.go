package settingdao

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

/** more_setting: id, data */

/// auto template
type Dao struct {
	sqlkit.Dao
}

func New(schema string, tx ...*sqlkit.Dao) *Dao {
	dao := &Dao{}
	dao.NewHelper(schema, tx...)
	return dao
}
func (dao *Dao) scanOne(sql string, args []interface{}) *map[string]interface{} {
	rows := dao.Query(sql, args...)
	defer rows.Close()
	for rows.Next() {
		m := map[string]interface{}{}
		err := rows.MapScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		return &m
	}
	return nil
}

////

func (dao *Dao) Set(data map[string]interface{}) {
	sql, args, err := sqlkit.Builder().Update(dao.GetTableD("more_setting")).Set("data", jsonkit.ToString(data)).Where("id=?", 1).ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	dao.Exec(sql, args...)
}

func (dao *Dao) Get() map[string]interface{} {
	sql, args := sqlkit.Builder().Select("data").From(dao.GetTableD("more_setting")).Where("id=?", 1).MustSql()
	rows := dao.Query(sql, args...)
	var data string
	defer rows.Close()
	for rows.Next() {
		_ = rows.Scan(&data)
	}
	return jsonkit.ParseMap(data)
}
