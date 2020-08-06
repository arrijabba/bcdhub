package database

import (
	"time"

	"github.com/lib/pq"
)

// DApp model
type DApp struct {
	ID               uint           `gorm:"primary_key" json:"-"`
	CreatedAt        time.Time      `json:"-"`
	UpdatedAt        time.Time      `json:"-"`
	DeletedAt        *time.Time     `json:"-"`
	Name             string         `json:"name"`
	ShortDescription string         `json:"short_description"`
	FullDescription  string         `json:"full_description"`
	Version          string         `json:"version"`
	License          string         `json:"license"`
	Authors          pq.StringArray `gorm:"type:varchar(128)[]" json:"authors"`
	SocialLinks      pq.StringArray `gorm:"type:varchar(1024)[]" json:"social_links"`
	Interfaces       pq.StringArray `gorm:"type:varchar(64)[]" json:"interfaces"`
	Categories       pq.StringArray `gorm:"type:varchar(32)[]" json:"categories"`

	Pictures  []Picture `json:"pictures"`
	Contracts []Alias   `json:"contracts"`
}

// Picture model
type Picture struct {
	ID     uint   `gorm:"primary_key,AUTO_INCREMENT" json:"-"`
	Link   string `json:"link"`
	Type   string `json:"type"`
	DAppID uint   `json:"-"`
}

// GetDApps -
func (d *db) GetDApps() (dapps []DApp, err error) {
	err = d.ORM.Preload("Pictures").Preload("Contracts").Find(&dapps).Error
	return
}

// GetDApp -
func (d *db) GetDApp(id uint) (dapp DApp, err error) {
	err = d.ORM.Preload("Pictures").Preload("Contracts").Where("id = ?", id).First(&dapp).Error
	return
}
