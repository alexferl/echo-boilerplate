package tasks_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers/tasks"
	"github.com/alexferl/echo-boilerplate/handlers/users"
)

func TestHandler_GetTask_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	task := tasks.NewTask()
	task.Create(user.Id)
	short := &tasks.TaskWithUsers{
		Id:          task.Id,
		Title:       task.Title,
		IsPrivate:   task.IsPrivate,
		IsCompleted: task.IsCompleted,
		CreatedAt:   task.CreatedAt,
		CreatedBy: &tasks.TaskUser{
			Id:       user.Id,
			Username: user.Username,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"Aggregate",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			[]*tasks.TaskWithUsers{short},
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHandler_GetTask_401(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodGet, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestHandler_GetTask_404(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"Aggregate",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			[]*tasks.TaskWithUsers{},
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestHandler_GetTask_410(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	task := tasks.NewTask()
	task.Create(user.Id)
	task.Delete(user.Id)
	short := &tasks.TaskWithUsers{
		Id:        task.Id,
		DeletedAt: task.DeletedAt,
		DeletedBy: task.DeletedBy,
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"Aggregate",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			[]*tasks.TaskWithUsers{short},
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusGone, resp.Code)
}

func TestHandler_PatchTask_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	payload := &tasks.TaskPatch{
		Title: "My Edited Task",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	task := tasks.NewTask()
	task.Create(user.Id)
	find := &tasks.Task{
		Model: &data.Model{
			CreatedBy: user.Id,
			DeletedAt: nil,
		},
		Title:       payload.Title,
		IsPrivate:   false,
		IsCompleted: false,
		CompletedAt: task.CompletedAt,
		CompletedBy: task.CompletedBy,
	}
	update := &tasks.TaskWithUsers{
		Id:          task.Id,
		Title:       task.Title,
		IsPrivate:   task.IsPrivate,
		IsCompleted: task.IsCompleted,
		CreatedAt:   task.CreatedAt,
		CreatedBy: &tasks.TaskUser{
			Id:       user.Id,
			Username: user.Username,
		},
	}

	req := httptest.NewRequest(http.MethodPatch, "/tasks/id", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			find,
			nil,
		).
		On(
			"UpdateById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			[]*tasks.TaskWithUsers{update},
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHandler_PatchTask_401(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestHandler_PatchTask_403(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	payload := &tasks.TaskPatch{
		Title: "My Edited Task",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	task := tasks.NewTask()
	task.Create(user.Id)
	task.Complete(user.Id)
	find := &tasks.Task{
		Model: &data.Model{
			DeletedBy: "",
		},
		Title:       payload.Title,
		IsPrivate:   false,
		IsCompleted: false,
		CompletedAt: task.CompletedAt,
		CompletedBy: task.CompletedBy,
	}

	req := httptest.NewRequest(http.MethodPatch, "/tasks/id", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			find,
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestHandler_PatchTask_404(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	payload := &tasks.TaskPatch{
		Title: "My Edited Task",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/id", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			nil,
			tasks.ErrTaskNotFound,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestHandler_PatchTask_410(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	payload := &tasks.TaskPatch{
		Title: "My Edited Task",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	task := tasks.NewTask()
	task.Create(user.Id)
	task.Delete(user.Id)
	find := &tasks.Task{
		Model: &data.Model{
			DeletedAt: task.DeletedAt,
		},
		Title:       payload.Title,
		IsPrivate:   false,
		IsCompleted: false,
		CompletedAt: task.CompletedAt,
		CompletedBy: task.CompletedBy,
	}

	req := httptest.NewRequest(http.MethodPatch, "/tasks/id", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			find,
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusGone, resp.Code)
}

func TestHandler_PatchTask_422(t *testing.T) {
	_, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	b := bytes.NewBuffer([]byte(`{"invalid": "invalid"}`))
	req := httptest.NewRequest(http.MethodPatch, "/tasks/id", b)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}
