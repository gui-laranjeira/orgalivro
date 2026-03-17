package main

import (
	"log"

	"orgalivro/backend/internal/config"
	"orgalivro/backend/internal/db"
	"orgalivro/backend/internal/handler"
	"orgalivro/backend/internal/repository"
	"orgalivro/backend/internal/router"
	"orgalivro/backend/internal/service"
)

func main() {
	cfg := config.Load()

	gormDB, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := db.Migrate(gormDB); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	profileRepo := repository.NewProfileRepo(gormDB)
	bookRepo := repository.NewBookRepo(gormDB)
	entryRepo := repository.NewEntryRepo(gormDB)

	profileSvc := service.NewProfileService(profileRepo)
	bookSvc := service.NewBookService(bookRepo)
	entrySvc := service.NewEntryService(entryRepo, profileRepo)
	isbnSvc := service.NewISBNService()

	ph := handler.NewProfileHandler(profileSvc)
	bh := handler.NewBookHandler(bookSvc)
	eh := handler.NewEntryHandler(entrySvc)
	ih := handler.NewISBNHandler(isbnSvc)

	r := router.New(cfg, ph, bh, eh, ih)
	log.Printf("Starting server on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
