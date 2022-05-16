package load

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/rs/zerolog/log"
	"github.com/spencerreeves/snippets/postgres_pub_sub/notification"
	"github.com/spencerreeves/snippets/thread"
	"math"
	"time"
)

type producerConfig struct {
	Payload string
	Repo    *Repo
}

func RunLoadTest(pool *pgx.ConnPool, channelName string) {
	repo := NewRepo(pool)

	// Determine the number of rows to create
	// Total bytes written / size of one payload
	rowCount := int(math.Ceil(25000000 / 250))
	payload := "{\"a\":\"" + RandomPayload(250-8) + "\"}"
	config := producerConfig{
		Payload: payload,
		Repo:    repo,
	}
	channel := make(chan *notification.Notification, rowCount)

	ctx, cancel := context.WithCancel(context.Background())
	notificationRepo := notification.NewRepo(pool)
	notificationService := notification.NewService(notificationRepo, channelName, func(notification *notification.Notification) { channel <- notification })
	err := notificationService.Start(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service start")
	}

	notes, inserted, errs := 0, 0, 0
	threadPool := thread.NewPool[*notification.Notification, *producerConfig](1, 4, &rowCount, &config, channel, handleNotification(&notes), handleError(&errs), producerHandler(&inserted))
	for !threadPool.Closed {
		if notes == rowCount {
			threadPool.Close(true)
		}
		//log.Debug().Int("Inserted", inserted).Int("Received", notes).Int("Errors", errs).Msgf("Percent complete: %0.2f", float64(notes+inserted)/(float64(rowCount)*2)*100)
		time.Sleep(time.Millisecond * 500)
	}
	cancel()

	for _, m := range threadPool.Metrics() {
		log.Info().Time("Start", m.StartTime).Time("End", m.EndTime).Dur("Busy", m.BusyDuration).Dur("Idle", m.IdleDuration).Dur("Total Duration", m.EndTime.Sub(m.StartTime)).Int("Processed", m.ProcessedCount).Int("Errors", m.ErrorCount).Send()
	}
}

func handleNotification(counter *int) func(elem *notification.Notification) error {
	return func(n *notification.Notification) error {
		if n != nil {
			*counter++
		}
		return nil
	}
}

func handleError(counter *int) func(id string, e error) {
	return func(id string, e error) {
		*counter++
	}
}

func producerHandler(counter *int) func(index int, config *producerConfig) (*notification.Notification, error) {
	return func(index int, config *producerConfig) (*notification.Notification, error) {
		if err := config.Repo.InsertToDB(config.Payload); err != nil {
			return nil, err
		}

		*counter++
		return nil, nil
	}
}
