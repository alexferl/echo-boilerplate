package data

import (
	"time"

	"github.com/rs/xid"
)

type Model struct {
	Id        string     `json:"id" bson:"id"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	CreatedBy string     `json:"created_by" bson:"created_by"`
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
	DeletedBy string     `json:"deleted_by" bson:"deleted_by"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedBy string     `json:"updated_by" bson:"updated_by"`
}

func (m *Model) Create(id string) {
	t := time.Now()
	m.CreatedAt = &t
	m.CreatedBy = id
}

func (m *Model) Delete(id string) {
	t := time.Now()
	m.DeletedAt = &t
	m.DeletedBy = id
}

func (m *Model) Update(id string) {
	t := time.Now()
	m.UpdatedAt = &t
	m.UpdatedBy = id
}

func NewModel() *Model {
	return &Model{
		Id: xid.New().String(),
	}
}
