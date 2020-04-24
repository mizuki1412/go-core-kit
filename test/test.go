package main

import (
	"log"
	"net/http"
	"runtime"
)

func httpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	panic(http.ListenAndServe(":8080", nil))
}

func main() {
	_, file, line, _ := runtime.Caller(1)
	log.Println(file, line)
	//gorm.Test()

	//corekit.Waiting()
}
