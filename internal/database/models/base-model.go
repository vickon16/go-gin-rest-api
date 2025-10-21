package models

import "time"

type BaseModel struct {
	CreatedAt *time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}
