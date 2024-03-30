package services

import (
	"permnews/internal/models/entity"
	"permnews/internal/storage/filters"
	"permnews/internal/storage/models"
	"permnews/internal/storage/repository"
)

type EventsRepository interface {
	GetAll(filter *filters.EventFilter) ([]*models.Event, error)
	Create(event *models.Event) error
}

type EventsService struct {
	repository repository.EventsRepository
}

func NewEventsService(repository repository.EventsRepository) *EventsService {
	return &EventsService{repository: repository}
}

func (es *EventsService) GetAll(filter *filters.EventFilter) ([]*entity.Event, error) {
	entities := make([]*entity.Event, 0)
	mds, err := es.repository.GetAll(filter)
	if err != nil {
		return nil, err
	}
	for _, m := range mds {
		e := entity.Event{
			Id:       m.Id,
			Title:    m.Title,
			Notice:   m.Notice,
			Source:   m.Source,
			Image:    m.Image,
			Category: m.Category,
		}
		entities = append(entities, &e)
	}
	return entities, nil
}

func (es *EventsService) Create(e *models.Event) error {
	m := models.Event{
		Id:       e.Id,
		Title:    e.Title,
		Notice:   e.Notice,
		Source:   e.Source,
		Image:    e.Image,
		Category: e.Category,
	}
	err := es.repository.Create(&m)
	if err != nil {
		return err
	}
	return nil
}
