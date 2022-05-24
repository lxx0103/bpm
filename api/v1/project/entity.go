package project

import "time"

type Project struct {
	ID              int64     `db:"id" json:"id"`
	OrganizationID  int64     `db:"organization_id" json:"organization_id"`
	TemplateID      int64     `db:"template_id" json:"template_id"`
	ClientID        int64     `db:"client_id" json:"client_id"`
	Name            string    `db:"name" json:"name"`
	Type            int       `db:"type" json:"type"`
	Location        string    `db:"location" json:"location"`
	Longitude       float64   `db:"longitude" json:"longitude"`
	Latitude        float64   `db:"latitude" json:"latitude"`
	CheckinDistance int       `db:"checkin_distance" json:"checkin_distance"`
	Status          int       `db:"status" json:"status"`
	Created         time.Time `db:"created" json:"created"`
	CreatedBy       string    `db:"created_by" json:"created_by"`
	Updated         time.Time `db:"updated" json:"updated"`
	UpdatedBy       string    `db:"updated_by" json:"updated_by"`
}
