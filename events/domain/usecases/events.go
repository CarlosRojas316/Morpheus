package usecases

import (
	"events/application"
	"events/domain/entities"
)

type CreateEvent interface {
	Create(event *application.CreateEventParams) (*entities.Event, error)
}

type FindEvents interface {
	FindAccountEvents(accountId string) ([]*entities.Event, error)
	FindEventById(eventId string) (*entities.Event, error)
	FindAll(state string, month, ageGroup int) ([]entities.Event, error)
}

type UpdateEvents interface {
	UpdateStatus(eventId string, status interface{}) error
	UpdateData(eventId string, data *application.UpdateEventDTO) (*entities.Event, error)
}
