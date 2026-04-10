package repository

type UserRepository interface {
	GetByUsername(username string) (int, string, string, error)
	IsUserExist(username string) bool
	Insert(username, password string) error
	UpdatePassword(username, password string) error
	GetProfile(id int, username string) (int, string, error)
}
