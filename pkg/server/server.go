package server

import (
	"fmt"
	"gofinalproject/pkg/api"
	"net/http"
)

func SetUpServer(dir, port string) {
	http.Handle("/", http.FileServer(http.Dir(dir)))
	api.Init()
	fmt.Printf("Todo server started on http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}

}
