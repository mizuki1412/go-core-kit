package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
)

func jsonDemo() {
	const jsonStr = `{"name":{"first":"Janet","last":"Prichard"},"age":[1,2,3]}`
	val := gjson.Get(jsonStr, "name.first")
	fmt.Println(val.Exists())
	fmt.Println(val.String())
	m := gjson.Parse(jsonStr).Value().(map[string]interface{})
	fmt.Println(m["age"].([]interface{})[1])

	d,_ :=json.Marshal(m)
	fmt.Println(d)
	fmt.Println(string(d))
}

func httpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	panic(http.ListenAndServe(":8080", nil))

}

func main() {

}