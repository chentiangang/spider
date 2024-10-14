package storage

type Storage[T any] interface {
	Save(data T) error
	Delete() error
}
