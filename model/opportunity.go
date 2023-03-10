package model

import (
	"time"

	"gorm.io/gorm"
)

type Opportunity struct {
	Id              uint      `gorm:"column:id" json:"id"`
	Code            string    `gorm:"column:code" json:"code"`
	ClientCode      string    `gorm:"column:client_code" json:"client_code"`
	PicEmail        string    `gorm:"column:pic_email" json:"pic_email"`
	OpportunityName string    `gorm:"column:opportunity_name" json:"opportunity_name"`
	Description     string    `gorm:"column:description" json:"description"`
	SalesEmail      string    `gorm:"column:sales_email" json:"sales_email"`
	Status          string    `gorm:"column:status" json:"status"`
	LastModified    time.Time `gorm:"column:last_modified" json:"last_modified" time_format:"2006-01-02 15:04:05"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (op *Opportunity) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	op.CreatedAt = now
	op.UpdatedAt = now
	return nil
}

func (op *Opportunity) BeforeUpdate(tx *gorm.DB) error {
	op.UpdatedAt = time.Now()
	return nil
}
