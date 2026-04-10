package repository

import "golearn-structured/internal/model"

type TodoRepository interface {
	Create(todo model.Todo) error
	GetAllByUserID(userID int) ([]model.Todo, error)
	GetOne(id, userID int) (model.Todo, error)
	Update(id, userID int, title, note, imageURL string) error
	Delete(id int, userID int) error
}
