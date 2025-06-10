package model

import "time"

type Task struct {
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
