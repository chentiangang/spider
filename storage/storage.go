package storage

type Storage interface {
	Save(data interface{}) error
	ReplaceOldData(data interface{}) error
}
