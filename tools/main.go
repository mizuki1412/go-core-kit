package main

import (
	"github.com/mizuki1412/go-core-kit/class"
	"net/http"
)

func httpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	panic(http.ListenAndServe(":8080", nil))
}

type Bean struct {
	Id     int64           `json:"id" pk:"true" db:"id"`
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

	SQL2Struct("/Users/ycj/Downloads/demo.sql", "/Users/ycj/Downloads/dest.go")

	//sqlkit.Update(&Bean{Id: 11, Name: class.String{String: "qww", Valid: true}})

}
