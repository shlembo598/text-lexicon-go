package postgres

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"

	"github.com/shlembo598/text-lexicon-go/internal/config"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

func NewPsqlDB(c *config.Config) (*sqlx.DB, error) {
	const op = "pkg.db.postgres.newPsqlDB"

	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Postrgres.Host,
		c.Postrgres.Port,
		c.Postrgres.User,
		c.Postrgres.Dbname,
		c.Postrgres.Password,
	)

	db, err := sqlx.Connect(c.Postrgres.Driver, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil

}
