package repository

import "golearn-structured/internal/model"

type TodoRepository interface {
	Create(todo model.Todo) error
	// Menggunakan List berupa Pointer ([]*model.Todo) agar yang dilempar hanya "alamat memorinya" saja.
	// Ini membuat aplikasi Anda ngebut (ringan) walau me-return 10.000 list data todos sekaligus.
	GetAllByUserID(userID int) ([]*model.Todo, error)
	GetOne(id, userID int) (model.Todo, error)
	Update(id, userID int, title, note, imageURL string) error
	Delete(id int, userID int) error
}
