package postgres

import (
	"database/sql"

	"github.com/demkowo/goquery/model"
)

const (
	GATHER_LINKS = "INSERT INTO links (url, details) VALUES ($1, $2) ON CONFLICT (url) DO NOTHING;"
)

type CollectorRepository interface {
	GatherAllLinks(links []*model.Links) error
}

type collectorRepository struct {
	db *sql.DB
}

func New(db *sql.DB) CollectorRepository {
	return &collectorRepository{
		db: db,
	}
}

func (r collectorRepository) GatherAllLinks(links []*model.Links) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		// If panic occurs, rollback the transaction
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
			// If an error occurred, rollback the transaction
		} else if err != nil {
			_ = tx.Rollback()
			// If all operations succeed, commit the transaction
		} else {
			err = tx.Commit()
		}
	}()

	// Prepare the statement
	stmt, err := tx.Prepare(GATHER_LINKS)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Iterate through the links and execute the statement for each link
	for _, link := range links {
		_, err = stmt.Exec(link.Url, link.Details)
		if err != nil {
			return err
		}
	}

	return nil
}
