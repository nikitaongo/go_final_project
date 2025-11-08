package server

import (
	"fmt"
	"gofinalproject/pkg/api"
	"net/http"
)

type Config struct {
	Password string
	Port     string
	DBFile   string
}

func StartServer(dir string, cfg Config) {
	http.Handle("/", http.FileServer(http.Dir(dir)))
	api.Init(cfg.Password)
	fmt.Printf("Todo server started on http://localhost:%s\n", cfg.Port)

	if cfg.Password != "" {
		fmt.Println("Authorization with password required.")
	} else {
		fmt.Println("No authorization with password required.")
	}

	err := http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		panic(err)
	}
}
