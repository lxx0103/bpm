package event

import (
	"bpm/api/v1/component"
	"bpm/core/database"
	"errors"
	"fmt"
	"math"
	"strings"
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
	NewCheckin(int64, NewCheckin) error
	GetCheckinList(CheckinFilter, int64) (int, *[]CheckinResponse, error)
	//WX API
	GetAssignedEvent(AssignedEventFilter, int64, int64, int64) (*[]MyEvent, error)
	GetAssignedAudit(AssignedAuditFilter, int64, int64, int64) (*[]MyEvent, error)
	GetProjectEvent(MyEventFilter) (*[]MyEvent, error)
	SaveEvent(int64, SaveEventInfo) error
	AuditEvent(int64, AuditEventInfo) error
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
	audits, err := query.GetAuditsByEventID(event.ID)
	if err != nil {
		return nil, err
	}
	event.Audit = audits
	return event, err
}

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
	if oldEvent.Assignable == 1 {
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
		}
	}
	oldEvent.NeedAudit = info.NeedAudit
	err = repo.DeleteEventAudit(eventID, info.User)
	if err != nil {
		return nil, err
	}
	err = repo.CreateEventAudit(eventID, info.AuditType, info.AuditTo, info.User)
	if err != nil {
		return nil, err
	}
	err = repo.UpdateEvent(eventID, *oldEvent, info.User)
	if err != nil {
		return nil, err
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
	audits, err := repo.GetAuditsByEventID(eventID)
	if err != nil {
		return nil, err
	}
	event.Audit = audits
	tx.Commit()
	return event, err
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

func (s *eventService) GetProjectEvent(filter MyEventFilter) (*[]MyEvent, error) {
	db := database.InitMySQL()
	query := NewEventQuery(db)
	myEvents, err := query.GetProjectEvent(filter)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return nil, err
		}
	}
	return myEvents, err
}

