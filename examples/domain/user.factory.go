package domain

import "github.com/hashemihossein/GO-CQRS/pkg/aggregate"

func NewUser(Username, Password, DateOfBirth string) *User {
	return &User{
		AggregateRoot: aggregate.GetNewAggregateRoot(),
		Username:      Username,
		Password:      Password,
		DateOfBirth:   DateOfBirth,
	}
}
