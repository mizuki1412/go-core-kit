package main

import (
	"fmt"
	"mizuki/project/core-kit/library/stringkit"
	"net/http"
	"time"
)

func httpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	panic(http.ListenAndServe(":8080", nil))

}

func main() {
	fmt.Println(stringkit.ToString(int8(7)))
	fmt.Println(stringkit.ToString(time.Now().Unix()))
	fmt.Println(stringkit.ToString(float32(7.101)))
	fmt.Println(stringkit.ToString("sddd"))
}