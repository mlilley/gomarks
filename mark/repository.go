package mark

type Repository interface {
	FindAll() ([]*Mark, error)
	FindByID(id string) (*Mark, error)
	Create(m *Mark) (*Mark, error)
	Update(m *Mark) error
	DeleteByID(id string) (bool, error)
	DeleteAll() error
}
