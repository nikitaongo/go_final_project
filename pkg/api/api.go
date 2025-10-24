package api

import "net/http"

// Init runs nextDayHandler for /api/nextdate endpoint
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
}
