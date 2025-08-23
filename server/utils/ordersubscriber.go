package utils

import (
	"github.com/codecrafter404/bubble/ent"
	"github.com/google/uuid"
)

type Subscriber struct {
	Id         uuid.UUID
	Reciever   chan *ent.Order
	Processing *int
}

func NewSubscriber() Subscriber {
	return Subscriber{
		Id:         uuid.New(),
		Reciever:   make(chan *ent.Order),
		Processing: nil,
	}
}
