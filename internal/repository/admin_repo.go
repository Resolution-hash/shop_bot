package repository

type AdminRepo interface {
	AddAdmin(int64) error
	deleteAdmin(int64) error
	IsAdmin(int64) bool
}

// type Admin struct {
// 	adminID int64
// }


