package main

import (
	"net/http"
	"os"

	"TP_Andreev/internal/transport/http/controller"
	"TP_Andreev/internal/transport/http/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Initialize router
	r := router.New()

	// Initialize controller
	pageCtrl := &controller.MainController{}

	// Register routes
	r.GET("/", pageCtrl.GetMainPage)

	// Start server
	http.ListenAndServe(":"+port, r)
}
