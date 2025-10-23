package main

import (
	"gofinalproject/pkg/db"
	"gofinalproject/pkg/server"
	"os"

	_ "modernc.org/sqlite"
)

func main() {

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	db.Init(dbFile)

	webDir := "./web"
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	server.SetUpServer(webDir, port)
}
