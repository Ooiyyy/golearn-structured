package model

type Todo struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Title    string `json:"title"`
	Note     string `json:"note"`
	ImageUrl string `json:"image_url"`
}
