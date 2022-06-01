package entities

import (
	"time"
)

type Person struct {
	ID        uint32    `json:"id,omitempty" gorm:"primaryKey"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"` // Use unix milli seconds as updating time
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"` // Use unix seconds as creating time`
}
