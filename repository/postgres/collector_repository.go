package postgres

import "database/sql"

const ()

type CollectorRepository interface {
	FindAllLinks()
}

type collectorRepository struct {
	db *sql.DB
}

func New(db *sql.DB) CollectorRepository {
	return &collectorRepository{
		db: db,
	}
}

func (s collectorRepository) FindAllLinks() {}
