package main

import (
	"log"

	"github.com/RakhaYandra/lifeos/internal/config"
	"github.com/RakhaYandra/lifeos/internal/handler"
	"github.com/RakhaYandra/lifeos/internal/repository"
	"github.com/RakhaYandra/lifeos/internal/service"
)

func main() {
	cfg := config.Load()
	db, err := repository.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close() //nolint:errcheck

	users := &repository.UserRepository{DB: db}
	settings := &repository.SettingsRepository{DB: db}
	areas := &repository.LifeAreaRepository{DB: db}

	authH := &handler.AuthHandler{
		Svc:      &service.AuthService{Users: users, Secret: cfg.JWTSecret},
		Users:    users,
		Settings: settings,
	}
	setH := &handler.SettingsHandler{Settings: settings}
	areaH := &handler.LifeAreaHandler{Areas: areas}

	r := handler.NewRouter(&handler.Deps{Auth: authH, Settings: setH, LifeArea: areaH}, cfg.JWTSecret, cfg.FrontendURL)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
