package main

type SubscriberStore interface {
	AddSubscriber(address string)
	HasSubscriber(address string) bool
}
