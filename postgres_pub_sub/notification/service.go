package notification

import (
	"context"
	"encoding/json"
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
	ID                  int64       `json:"id"`
	JobID               int64       `json:"job_id"`
	Status              string      `json:"status"`
	PreviousStatus      string      `json:"prev_status"`
	CreateTime          time.Time   `json:"create_time"`
	Data                interface{} `json:"data"`
	ChannelName         string      `json:"channel_name"`
	InternalChannelName string      `json:"notification_channel"`
	ChannelID           uint32      `json:"channel_id"`
}

type service struct {
	closeCh     chan struct{}
	ChannelName string
	ProcessFn   func(notification *Notification)
	Repo        *repo
}

func NewService(repo *repo, channelName string, processFn func(notification *Notification)) *service {
	if processFn == nil {
		processFn = func(n *Notification) {
			j, err := json.Marshal(n)
			if err != nil {
				log.Warn().Err(err).Msg("failed to marshal notification")
			} else {
				log.Debug().RawJSON("notification", j).Send()
			}
		}
	}

	return &service{
		ChannelName: channelName,
		Repo:        repo,
		ProcessFn:   processFn,
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

	go func() {
		for {
			select {
			case _ = <-s.closeCh:
				return
			default:
			}

			notification, err := conn.WaitForNotification(ctx)
			if err != nil {
				return
			}

			n := Notification{
				ChannelID:           notification.PID,
				ChannelName:         s.ChannelName,
				InternalChannelName: notification.Channel,
			}

			//log.Debug().Str("notification_channel", notification.Channel).Uint32("channel_id", notification.PID).Str("payload", notification.Payload).Send()
			if err = json.Unmarshal([]byte(notification.Payload), &n); err != nil {
				log.Warn().Err(err).Msg("failed json unmarshal on notification")
			}

			s.ProcessFn(&n)
		}
	}()

	return nil
}

func (s service) Stop() {
	s.closeCh <- struct{}{}
}
