// Package main Notes API server.
//
// @title           Notes API
// @version         1.0
// @description     Учебный REST API для заметок (CRUD).
// @contact.name    Backend Course
// @contact.email   example@university.ru
// @BasePath        /api/v1
package main

import (
	"log"
	"net/http"

	"example.com/pz11-notes-api/internal/http"
	"example.com/pz11-notes-api/internal/repo"

	httpSwagger "github.com/swaggo/http-swagger"
	
	_ "example.com/pz11-notes-api/docs"
	// _ нужен потому что пакет docs используем только для регистрации Swagger при загрузке приложения, а сам пакет не вызываем.
)

func main() {
	// Инициализация репозитория
	noteRepo := repo.NewNoteRepoMem()

	// Создание маршрутизатора
	r := router.NewRouter(noteRepo)

	r.Get("/docs/*", httpSwagger.WrapHandler) // для chi: r.Get

	// Запуск сервера
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
