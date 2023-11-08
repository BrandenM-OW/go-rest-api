package models

import (
	"time"
)

type File struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
	Name        string
	Owner       string
	ContentType string
	Size        int64
}
