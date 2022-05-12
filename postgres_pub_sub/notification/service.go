package notification

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type SubMessage struct {
	JobID          int64  `json:"job_id"`
	NotificationID int64  `json:"notification_id"`
	Status         string `json:"status"`
}

type Notification struct {
	ID             int64       `json:"id"`
	JobID          int64       `json:"job_id"`
	Status         string      `json:"status"`
	PreviousStatus string      `json:"prev_status"`
	CreateTime     time.Time   `json:"create_time"`
	Data           interface{} `json:"data"`
}

type service struct {
	ChannelName string
	cancelCtx   context.CancelFunc
	Repo        *repo
}

func NewService(repo *repo, channelName string) *service {
	return &service{
		ChannelName: channelName,
		Repo:        repo,
	}
}

func (s service) Start(ctx context.Context) error {
	conn, err := s.Repo.DB.Acquire()
	if err != nil {
		return errors.Wrap(err, "failed to acquire connection")
	}

	if err = conn.Listen(s.ChannelName); err != nil {
		return errors.Wrap(err, "unable to listen to channel")
	}

	log.Debug().Msg("Notification:Start")

	for {
		notification, err := conn.WaitForNotification(ctx)
		if err != nil {
			return errors.Wrap(err, "error listening to channel")
		}

		log.Debug().Str("notification_channel", notification.Channel).Uint32("channel_id", notification.PID).Str("payload", notification.Payload).Send()
		if err = process(notification.Payload); err != nil {
			return errors.Wrap(err, "error in notification callback")
		}
	}
}

func process(payload string) error {
	// Metrics we care about
	// How many notifications can we put per second
	// How many notifications can we process per second
	// When does the system start
	return nil
}
