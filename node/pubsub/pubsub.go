package pubsub

import (
	"sync"

	"github.com/google/uuid"
)

// Topic is a pubsub topic
// T is the type of the messages
// C() is the channel where the messages are published
type Topic[T any] struct {
	*sync.RWMutex
	c           chan T
	subChannels map[string]Channel[T]
}

func (t *Topic[T]) start() {
	// push messages to all subscribers
	for msg := range t.c {
		t.RLock()
		for _, sub := range t.subChannels {
			sub <- msg
		}
		t.RUnlock()
	}
}

func (t *Topic[T]) C() Channel[T] {
	return t.c
}

func NewTopic[T any]() *Topic[T] {
	t := &Topic[T]{
		RWMutex:     &sync.RWMutex{},
		c:           make(Channel[T]),
		subChannels: make(map[string]Channel[T]),
	}

	go t.start()

	return t
}

// Sub returns a subscriber to the topic
// the subscriber will receive all messages published to the topic
// until it is unsubscribed.
//
// Remember to unsubscribe when you are done with the subscriber.
func (t *Topic[T]) Sub() *Subscriber[T] {
	t.Lock()
	defer t.Unlock()

	sub := &Subscriber[T]{
		topic: t,
		id:    uuid.New().String(),
	}

	t.subChannels[sub.id] = make(Channel[T])
	return sub
}

type Channel[T any] chan T

type Subscriber[T any] struct {
	topic *Topic[T]
	id    string
}

func (s *Subscriber[T]) UnSub() {
	s.topic.Lock()
	defer s.topic.Unlock()

	// close the channel
	close(s.topic.subChannels[s.id])
	// delete the channel from the topic
	delete(s.topic.subChannels, s.id)
}

func (s *Subscriber[T]) C() Channel[T] {
	return s.topic.subChannels[s.id]
}
