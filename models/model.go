package models

import (
	"time"

	"github.com/rs/xid"
)

type Model struct {
	Id        string     `json:"id" bson:"id"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	CreatedBy any        `json:"created_by" bson:"created_by"`
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
	DeletedBy any        `json:"deleted_by" bson:"deleted_by"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedBy any        `json:"updated_by" bson:"updated_by"`
}

//func (m *Model) MarshalBSON() ([]byte, error) {
//	log.Debug().Interface("MARSH", m).Send()
//	var r Ref
//	err := util.DocToStruct(m.CreatedBy.(primitive.D), &r)
//	if err != nil {
//		log.Error().Err(err).Msg("MARSH FAIL")
//		return nil, err
//	}
//	m.CreatedBy = r
//
//	log.Debug().Interface("RRRRR", m).Send()
//
//	return bson.Marshal(r)
//}

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

func NewModel() Model {
	return Model{Id: xid.New().String()}
}
