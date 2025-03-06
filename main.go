package main

import (
	//"fmt"
	"fmt"
	"log"
	"photo_editor/config"
	"photo_editor/repository/ram_storage"
	"photo_editor/usecases/service"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"photo_editor/api/http"
	_ "photo_editor/docs"
	pkgHttp "photo_editor/pkg/http"
)

// @title photo_editor
// @version 1.0
// @description This is a photo editor.

// @host localhost:8080
// @BasePath /
func main() {
	appFlags := config.ParseFlags()
	var cfg config.HTTPConfig
	config.MustLoad(appFlags.ConfigPath, &cfg)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)


	taskRepo := ram_storage.NewTask()
	sessionRepo := ram_storage.NewSession()
	userRepo := ram_storage.NewUser()

	taskService := service.NewTask(taskRepo, sessionRepo)
	userService := service.NewUser(userRepo, sessionRepo)

	taskHandlers := http.NewTaskHandler(taskService)
	userHandlers := http.NewUserHandler(userService)


	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	taskHandlers.WithTaskHandlers(r)
	userHandlers.WithUserHandlers(r)

	log.Printf("Starting server on %s", addr)
	if err := pkgHttp.CreateAndRunServer(r, addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}