package main

import (
	"log"

	"net/http"

	"settlements/internal/config"
	"settlements/internal/db"
	"settlements/internal/db/migrations"
	"settlements/internal/repo"
	"settlements/internal/service"
	"settlements/internal/transport/http/controller"
	"settlements/internal/transport/http/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("congif load failed: %v", err)
	}

	db, err := db.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}

	//Initialize router
	r := router.New()

	repo := repo.New(db)

	service := service.New(repo)

	// Initialize controller
	pageCtrl := controller.New(service)

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register routes
	r.GET("/", pageCtrl.GetMainPage)

	// Start server with both router and static handler
	http.Handle("/", r)
	http.ListenAndServe(":"+cfg.Server.Port, nil)
}
