package load

import (
	"github.com/jackc/pgx"
	"math/rand"
	"time"
)

type Repo struct {
	pool *pgx.ConnPool
}

func NewRepo(pool *pgx.ConnPool) *Repo {
	return &Repo{
		pool: pool,
	}
}

func (r Repo) InsertToDB(data string) error {
	_, err := r.pool.Exec(`INSERT INTO job (data) VALUES ($1);`, data)
	return err
}

func RandomPayload(payloadSize int) string {
	rand.Seed(time.Now().UnixNano())

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, payloadSize)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
