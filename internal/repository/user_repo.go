package repository

type UserRepo interface {
	AddUser() error
}

type User struct {
	UserID     int
	First_name string
	Last_name  string
	User_name  string
}
