package main

import "log"

type NotificationsService struct {
	// Channel to receive transactions that matches the from address
	fromCh chan Transaction

	// Channel to receive transactions that matches the to address
	toCh chan Transaction
}

func NewNotificationsService() *NotificationsService {
	return &NotificationsService{
		fromCh: make(chan Transaction),
		toCh:   make(chan Transaction),
	}
}

func (ns NotificationsService) Start() {
	for {
		select {
		case transaction := <-ns.fromCh:
			log.Printf("%v sent a transaction: %v\n", transaction.From, transaction.Hash)
		case transaction := <-ns.toCh:
			log.Printf("%v received a transaction: %v\n", transaction.To, transaction.Hash)
		}
	}
}
