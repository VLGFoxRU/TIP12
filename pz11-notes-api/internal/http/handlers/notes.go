package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/pz11-notes-api/internal/core"
	"example.com/pz11-notes-api/internal/repo"
	"github.com/go-chi/chi/v5"
)

type NoteHandler struct {
	repo *repo.NoteRepoMem
}

func NewNoteHandler(repo *repo.NoteRepoMem) *NoteHandler {
	return &NoteHandler{repo: repo}
}

// CreateNote godoc
// @Summary Создать новую заметку
// @Description Создаёт и сохраняет новую заметку в системе
// @Tags notes
// @Accept json
// @Produce json
// @Param input body core.CreateNoteRequest true "Данные новой заметки"
// @Success 201 {object} core.Note "Заметка успешно создана"
// @Failure 400 {object} map[string]string "Некорректные данные (пустые поля)"
// @Failure 500 {object} map[string]string "Ошибка при создании заметки"
// @Router /notes [post]
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req core.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	note := core.Note{
		Title:   req.Title,
		Content: req.Content,
	}

	id, err := h.repo.Create(note)
	if err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	note.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

// GetNote godoc
// @Summary Получить заметку по ID
// @Description Возвращает полную информацию о конкретной заметке по её идентификатору
// @Tags notes
// @Produce json
// @Param id path int true "Идентификатор заметки"
// @Success 200 {object} core.Note "Заметка найдена"
// @Failure 400 {object} map[string]string "Некорректный формат ID"
// @Failure 404 {object} map[string]string "Заметка не найдена"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /notes/{id} [get]
func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	note, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// GetAllNotes godoc
// @Summary Получить список всех заметок
// @Description Возвращает список всех заметок с поддержкой пагинации и фильтра по названию
// @Tags notes
// @Produce json
// @Param page query int false "Номер страницы (по умолчанию 1)" default(1)
// @Param limit query int false "Размер страницы (по умолчанию 10)" default(10)
// @Param q query string false "Поиск по названию (title)"
// @Success 200 {array} core.Note "Список заметок"
// @Header 200 {integer} X-Total-Count "Общее количество заметок в системе"
// @Failure 400 {object} map[string]string "Некорректные параметры пагинации"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /notes [get]
func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes := h.repo.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// UpdateNote godoc
// @Summary Обновить заметку (частичное обновление)
// @Description Обновляет один или несколько полей заметки. Поля, не указанные в запросе, не изменяются
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор заметки для обновления"
// @Param input body core.UpdateNoteRequest true "Поля для обновления (все опциональны)"
// @Success 200 {object} core.Note "Заметка успешно обновлена"
// @Failure 400 {object} map[string]string "Некорректный ID или тело запроса"
// @Failure 404 {object} map[string]string "Заметка не найдена"
// @Failure 500 {object} map[string]string "Ошибка при обновлении заметки"
// @Router /notes/{id} [patch]
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req core.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	note, err := h.repo.Update(id, req)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// DeleteNote godoc
// @Summary Удалить заметку
// @Description Удаляет заметку из системы по её идентификатору. Операция необратима
// @Tags notes
// @Param id path int true "Идентификатор заметки для удаления"
// @Success 204 "Заметка успешно удалена (No Content)"
// @Failure 400 {object} map[string]string "Некорректный формат ID"
// @Failure 404 {object} map[string]string "Заметка не найдена"
// @Failure 500 {object} map[string]string "Ошибка при удалении заметки"
// @Router /notes/{id} [delete]
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
