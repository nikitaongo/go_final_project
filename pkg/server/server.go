package server

import (
	"fmt"
	"net/http"
)

func SetUpServer(dir, port string) {
	http.Handle("/", http.FileServer(http.Dir(dir)))
	fmt.Printf("Todo server started on http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}

}
