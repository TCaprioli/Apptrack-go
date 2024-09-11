// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"

	null "github.com/guregu/null/v5"
)

type Application struct {
	ID              int32       `json:"id"`
	JobTitle        string      `json:"job_title"`
	Company         string      `json:"company"`
	Location        null.String `json:"location"`
	ApplicationDate null.Time   `json:"application_date"`
	UserID          int32       `json:"user_id"`
	Status          string      `json:"status"`
	Notes           null.String `json:"notes"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type User struct {
	ID         int32       `json:"id"`
	Email      string      `json:"email"`
	Password   string      `json:"password"`
	GoogleID   null.String `json:"google_id"`
	LinkedInID null.String `json:"linked_in_id"`
	Name       null.String `json:"name"`
	UpdatedAt  time.Time   `json:"updated_at"`
	CreatedAt  time.Time   `json:"created_at"`
}