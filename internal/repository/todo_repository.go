package repository

import (
	"database/sql"
)

type Todo struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Title    string `json:"title"`
	Note     string `json:"note"`
	ImageUrl string `json:"image_url"`
}

type TodoRepository struct {
	DB *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{DB: db}
}

func (r *TodoRepository) Create(todo Todo) error {
	_, err := r.DB.Exec("INSERT INTO todos (user_id, title, note, image_url) VALUES (?,?,?,?)", todo.UserID, todo.Title, todo.Note, todo.ImageUrl)
	return err
}

func (r *TodoRepository) GetAllByUserID(userID int) ([]Todo, error) {

	rows, err := r.DB.Query("SELECT id, user_id, title, note, image_url FROM todos WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo

		err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Note, &t.ImageUrl)
		if err == nil {
			todos = append(todos, t)
		}
	}
	return todos, nil
}

func (r *TodoRepository) GetOne(id, userID int) (Todo, error) {
	var t Todo

	err := r.DB.QueryRow("SELECT id, user_id, title, note, image_url FROM todos WHERE id = ? AND user_id = ?", id, userID).Scan(&t.ID, &t.UserID, &t.Title, &t.Note, &t.ImageUrl)
	return t, err
}
func (r *TodoRepository) Update(id, userID int, title, note, imageURL string) error {
	_, err := r.DB.Exec("UPDATE todos SET title = ?, note = ?, image_url = ? WHERE id = ? AND user_id = ?", title, note, imageURL, id, userID)
	return err
}

func (r *TodoRepository) Delete(id int, userID int) error {
	_, err := r.DB.Exec("DELETE FROM todos WHERE id = ? AND user_id = ?", id, userID)
	return err
}
