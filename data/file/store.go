package file

import "gorm.io/gorm"

type Store struct {
	DB *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) Create(p *Model) error {
	return s.DB.Create(p).Error
}

func (s *Store) Update(p *Model) error {
	return s.DB.Model(Model{}).Updates(p).Error
}

func (s *Store) GetById(id string) (*Model, error) {
	var retVal Model
	if err := s.DB.Where(Model{ID: id}).First(&retVal).Error; err != nil {
		return nil, err
	}
	return &retVal, nil
}
