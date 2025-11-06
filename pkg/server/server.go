package server

import (
	"fmt"
	"gofinalproject/pkg/api"
	"net/http"
	"os"
)

func SetUpServer(dir, port string) {
	http.Handle("/", http.FileServer(http.Dir(dir)))
	api.Init()
	fmt.Printf("Todo server started on http://localhost:%s\n", port)

	if os.Getenv("TODO_PASSWORD") != "" {
		fmt.Println("Authorization with password required.")
	} else {
		fmt.Println("No authorization with password required.")
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
