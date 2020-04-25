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
	//err:= filekit.WriteFile("./t/t.txt", []byte("sddda"))
	//log.Println(afero.Exists(afero.NewOsFs(),"./test/ss"))
}
