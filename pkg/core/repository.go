package core

// Repository defines the API repository implementation should follow.
type Repository interface {
	Find(id string) (*Kudo, error)
	FindAll(selector map[string]interface{}) ([]*Kudo, error)
	Delete(kudo *Kudo) error
	Update(kudo *Kudo) error
	Create(kudo ...*Kudo) error
	Count() (int, error)
}
