package event

import (
	"bpm/core/database"
)

type eventService struct {
}

func NewEventService() EventService {
	return &eventService{}
}

// EventService represents a service for managing events.
type EventService interface {
	//Event Management
	GetEventByID(int64) (*Event, error)
	NewEvent(EventNew) (*Event, error)
	GetEventList(EventFilter) (int, *[]Event, error)
	UpdateEvent(int64, EventNew) (*Event, error)
}

func (s *eventService) GetEventByID(id int64) (*Event, error) {
	db := database.InitMySQL()
	query := NewEventQuery(db)
	event, err := query.GetEventByID(id)
	return event, err
}

func (s *eventService) NewEvent(info EventNew) (*Event, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewEventRepository(tx)
	eventID, err := repo.CreateEvent(info)
	if err != nil {
		return nil, err
	}
	event, err := repo.GetEventByID(eventID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return event, err
}

func (s *eventService) GetEventList(filter EventFilter) (int, *[]Event, error) {
	db := database.InitMySQL()
	query := NewEventQuery(db)
	count, err := query.GetEventCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetEventList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *eventService) UpdateEvent(eventID int64, info EventNew) (*Event, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewEventRepository(tx)
	_, err = repo.UpdateEvent(eventID, info)
	if err != nil {
		return nil, err
	}
	event, err := repo.GetEventByID(eventID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return event, err
}