func (s *eventService) SaveEvent(eventID int64, info SaveEventInfo) error {
	db := database.InitMySQL()
	query := NewEventQuery(db)
	active, err := query.CheckActive(eventID)
	if err != nil {
		return err
	}
	if !active {
		msg := "?????????????????????"
		return errors.New(msg)
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewEventRepository(tx)
	componentRepo := component.NewComponentRepository(tx)
	event, err := repo.GetEventByID(eventID, 0)
	if err != nil {
		return err
	}
	if event.Status != 1 && event.Status != 3 {
		msg := "??????????????????"
		return errors.New(msg)
	}
	assignExist, err := repo.CheckAssign(eventID, info.UserID, info.PositionID)
	if err != nil {
		return err
	}
	if assignExist == 0 {
		msg := "????????????????????????"
		return errors.New(msg)
	}
	components, err := componentRepo.GetComponentByEventID(eventID)
	if err != nil {
		return err
	}
	for i := 0; i < len(info.Components); i++ {
		componentInfo := info.Components[i]
		for j := 0; j < len(*components); j++ {
			toUpdate := (*components)[j]
			if componentInfo.ID == toUpdate.ID {
				if toUpdate.Patterns != "" {
					patternArr := strings.Split(toUpdate.Patterns, "|")
					if len(patternArr) != 2 {
						msg := "??????????????????"
						return errors.New(msg)
					}
					switch patternArr[0] {
					case "oneof":
						valid := false
						valueArr := strings.Split(patternArr[1], ";")
						for k := 0; k < len(valueArr); k++ {
							if componentInfo.Value == valueArr[k] {
								valid = true
							}
						}
						if !valid {
							msg := toUpdate.Name + "??????????????????"
							return errors.New(msg)
						}
					case "mul":
						valid := false
						inputArr := strings.Split(componentInfo.Value, ";")
						valueArr := strings.Split(patternArr[1], ";")
						for k := 0; k < len(inputArr); k++ {
							valid = false
							for l := 0; l < len(valueArr); l++ {
								if inputArr[k] == valueArr[l] {
									valid = true
								}
							}
						}
						if !valid {
							msg := toUpdate.Name + "??????????????????"
							return errors.New(msg)
						}
					default:
						msg := toUpdate.Name + "??????????????????"
						return errors.New(msg)
					}

				}
				err := componentRepo.SaveComponent(toUpdate.ID, componentInfo.Value, info.User)
				if err != nil {
					return err
				}
				break
			}
		}
	}
	requiredCount, err := componentRepo.CheckRequired(eventID)
	if err != nil {
		return err
	}
	if requiredCount != 0 {
		msg := "???" + fmt.Sprintf("%v", requiredCount) + "??????????????????"
		return errors.New(msg)
	}
	err = repo.CompleteEvent(eventID, info.User)
	if err != nil {
		return err
	}
	fmt.Println("needAudit:", event.NeedAudit)
	if event.NeedAudit == 2 {
		err = repo.AuditEvent(eventID, true, "SYSTEM", "????????????")
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func (s *eventService) AuditEvent(eventID int64, info AuditEventInfo) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewEventRepository(tx)
	event, err := repo.GetEventByID(eventID, 0)
	if err != nil {
		return err
	}
	if event.Status != 2 {
		msg := "?????????????????????"
		return errors.New(msg)
	}
	assignExist, err := repo.CheckAudit(eventID, info.UserID, info.PositionID)
	if err != nil {
		return err
	}
	if assignExist == 0 {
		msg := "????????????????????????"
		return errors.New(msg)
	}
	approved := true
	if info.Result != 1 {
		approved = false
	}
	err = repo.AuditEvent(eventID, approved, info.User, info.Content)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
func (s *eventService) GetAssignedAudit(filter AssignedAuditFilter, userID int64, positionID int64, organizationID int64) (*[]MyEvent, error) {
	var activeEvents []MyEvent
	db := database.InitMySQL()
	query := NewEventQuery(db)
	assignedAudit, err := query.GetAssignedAudit(userID, positionID)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(assignedAudit); i++ {
		active, err := query.CheckActive(assignedAudit[i])
		if err != nil {
			return nil, err
		}
		if !active {
			continue
		}
		activeEvent, err := query.GetAssignedAuditByID(assignedAudit[i], filter.Status)
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

func (s *eventService) NewCheckin(eventID int64, info NewCheckin) error {
	db := database.InitMySQL()
	query := NewEventQuery(db)
	active, err := query.CheckActive(eventID)
	if err != nil {
		return err
	}
	if !active {
		msg := "?????????????????????"
		return errors.New(msg)
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewEventRepository(tx)
	event, err := repo.GetEventByID(eventID, info.OrganizationID)
	if err != nil {
		msg := "??????????????????"
		return errors.New(msg)
	}
	if event.NeedCheckin == 0 {
		msg := "?????????????????????"
		return errors.New(msg)
	}
	if event.Status != 1 && event.Status != 3 {
		msg := "??????????????????"
		return errors.New(msg)
	}
	assignExist, err := repo.CheckAssign(eventID, info.UserID, info.PositionID)
	if err != nil {
		msg := "??????????????????"
		return errors.New(msg)
	}
	if assignExist == 0 {
		msg := "????????????????????????"
		return errors.New(msg)
	}
	projectLongitude, projectLatitude, projectDistance, err := repo.GetProjectLocation(event.ProjectID, info.OrganizationID)
	if err != nil {
		msg := "??????????????????"
		return errors.New(msg)
	}
	distance := getDistance(projectLatitude, projectLongitude, info.Latitude, info.Longitude)
	if projectDistance < distance && projectDistance != 0 {
		msg := "?????????????????????:" + fmt.Sprintf("%v", distance) + "???"
		return errors.New(msg)
	}
	info.Distance = distance
	checkinExist, err := repo.CheckCheckin(eventID, info.UserID)
	if err != nil {
		return err
	}
	if checkinExist >= 2 {
		msg := "????????????????????????"
		return errors.New(msg)
	}
	if checkinExist == 1 {
		info.CheckinType = 2
	} else {
		info.CheckinType = 1
	}
	err = repo.doCheckin(eventID, info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func getDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) int {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lng1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lng2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	dist := 2 * r * math.Asin(math.Sqrt(h))
	return int(dist)
}
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func (s *eventService) GetCheckinList(filter CheckinFilter, organizationID int64) (int, *[]CheckinResponse, error) {
	if organizationID != 0 && organizationID != filter.OrganizationID {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewEventQuery(db)
	count, err := query.GetCheckinCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetCheckinList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}
