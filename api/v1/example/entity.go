package example

import "time"

type Example struct {
	ID             int64     `db:"id" json:"id"`
	OrganizationID int64     `db:"organization_id" json:"organization_id"`
	Name           string    `db:"name" json:"name"`
	Cover          string    `db:"cover" json:"cover"`
	Notes          string    `db:"notes" json:"notes"`
	Description    string    `db:"description" json:"description"`
	Description2   string    `db:"description2" json:"description2"`
	Style          string    `db:"style" json:"style"`
	Type           string    `db:"type" json:"type"`
	Room           string    `db:"room" json:"room"`
	ExampleType    int       `db:"example_type" json:"example_type"`
	FinderUserName string    `db:"finder_user_name" json:"finder_user_name"`
	FeedID         string    `db:"feed_id" json:"feed_id"`
	Priority       int       `db:"priority" json:"priority"`
	Building       string    `db:"building" json:"building"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
