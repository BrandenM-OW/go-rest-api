package models

import (
	"net/url"

	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title   string
	Content string
}

type FileItem struct {
	gorm.Model
	Name            string
	url             url.URL
	lastModified    string
	cacheExpiration string
	CachedUrl       string
	created         string
	owner           string
}
