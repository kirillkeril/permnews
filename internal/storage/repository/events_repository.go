package repository

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"permnews/internal/storage/filters"
	"permnews/internal/storage/models"
	"strings"
)

type EventsRepository struct {
	storage Storage
}

func NewEventsRepository(storage Storage) EventsRepository {
	return EventsRepository{storage: storage}
}

func (r EventsRepository) GetAll(filter *filters.EventFilter) ([]*models.Event, error) {
	events := make([]*models.Event, 0)
	db := r.storage.GetDb()
	if filter == nil {
		rows, err := db.Query(`select * from events`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			ev := &models.Event{}
			err = rows.Scan(&ev.Id, &ev.Title, &ev.Notice, &ev.Source, &ev.Image, &ev.Link, &ev.Category)
			events = append(events, ev)
		}
		if err != nil {
			return nil, err
		}
	} else {
		q := `select * from events`
		filterList := make([]string, 0)
		if filter.Id != "" {
			filterList = append(filterList, fmt.Sprintf(`id = '%v'`, filter.Id))
		}
		if filter.Source != "" {
			filterList = append(filterList, fmt.Sprintf(`source = '%v'`, filter.Source))
		}
		if filter.Category != "" {
			filterList = append(filterList, fmt.Sprintf(`category = '%v'`, filter.Category))
		}
		if filter.Title != "" {
			filterList = append(filterList, fmt.Sprintf(`title like '%v'`, filter.Title))
		}
		f := strings.Join(filterList, " and ")
		if len(f) > 0 {
			q += ` where ` + f
		}
		if filter.Skip > 0 {
			q += fmt.Sprintf(` offset %v`, filter.Skip)
		}
		if filter.Limit > 0 {
			q += fmt.Sprintf(` limit %v`, filter.Limit)
		}
		log.Println(q)
		rows, err := db.Query(q)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			ev := &models.Event{}
			err = rows.Scan(&ev.Id, &ev.Title, &ev.Notice, &ev.Source, &ev.Image, &ev.Link, &ev.Category)
			events = append(events, ev)
		}
		if err != nil {
			return nil, err
		}
	}
	return events, nil
}

func (r EventsRepository) Create(event *models.Event) error {
	db := r.storage.GetDb()
	_, err := db.Exec(`
			insert into events (id, title, notice, source, image, link, category) 
			values ($1, $2, $3, $4, $5, $6, $7)`, uuid.New(), event.Title, event.Notice, event.Source, event.Image, event.Link, event.Category)
	if err != nil {
		return err
	}
	return nil
}
