package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	app "github.com/alexferl/echo-boilerplate"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/models"
	_ "github.com/alexferl/echo-boilerplate/testing"
)

type TaskHandlerTestSuite struct {
	suite.Suite
	svc         *handlers.MockTaskService
	server      *server.Server
	user        *models.User
	accessToken []byte
}

func (s *TaskHandlerTestSuite) SetupTest() {
	svc := handlers.NewMockTaskService(s.T())
	h := handlers.NewTaskHandler(openapi.NewHandler(), svc)
	user := models.NewUser("test@example.com", "test")
	user.Id = "1"
	user.Create(user.Id)
	access, _, _ := user.Login()

	s.svc = svc
	s.server = app.NewTestServer(h)
	s.user = user
	s.accessToken = access
}

func TestTaskHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(TaskHandlerTestSuite))
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Get_200() {
	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	task := models.NewTask()
	task.Create(s.user.Id)
	task.CreatedBy = s.user

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil)

	s.server.ServeHTTP(resp, req)

	var result models.TaskResponse
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), s.user.Id, result.CreatedBy.Id)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Get_401() {
	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Get_404() {
	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(nil, data.ErrNoDocuments)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
	assert.Equal(s.T(), "task not found", result.Message)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Get_410() {
	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	task := models.NewTask()
	task.Create(s.user.Id)
	task.Delete(s.user.Id)
	task.CreatedBy = s.user
	task.DeletedBy = s.user

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusGone, resp.Code)
	assert.Equal(s.T(), "task was deleted", result.Message)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Update_200() {
	title := "My Edited Task"
	payload := &handlers.UpdateTaskRequest{
		Title: &title,
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	task := models.NewTask()
	task.Create(s.user.Id)
	task.CreatedBy = s.user

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil)

	task.Update(s.user.Id)
	task.UpdatedBy = s.user

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(task, nil)

	s.server.ServeHTTP(resp, req)

	var result models.TaskResponse
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), title, result.Title)
	assert.Equal(s.T(), s.user.Id, result.UpdatedBy.Id)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Update_401() {
	req := httptest.NewRequest(http.MethodPut, "/tasks/1", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Update_403() {
	title := "My Edited Task"
	payload := &handlers.UpdateTaskRequest{
		Title: &title,
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	user := &models.User{Model: models.Model{Id: "2"}}
	task := models.NewTask()
	task.Create(user.Id)
	task.CreatedBy = user

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusForbidden, resp.Code)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Update_404() {
	title := "My Edited Task"
	payload := &handlers.UpdateTaskRequest{
		Title: &title,
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(nil, data.ErrNoDocuments)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
	assert.Equal(s.T(), "task not found", result.Message)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Update_410() {
	title := "My Edited Task"
	payload := &handlers.UpdateTaskRequest{
		Title: &title,
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	task := models.NewTask()
	task.Create(s.user.Id)
	task.Delete(s.user.Id)
	task.CreatedBy = s.user
	task.DeletedBy = s.user

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusGone, resp.Code)
	assert.Equal(s.T(), "task was deleted", result.Message)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Update_422() {
	payload := &handlers.UpdateTaskRequest{}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnprocessableEntity, resp.Code)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Transition_200() {
	t := true
	payload := &handlers.TransitionTaskRequest{
		Completed: &t,
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1/transition", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	task := models.NewTask()
	task.Create(s.user.Id)
	task.CreatedBy = s.user

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil)

	task.Complete(s.user.Id)
	task.CompletedBy = s.user

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(task, nil)

	s.server.ServeHTTP(resp, req)

	var result models.TaskResponse
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), s.user.Id, result.CompletedBy.Id)
	assert.True(s.T(), result.Completed)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Transition_401() {
	req := httptest.NewRequest(http.MethodPut, "/tasks/1/transition", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Transition_404() {
	t := true
	payload := &handlers.TransitionTaskRequest{
		Completed: &t,
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1/transition", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(nil, data.ErrNoDocuments)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
	assert.Equal(s.T(), "task not found", result.Message)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Transition_410() {
	t := true
	payload := &handlers.TransitionTaskRequest{
		Completed: &t,
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1/transition", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	task := models.NewTask()
	task.Create(s.user.Id)
	task.Delete(s.user.Id)
	task.CreatedBy = s.user
	task.DeletedBy = s.user

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusGone, resp.Code)
	assert.Equal(s.T(), "task was deleted", result.Message)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Transition_422() {
	payload := &handlers.TransitionTaskRequest{}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1/transition", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnprocessableEntity, resp.Code)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Delete_200() {
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	task := models.NewTask()
	task.Create(s.user.Id)
	task.CreatedBy = s.user

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil)
	s.svc.EXPECT().
		Delete(mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusNoContent, resp.Code)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Delete_403() {
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	user := &models.User{Model: models.Model{Id: "2"}}
	task := models.NewTask()
	task.Create(user.Id)
	task.CreatedBy = user

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusForbidden, resp.Code)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Delete_404() {
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	task := models.NewTask()
	task.Create(s.user.Id)
	task.CreatedBy = s.user

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(nil, data.ErrNoDocuments)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
	assert.Equal(s.T(), "task not found", result.Message)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Delete_410() {
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	task := models.NewTask()
	task.Create(s.user.Id)
	task.CreatedBy = s.user
	task.Delete(s.user.Id)
	task.DeletedBy = s.user

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusGone, resp.Code)
	assert.Equal(s.T(), "task was deleted", result.Message)
}

func createTasks(num int, user *models.User) models.Tasks {
	result := make(models.Tasks, 0)

	for i := 1; i <= num; i++ {
		newTask := models.NewTask()
		newTask.Create(user.Id)
		newTask.CreatedBy = user
		result = append(result, *newTask)
	}

	return result
}

func (s *TaskHandlerTestSuite) TestTaskHandler_List_200() {
	req := httptest.NewRequest(http.MethodGet, "/tasks?per_page=1&page=2", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	num := 10
	tasks := createTasks(num, s.user)

	s.svc.EXPECT().
		Find(mock.Anything, mock.Anything).
		Return(int64(num), tasks, nil)

	s.server.ServeHTTP(resp, req)

	var result handlers.ListTasksResponse
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	h := resp.Header()
	link := `<http://example.com/tasks?per_page=1&page=3>; rel=next, ` +
		`<http://example.com/tasks?per_page=1&page=10>; rel=last, ` +
		`<http://example.com/tasks?per_page=1&page=1>; rel=first, ` +
		`<http://example.com/tasks?per_page=1&page=1>; rel=prev`

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), num, len(result.Tasks))
	assert.Equal(s.T(), "2", h.Get("X-Page"))
	assert.Equal(s.T(), "1", h.Get("X-Per-Page"))
	assert.Equal(s.T(), "10", h.Get("X-Total"))
	assert.Equal(s.T(), "10", h.Get("X-Total-Pages"))
	assert.Equal(s.T(), "3", h.Get("X-Next-Page"))
	assert.Equal(s.T(), "1", h.Get("X-Prev-Page"))
	assert.Equal(s.T(), link, h.Get("Link"))
}

func (s *TaskHandlerTestSuite) TestTaskHandler_List_401() {
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Create_200() {
	payload := &handlers.CreateTaskRequest{Title: "Test"}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Create(mock.Anything, mock.Anything, mock.Anything).
		Return(&models.Task{Model: models.Model{CreatedBy: s.user}, Title: payload.Title}, nil)

	s.server.ServeHTTP(resp, req)

	var result models.TaskResponse
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), payload.Title, result.Title)
	assert.Equal(s.T(), s.user.Id, result.CreatedBy.Id)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Create_401() {
	payload := &handlers.CreateTaskRequest{Title: "Test"}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Create_422() {
	payload := &handlers.CreateTaskRequest{}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result models.TaskResponse
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusUnprocessableEntity, resp.Code)
}
