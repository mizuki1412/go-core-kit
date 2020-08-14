//go:generate pkger -include /swagger-ui
package main

import (
	"github.com/mizuki1412/go-core-kit/cmd"
	"net/http"
)

func httpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	panic(http.ListenAndServe(":8080", nil))
}

func main() {
	cmd.Execute()
}
