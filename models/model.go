package models

import (
	"time"

	"github.com/rs/xid"
)

// Model is the base model for all models.
// NewModel should be used unless you know what you're doing.
type Model struct {
	Id        string     `json:"id" bson:"id"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	CreatedBy any        `json:"created_by" bson:"created_by"`
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
	DeletedBy any        `json:"deleted_by" bson:"deleted_by"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedBy any        `json:"updated_by" bson:"updated_by"`
}

// Ref is a reference to another document
type Ref struct {
	Id string `json:"id" bson:"id"`
}

func (m *Model) Create(id string) {
	t := time.Now()
	m.CreatedAt = &t
	m.CreatedBy = &Ref{Id: id}
}

func (m *Model) Delete(id string) {
	t := time.Now()
	m.DeletedAt = &t
	m.DeletedBy = &Ref{Id: id}
}

func (m *Model) Update(id string) {
	t := time.Now()
	m.UpdatedAt = &t
	m.UpdatedBy = &Ref{Id: id}
}

func NewModel() *Model {
	return &Model{Id: xid.New().String()}
}
