package domain

import "github.com/hashemihossein/GO-CQRS/pkg/aggregate"

type User struct {
	*aggregate.AggregateRoot
	ID          string
	Username    string
	Password    string
	DateOfBirth string
}
