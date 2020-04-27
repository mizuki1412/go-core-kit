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

func main() {
	//log.Println(stringkit.Split("var xxx\tsss,233","[ ,\t]+"))
	//SQL2Struct("/Users/ycj/Downloads/demo.sql", "/Users/ycj/Downloads/dest.go")
}
