package repository

type UserRepo interface {
	AddUser(User) (bool,error)
}

type User struct {
	UserID     int
	First_name string
	Last_name  string
	User_name  string
	IsAdmin    int
}
