package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	m := NewModel()
	assert.NotEqual(t, "", m.Id)

	id := "1"
	m.Create(id)
	assert.Equal(t, id, m.CreatedBy.(*Ref).Id)
	assert.NotNil(t, m.CreatedAt)

	m.Update(id)
	assert.Equal(t, id, m.UpdatedBy.(*Ref).Id)
	assert.NotNil(t, m.UpdatedAt)

	m.Delete(id)
	assert.Equal(t, id, m.DeletedBy.(*Ref).Id)
	assert.NotNil(t, m.DeletedAt)
}
