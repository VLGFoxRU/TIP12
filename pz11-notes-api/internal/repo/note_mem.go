package repo

import (
	"errors"
	"sync"
	"time"

	"example.com/pz11-notes-api/internal/core"
)

type NoteRepoMem struct {
	mu    sync.Mutex
	notes map[int64]*core.Note
	next  int64
}

func NewNoteRepoMem() *NoteRepoMem {
	return &NoteRepoMem{
		notes: make(map[int64]*core.Note),
		next:  1,
	}
}

// Create - создание новой заметки
func (r *NoteRepoMem) Create(n core.Note) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n.ID = r.next
	n.CreatedAt = time.Now()
	r.notes[n.ID] = &n
	r.next++

	return n.ID, nil
}

// GetByID - получение заметки по ID
func (r *NoteRepoMem) GetByID(id int64) (*core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	note, ok := r.notes[id]
	if !ok {
		return nil, errors.New("note not found")
	}
	return note, nil
}

// GetAll - получение всех заметок
func (r *NoteRepoMem) GetAll() []*core.Note {
	r.mu.Lock()
	defer r.mu.Unlock()

	notes := make([]*core.Note, 0, len(r.notes))
	for _, n := range r.notes {
		notes = append(notes, n)
	}
	return notes
}

// Update - обновление заметки
func (r *NoteRepoMem) Update(id int64, n core.UpdateNoteRequest) (*core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	note, ok := r.notes[id]
	if !ok {
		return nil, errors.New("note not found")
	}

	if n.Title != "" {
		note.Title = n.Title
	}
	if n.Content != "" {
		note.Content = n.Content
	}
	now := time.Now()
	note.UpdatedAt = &now

	return note, nil
}

// Delete - удаление заметки
func (r *NoteRepoMem) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.notes[id]; !ok {
		return errors.New("note not found")
	}
	delete(r.notes, id)
	return nil
}
