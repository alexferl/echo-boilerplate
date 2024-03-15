package models

import (
	"time"

	"github.com/rs/xid"
)

// Model is the base model for all models.
// NewModel should be used unless you know what you're doing.
type Model struct {
	Id        string     `bson:"id"`
	CreatedAt *time.Time `bson:"created_at"`
	CreatedBy any        `bson:"created_by"`
	DeletedAt *time.Time `bson:"deleted_at"`
	DeletedBy any        `bson:"deleted_by"`
	UpdatedAt *time.Time `bson:"updated_at"`
	UpdatedBy any        `bson:"updated_by"`
}

// Ref is a reference to another document.
type Ref struct {
	Id string `json:"id" bson:"id"`
}

func NewModel() *Model {
	return &Model{Id: xid.New().String()}
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
