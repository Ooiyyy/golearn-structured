package mysql

import (
	"database/sql"
	"golearn-structured/internal/model"
	"golearn-structured/internal/repository"
)

type todoRepositoryImpl struct {
	DB *sql.DB
}

func NewTodoRepository(db *sql.DB) repository.TodoRepository {
	return &todoRepositoryImpl{DB: db}
}

func (r *todoRepositoryImpl) Create(todo model.Todo) error {
	// Exec dipakai untuk statement yang mengubah data (INSERT/UPDATE/DELETE).
	_, err := r.DB.Exec("INSERT INTO todos (user_id, title, note, image_url) VALUES (?,?,?,?)", todo.UserID, todo.Title, todo.Note, todo.ImageUrl)
	return err
}

func (r *todoRepositoryImpl) GetAllByUserID(userID int) ([]model.Todo, error) {

	// Query mengembalikan banyak baris todo milik satu user.
	rows, err := r.DB.Query("SELECT id, user_id, title, note, image_url FROM todos WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var t model.Todo

		// Scan menyalin kolom SQL ke field struct via pointer.
		err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Note, &t.ImageUrl)
		if err == nil {
			todos = append(todos, t)
		}
	}
	return todos, nil
}

func (r *todoRepositoryImpl) GetOne(id, userID int) (model.Todo, error) {
	var t model.Todo

	// QueryRow cocok untuk satu baris target saat edit.
	err := r.DB.QueryRow("SELECT id, user_id, title, note, image_url FROM todos WHERE id = ? AND user_id = ?", id, userID).Scan(&t.ID, &t.UserID, &t.Title, &t.Note, &t.ImageUrl)
	return t, err
}
func (r *todoRepositoryImpl) Update(id, userID int, title, note, imageURL string) error {
	// Placeholder (?) menjaga parameter aman dan mudah di-bind.
	_, err := r.DB.Exec("UPDATE todos SET title = ?, note = ?, image_url = ? WHERE id = ? AND user_id = ?", title, note, imageURL, id, userID)
	return err
}

func (r *todoRepositoryImpl) Delete(id int, userID int) error {
	_, err := r.DB.Exec("DELETE FROM todos WHERE id = ? AND user_id = ?", id, userID)
	return err
}
