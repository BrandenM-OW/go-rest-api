package models

import (
	"net/url"
	"time"
)

type File struct {
	ID              uint `gorm:"primarykey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time `gorm:"index"`
	Name            string
	url             url.URL
	lastModified    string
	cacheExpiration string
	CachedUrl       string
	created         string
	owner           string
}
