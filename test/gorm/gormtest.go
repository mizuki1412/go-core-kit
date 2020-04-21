package gorm

import (
	"github.com/jinzhu/gorm"
	"log"
	"mizuki/project/core-kit/service/logkit"
)

func Test() {
	db, err := gorm.Open("postgres", "host=47.96.160.63 port=5432 user=postgres password=503149 dbname=teach-platform-server sslmode=disable")
	if err!=nil{
		logkit.Fatal(err.Error())
	}
	rows, err :=db.Raw("select * from admin_user limit ?", 10).Rows()
	if err!=nil{
		logkit.Fatal(err.Error())
	}
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache { //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}
	var list []map[string]interface{} //返回的切片
	for rows.Next() {
		_ = rows.Scan(cache...)

		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{}) //取实际类型
		}
		list = append(list, item)
	}
	log.Println(list)
	defer rows.Close()
	defer db.Close()
}
