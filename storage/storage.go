package storage

type Storage interface {
	Save() error
	Delete() error
}

type ProjectSummaryRespStorage struct {
}

func (s *ProjectSummaryRespStorage) Save() error {
	return nil
}
