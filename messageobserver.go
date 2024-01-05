package gojourney

import (
	"errors"
	"time"

	"github.com/kiwioneone/gojourney/discord"
)

type MessageFilterFunc func(msg *discord.Message) (bool, error)

type CommonMessageObserver struct {
	filterFunc MessageFilterFunc
	ch         chan *discord.Message
	errCh      chan error
}

func NewCommonMessageObserver(filterFunc MessageFilterFunc) *CommonMessageObserver {
	return &CommonMessageObserver{
		filterFunc: filterFunc,
		ch:         make(chan *discord.Message),
		errCh:      make(chan error),
	}
}

func (o *CommonMessageObserver) Observe(msg *discord.Message) {
	hit, err := o.filterFunc(msg)
	if hit {
		o.ch <- msg
	}

	if err != nil {
		o.errCh <- err
	}
}

func (o *CommonMessageObserver) WaitMsg(timeout time.Duration) (*discord.Message, error) {
	defer close(o.ch)
	defer close(o.errCh)
	select {
	case result := <-o.ch:
		return result, nil
	case err := <-o.errCh:
		return nil, err
	case <-time.After(timeout):
		return nil, errors.New("timeout")
	}
}
