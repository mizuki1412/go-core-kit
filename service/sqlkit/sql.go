package sqlkit

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"mizuki/project/core-kit/service/configkit"
)

/**
数据库连接和CRUD的通用接口
rows转struct和map
在项目的main中须要导入driver
//todo middleware中加入事务
*/

var driver string
var db *sqlx.DB

func Driver() string {
	if db == nil {
		DB()
	}
	return driver
}
func DB() *sqlx.DB {
	if !configkit.Exist(ConfigKeyDBDriver) || !configkit.Exist(ConfigKeyDBHost) || !configkit.Exist(ConfigKeyDBPort) || !configkit.Exist(ConfigKeyDBPwd) || !configkit.Exist(ConfigKeyDBUser) || !configkit.Exist(ConfigKeyDBName) {
		panic("sqlkit: database config error")
	}
	if db == nil {
		driver = configkit.GetString(ConfigKeyDBDriver, "")
		db = sqlx.MustConnect(driver, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", configkit.GetString(ConfigKeyDBHost, ""), configkit.GetInt(ConfigKeyDBPort, 0), configkit.GetString(ConfigKeyDBUser, ""), configkit.GetString(ConfigKeyDBPwd, ""), configkit.GetString(ConfigKeyDBName, "")))
		//db.SetConnMaxLifetime()
	}
	return db
}

// dest use pointer
func QueryStruct(dest interface{}, sql string, args []interface{}, err error) {
	if err != nil {
		// todo
		panic(err)
	}
	log.Println(sql)
	err = DB().Select(dest, sql, args...)
	if err != nil {
		//todo
	}
}

func QueryMap(sql string, args []interface{}, err error) []map[string]interface{} {
	if err != nil {
		// todo
		panic(err)
	}
	log.Println(sql)
	list := []map[string]interface{}{}
	rows, _ := DB().Queryx(sql, args...)
	for rows.Next() {
		m := map[string]interface{}{}
		err := rows.MapScan(m)
		if err != nil {
			// todo
		}
		list = append(list, m)
	}
	return list
}

func Insert() {

}

func Update() {

}

func Delete() {

}
