package service

import (
	"github.com/jcuga/golongpoll"
	"github.com/rs/zerolog/log"
)

func (uc *useCaseImpl) OnLongpollStart() <-chan *golongpoll.Event {
	return uc.eventChan
}

func (uc *useCaseImpl) OnPublish(event *golongpoll.Event) {
	log.Info().Str("id", event.ID.String()).Str("Category", event.Category).Msg("long poll shutdown")
}

func (uc *useCaseImpl) OnShutdown() {
	log.Info().Msg("long poll shutdown")
}
