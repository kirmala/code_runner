package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"photo_editor/repository/ram_storage"
	"photo_editor/usecases/service"
	"log"

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
	addr := flag.String("addr", ":8080", "address for http server")

	taskRepo := ram_storage.NewTask()
	taskService := service.NewTask(taskRepo)
	taskHandlers := http.NewHandler(taskService)

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	taskHandlers.WithTaskHandlers(r)

	log.Printf("Starting server on %s", *addr)
	if err := pkgHttp.CreateAndRunServer(r, *addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}