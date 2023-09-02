package main

type InMemorySubscriberStore struct {
	subscribers map[string]bool
}

func NewInMemorySubscriberStore() *InMemorySubscriberStore {
	return &InMemorySubscriberStore{
		subscribers: make(map[string]bool),
	}
}

func (imss *InMemorySubscriberStore) AddSubscriber(address string) {
	imss.subscribers[address] = true
}

func (imss *InMemorySubscriberStore) HasSubscriber(address string) bool {
	return imss.subscribers[address]
}
