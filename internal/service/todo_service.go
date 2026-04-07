package service

import (
	"fmt"
	"golearn-structured/internal/repository"
	"os"
	"strings"
)

type TodoService struct {
	repo *repository.TodoRepository
}

func NewTodoService(repo *repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(userID int, title, note, ImageUrl string) error {
	if title == "" {
		return fmt.Errorf("Judul wajib diisi")
	}

	todo := repository.Todo{
		UserID:   userID,
		Title:    title,
		Note:     note,
		ImageUrl: ImageUrl,
	}
	return s.repo.Create(todo)
}

func (s *TodoService) GetUserTodos(userID int) ([]repository.Todo, error) {
	return s.repo.GetAllByUserID(userID)
}

func (s *TodoService) EditTodo(id, userID int, title, note, imageURL string) error {
	oldTodo, err := s.repo.GetOne(id, userID)

	if err != nil {
		return fmt.Errorf("data tidak ditemukan")
	}

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
	return s.repo.Delete(id, userID)
}
