package main

import (
	"net/http"
)

func httpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	panic(http.ListenAndServe(":8080", nil))

}

func main() {
	//logkit.Info("tets info")
	//logkit.Info("tets info2", logkit.Param{Key: "key1", Val: 1})
	//logkit.Info("tets info3")
	//defer logkit.Sync()
}
