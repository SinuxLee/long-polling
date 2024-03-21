package service

import (
	"context"
	"time"

	"github.com/jcuga/golongpoll"
	"template/internal/store"
)

var _ UseCase = (*useCaseImpl)(nil)

type config interface {
	GetThirdParty() string
}

type UseCase interface {
	Init() error
	NewDistLock(key string) store.DistLock
	Hello(ctx context.Context, name string) (string, error)
	GetLongPoller() *golongpoll.LongpollManager
}

func NewUseCase(d store.Dao, conf config) UseCase {
	return &useCaseImpl{
		dao:       d,
		conf:      conf,
		eventChan: make(chan *golongpoll.Event, 64),
	}
}

type useCaseImpl struct {
	dao       store.Dao
	conf      config
	eventChan chan *golongpoll.Event
	longPoll  *golongpoll.LongpollManager
}

func (uc *useCaseImpl) Init() error {
	longPoll, err := golongpoll.StartLongpoll(golongpoll.Options{
		LoggingEnabled:                 true,
		MaxLongpollTimeoutSeconds:      30,
		EventTimeToLiveSeconds:         600,
		DeleteEventAfterFirstRetrieval: false,
	})
	if err != nil {
		return err
	}

	uc.longPoll = longPoll

	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			longPoll.Publish("libz", "haha")
		}
	}()

	return nil
}

func (uc *useCaseImpl) GetLongPoller() *golongpoll.LongpollManager {
	return uc.longPoll
}

func (uc *useCaseImpl) NewDistLock(key string) store.DistLock {
	return uc.dao.NewDistLock(key)
}
