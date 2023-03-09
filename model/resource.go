package model

import (
	"time"

	"gorm.io/gorm"
)

type Resource struct {
	Id              uint    `gorm:"column:id" json:"id"`
	OpportunityId   uint    `gorm:"foreignKey:Opportunity;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"opportunity_id"`
	Qty             int64     `gorm:"column:qty" json:"qty" binding:"required"`
	Position        string    `gorm:"column:position" json:"position" binding:"required"`
	Level           string    `gorm:"column:level" json:"level" binding:"required"`
	Ctc             float64   `gorm:"column:ctc" json:"ctc"`
	ProjectDuration int64     `gorm:"column:project_duration" json:"project_duration" binding:"required"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (rsc *Resource) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	rsc.CreatedAt = now
	rsc.UpdatedAt = now
	return nil
}

func (rsc *Resource) BeforeUpdate(tx *gorm.DB) error {
	rsc.UpdatedAt = time.Now()
	return nil
}
