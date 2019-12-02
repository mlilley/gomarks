package user

type Repository interface {
	FindAll() ([]*User, error)
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	Create(u *User) (*User, error)
	Update(u *User) error
	DeleteByID(id string) (bool, error)
	DeleteByEmail(email string) (bool, error)
	DeleteAll() error
}
