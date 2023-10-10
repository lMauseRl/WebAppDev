package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/lMauseRl/WebAppDev/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetPeriodsByID(id int) (*ds.Periods, error) {
	periods := &ds.Periods{}

	err := r.db.First(periods, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return periods, nil
}

func (r *Repository) DeletePeriods(id int) error {
	return r.db.Exec("UPDATE periods SET status = 'удалён' WHERE id=?", id).Error
}

func (r *Repository) CreatePeriods(periods ds.Periods) error {
	return r.db.Create(periods).Error
}

func (r *Repository) GetAllPeriods() ([]ds.Periods, error) {
	var periods []ds.Periods
	err := r.db.Find(&periods, "status = 'действует'").Error
	if err != nil {
		return nil, err
	}

	return periods, nil
}
