package domain_events

import (
	"log"

	"github.com/hashemihossein/GO-CQRS/examples/domain"
)

type UserCreatedEvent struct {
	User *domain.User
}

func (event *UserCreatedEvent) Handle() error {
	// handling the event
	log.Println("User created: ", event.User)
	return nil
}
