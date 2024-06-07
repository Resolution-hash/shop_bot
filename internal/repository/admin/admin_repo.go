package repository

type AdminRepo interface {
	AddAdmin(int64) error
	DeleteAdmin(int64) error
	IsAdmin(int64) (bool, error)
}

// type Admin struct {
// 	adminID int64
// }


