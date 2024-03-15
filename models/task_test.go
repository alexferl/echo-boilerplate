package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestTask(t *testing.T) {
	task := NewTask()
	assert.NotEqual(t, "", task.Id)

	id := "1"
	task.Create(id)
	assert.Equal(t, id, task.CreatedBy.(*Ref).Id)
	assert.NotNil(t, task.CreatedAt)

	task.Update(id)
	assert.Equal(t, id, task.UpdatedBy.(*Ref).Id)
	assert.NotNil(t, task.UpdatedAt)

	task.Delete(id)
	assert.Equal(t, id, task.DeletedBy.(*Ref).Id)
	assert.NotNil(t, task.DeletedAt)

	task.Complete(id)
	assert.Equal(t, id, task.CompletedBy.(*Ref).Id)
	assert.NotNil(t, task.CompletedAt)

	task.Incomplete()
	assert.Nil(t, task.CompletedBy)
	assert.Nil(t, task.CompletedAt)
}

func TestTask_CustomBSON(t *testing.T) {
	task := NewTask()
	id := "1"
	user := NewUser("test@example.com", "test")
	user.Id = id
	task.CompletedBy = user
	task.CreatedBy = user
	task.DeletedBy = user
	task.UpdatedBy = user

	b, _ := bson.Marshal(task)

	var m Task
	_ = bson.Unmarshal(b, &m)

	assert.Equal(t, id, m.CompletedBy.(*User).Id)
	assert.Equal(t, id, m.CreatedBy.(*User).Id)
	assert.Equal(t, id, m.DeletedBy.(*User).Id)
	assert.Equal(t, id, m.UpdatedBy.(*User).Id)
	assert.IsType(t, &UserRef{}, m.Response().CompletedBy)
	assert.IsType(t, &UserRef{}, m.Response().CreatedBy)
	assert.IsType(t, &UserRef{}, m.Response().UpdatedBy)
}

func TestTasks(t *testing.T) {
	user := NewUser("test@example.com", "test")

	task1 := NewTask()
	task1.Create("1")
	task1.CreatedBy = user

	task2 := NewTask()
	task2.Create("1")
	task2.CreatedBy = user

	tasks := Tasks{*task1, *task2}
	resp := tasks.Response()

	assert.Len(t, resp.Tasks, 2)
}
