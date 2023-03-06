package event

import (
	"bpm/core/database"
	"bpm/core/queue"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type EventActiveChanged struct {
	ProjectID int64 `json:"project_id"`
}

func Subscribe(conn *queue.Conn) {
	conn.StartConsumer("UpdateActiveEvent", "EventActiveChanged", UpdateActiveEvent)
}

func UpdateActiveEvent(d amqp.Delivery) bool {
	if d.Body == nil {
		return false
	}
	var EventActiveChanged EventActiveChanged
	err := json.Unmarshal(d.Body, &EventActiveChanged)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "5")
			return false
		}
	}
	err = setEventActive(EventActiveChanged.ProjectID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		return true
	}
}

func setEventActive(projectID int64) error {
	db := database.InitMySQL()
	query := NewEventQuery(db)
	var filter MyEventFilter
	filter.ProjectID = projectID
	filter.Status = "active"
	events, err := query.GetProjectEvent(filter)
	if err != nil {
		fmt.Println(err.Error() + "4")
		return err
	}
	var actives []int64
	for _, event := range *events {
		active, err := query.CheckActive(event.ID)
		if err != nil {
			fmt.Println(err.Error() + "3")
			return err
		}
		if !active {
			continue
		}
		actives = append(actives, event.ID)
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewEventRepository(tx)
	for _, active := range actives {
		err := repo.SetEventActive(active)
		if err != nil {
			return err
		}
	}
	all, completed, err := repo.GetProjectProgress(projectID)
	if err != nil {
		return err
	}
	progress := completed * 100 / all
	err = repo.UpdateProjectProgress(projectID, progress)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
