package main

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spencerreeves/snippets/postgres_pub_sub/load"
	"github.com/spencerreeves/snippets/postgres_pub_sub/notification"
	"os"
	"sync"
	"time"
)

func main() {
	c, err := LoadConfig()

	if c.Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if err != nil {
		log.Panic().Err(err).Msg("unable to read config")
	}

	connConfig, err := pgx.ParseConnectionString(c.DbUrl)
	if err != nil {
		log.Panic().Err(err).Msg("unable to read db config")
	}

	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		AfterConnect:   nil,
		MaxConnections: 20,
		AcquireTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Panic().Err(err).Msg("unable to connect to db")
	}

	switch c.TestType {
	case "load":
		load.RunLoadTest(connPool, c.ChannelName)
	default:
		LogAllNotifications(connPool, c.ChannelName)
	}
}

func LogAllNotifications(pool *pgx.ConnPool, channel string) {
	notificationRepo := notification.NewRepo(pool)
	notificationService := notification.NewService(notificationRepo, channel, nil)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// This is a blocking call that will run until there is an error in the notification stream
		ctx := context.Background()
		if err := notificationService.Start(ctx); err != nil {
			log.Error().Err(err).Msg("error in notification service")
		}
	}()

	wg.Wait()
}
