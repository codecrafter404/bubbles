package gql

import (
	"github.com/rs/zerolog/log"
)

func NotifyEvent(event NotificationType, sinks []*chan NotificationType) {
	for _, s := range sinks {
		select {
		case *s <- event:
		default:
			log.Warn().Any("type", event).Msg("Notification could not be send due to full channel")
		}
	}
}
