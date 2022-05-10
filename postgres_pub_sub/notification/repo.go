package notification

import (
	"github.com/jackc/pgx"
	"time"
)

const minReconnectInterval, maxReconnectInterval = 10 * time.Second, time.Minute

type repo struct {
	DB *pgx.ConnPool
}

func NewRepo(db *pgx.ConnPool) *repo {
	return &repo{
		DB: db,
	}
}
