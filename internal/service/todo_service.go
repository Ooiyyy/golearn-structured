package service

import (
	"fmt"
	"golearn-structured/internal/model"
	"golearn-structured/internal/repository"
	"os"
	"strings"
)

type TodoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(userID int, title, note, ImageUrl string) error {
	// Rule bisnis minimum: judul wajib diisi.
	if title == "" {
		return fmt.Errorf("Judul wajib diisi")
	}

	// Struct dipakai sebagai kontrak data antar layer service -> repository.
	todo := model.Todo{
		UserID:   userID,
		Title:    title,
		Note:     note,
		ImageUrl: ImageUrl,
	}
	return s.repo.Create(todo)
}

func (s *TodoService) GetUserTodos(userID int) ([]*model.Todo, error) {
	// Service tetap tipis saat tidak ada aturan tambahan.
	return s.repo.GetAllByUserID(userID)
}

func (s *TodoService) EditTodo(id, userID int, title, note, imageURL string) error {
	// Ambil data lama untuk mendukung partial update.
	oldTodo, err := s.repo.GetOne(id, userID)

	if err != nil {
		return fmt.Errorf("data tidak ditemukan")
	}

	// Fallback ke nilai lama jika field tidak dikirim client.
	finalTitle := title
	if finalTitle == "" {
		finalTitle = oldTodo.Title
	}
	finalNote := note
	if finalNote == "" {
		finalNote = oldTodo.Note
	}
	finalImageURL := imageURL
	if finalImageURL == "" {
		finalImageURL = oldTodo.ImageUrl
	} else {
		if oldTodo.ImageUrl != "" {
			// Jika ada gambar baru, gambar lama dibersihkan dari storage.
			oldPath := strings.Replace(oldTodo.ImageUrl, "http://localhost:8080/", "", 1)
			os.Remove(oldPath)
		}
	}
	// if title == "" {
	// 	return fmt.Errorf("judul wajib diisi")
	// }
	return s.repo.Update(id, userID, finalTitle, finalNote, finalImageURL)
}

func (s *TodoService) DeleteTodo(id, userID int) error {
	// Penghapusan data tetap lewat repository agar query SQL terpusat.
	return s.repo.Delete(id, userID)
}
