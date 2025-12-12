package router

import (
	"github.com/go-chi/chi/v5"
	"example.com/pz11-notes-api/internal/http/handlers"
	"example.com/pz11-notes-api/internal/repo"
)

func NewRouter(noteRepo *repo.NoteRepoMem) chi.Router {
	r := chi.NewRouter()

	// Инициализация обработчиков
	noteHandler := handlers.NewNoteHandler(noteRepo)

	// Маршруты API v1
	r.Route("/api/v1/notes", func(r chi.Router) {
		r.Post("/", noteHandler.CreateNote)        // POST /api/v1/notes
		r.Get("/", noteHandler.GetAllNotes)        // GET /api/v1/notes
		r.Get("/{id}", noteHandler.GetNote)        // GET /api/v1/notes/{id}
		r.Patch("/{id}", noteHandler.UpdateNote)   // PATCH /api/v1/notes/{id}
		r.Delete("/{id}", noteHandler.DeleteNote)  // DELETE /api/v1/notes/{id}
	})

	return r
}
