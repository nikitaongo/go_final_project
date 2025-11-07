package main

import (
	"gofinalproject/pkg/db"
	"gofinalproject/pkg/server"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	cfg := server.Config{
		Password: os.Getenv("TODO_PASSWORD"),
		Port:     getenvDefault("TODO_PORT", "7540"),
		DBFile:   getenvDefault("TODO_DBFILE", "scheduler.db"),
	}

	db.Init(cfg.DBFile)
	defer db.Db.Close()

	webDir := "./web"
	server.StartServer(webDir, cfg)
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
