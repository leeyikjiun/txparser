package main

import (
	"testing"
)

func TestInMemorySubscriberStore_HasSubscriber(t *testing.T) {
	imss := NewInMemorySubscriberStore()
	address := "some_address"
	if imss.HasSubscriber(address) {
		t.Fatalf("Subscriber should not contain %v", address)
	}

	imss.AddSubscriber(address)
	if !imss.HasSubscriber(address) {
		t.Fatalf("Subscriber should contain %v", address)
	}
}
