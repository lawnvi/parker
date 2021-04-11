package model

import "time"

type BaseModel struct {
	ID uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `sql:"index"`
}
