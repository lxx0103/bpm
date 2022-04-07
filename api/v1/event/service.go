package event

import (
	"bpm/core/database"
	"errors"
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
	// NewEvent(EventNew, int64) (*Event, error)
	GetEventList(EventFilter, int64) (int, *[]Event, error)
	UpdateEvent(int64, EventUpdate, int64) (*Event, error)
	DeleteEventByProjectID(int64, int64, string) error
	//WX API
	GetAssignedEvent(AssignedEventFilter, int64, int64, int64) (*[]MyEvent, error)
	GetMyEvent(MyEventFilter, string) (*[]MyEvent, error)
}

func (s *eventService) GetEventByID(id int64) (*Event, error) {
	db := database.InitMySQL()
	query := NewEventQuery(db)
	event, err := query.GetEventByID(id)
	if err != nil {
		return nil, err
	}
	assigns, err := query.GetAssignsByEventID(event.ID)
	if err != nil {
		return nil, err
	}
	event.Assign = assigns
	pres, err := query.GetPresByEventID(event.ID)
	if err != nil {
		return nil, err
	}
	event.PreID = pres
	return event, err
}

// func (s *eventService) NewEvent(info EventNew, organizationID int64) (*Event, error) {
// 	db := database.InitMySQL()
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback()
// 	repo := NewEventRepository(tx)
// 	projectExist, err := repo.CheckProjectExist(info.ProjectID, organizationID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if projectExist == 0 {
// 		msg := "项目不存在"
// 		return nil, errors.New(msg)
// 	}
// 	exist, err := repo.CheckNameExist(info.Name, info.ProjectID, 0)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if exist != 0 {
// 		msg := "事件名称重复"
// 		return nil, errors.New(msg)
// 	}
// 	eventID, err := repo.CreateEvent(info)
// 	if err != nil {
// 		return nil, err
// 	}
// 	event, err := repo.GetEventByID(eventID, organizationID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = repo.CreateEventAssign(eventID, info.AssignType, info.AssignTo, info.User)
// 	if err != nil {
// 		return nil, err
// 	}
// 	assigns, err := repo.GetAssignsByEventID(eventID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	event.Assign = assigns
// 	err = repo.CreateEventPre(eventID, info.PreID, info.User)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pres, err := repo.GetPresByEventID(eventID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	event.PreID = pres
// 	tx.Commit()
// 	return event, err
// }

func (s *eventService) GetEventList(filter EventFilter, organizationID int64) (int, *[]Event, error) {
	db := database.InitMySQL()
	query := NewEventQuery(db)
	count, err := query.GetEventCount(filter, organizationID)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetEventList(filter, organizationID)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *eventService) UpdateEvent(eventID int64, info EventUpdate, organizationID int64) (*Event, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewEventRepository(tx)
	oldEvent, err := repo.GetEventByID(eventID, organizationID)
	if err != nil {
		return nil, err
	}
	if oldEvent.Assignable != 1 {
		msg := "不能修改分配"
		return nil, errors.New(msg)
	} else {
		if info.AssignType != 0 {
			oldEvent.AssignType = info.AssignType
			err = repo.DeleteEventAssign(eventID, info.User)
			if err != nil {
				return nil, err
			}
			err = repo.CreateEventAssign(eventID, info.AssignType, info.AssignTo, info.User)
			if err != nil {
				return nil, err
			}
			err = repo.UpdateEvent(eventID, *oldEvent, info.User)
			if err != nil {
				return nil, err
			}
		}
	}
	event, err := repo.GetEventByID(eventID, organizationID)
	if err != nil {
		return nil, err
	}
	assigns, err := repo.GetAssignsByEventID(eventID)
	if err != nil {
		return nil, err
	}
	event.Assign = assigns
	pres, err := repo.GetPresByEventID(eventID)
	if err != nil {
		return nil, err
	}
	event.PreID = pres
	tx.Commit()
	return event, err
}

func (s *eventService) DeleteEventByProjectID(projectID int64, organizationID int64, user string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewEventRepository(tx)
	err = repo.DeleteEventByProjectID(projectID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *eventService) GetAssignedEvent(filter AssignedEventFilter, userID int64, positionID int64, organizationID int64) (*[]MyEvent, error) {
	var activeEvents []MyEvent
	db := database.InitMySQL()
	query := NewEventQuery(db)
	assigned, err := query.GetAssigned(userID, positionID)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(assigned); i++ {
		active, err := query.CheckActive(assigned[i])
		if err != nil {
			return nil, err
		}
		if !active {
			continue
		}
		activeEvent, err := query.GetAssignedEventByID(assigned[i], filter.Status)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return nil, err
			}
			continue
		}
		activeEvents = append(activeEvents, *activeEvent)
	}
	return &activeEvents, err
}

func (s *eventService) GetMyEvent(filter MyEventFilter, createdBy string) (*[]MyEvent, error) {
	db := database.InitMySQL()
	query := NewEventQuery(db)
	myEvents, err := query.GetMyEvent(filter, createdBy)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return nil, err
		}
	}
	return myEvents, err
}
