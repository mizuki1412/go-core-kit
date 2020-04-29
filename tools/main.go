package main

import (
	"mizuki/project/core-kit/class"
	"net/http"
)

func httpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	panic(http.ListenAndServe(":8080", nil))
}

type Bean struct {
	Id     int64           `json:"id"`
	Name   class.String    `json:"name"`
	Age    class.Int32     `json:"age"`
	Extend class.MapString `json:"extend"`
	Dt     class.Time      `json:"dt,omitempty"`
	List1  class.ArrString `json:"list1"`
	List2  class.ArrInt    `json:"list2"`
}

type loginByUsernameParam struct {
	Username string      `form:"username" description:"用户名" validate:"required"`
	Pwd      class.Int32 `form:"pwd" validate:"required" default:"12"`
	Schema   string      `form:"schema" default:"public"`
}

func main() {
	//SQL2Struct("/Users/ycj/Downloads/demo.sql", "/Users/ycj/Downloads/dest.go")
	//v := loginByUsernameParam{
	//	Username: "z",
	//	Pwd:      class.Int32{Valid: false},
	//	Schema:   "pub",
	//}
	//rt := reflect.TypeOf(&v).Elem()
	//rv := reflect.ValueOf(&v).Elem()
	//for i := 0; i < rt.NumField(); i++ {
	//	log.Println(1, rv.Field(i))
	//	log.Println(2, rt.Field(i).Name)
	//}
}
