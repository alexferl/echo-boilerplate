package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexferl/echo-openapi"
	api "github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/services"
)

type TaskHandlerTestSuite struct {
	suite.Suite
	svc         *handlers.MockTaskService
	userSvc     *handlers.MockUserService
	server      *api.Server
	user        *models.User
	accessToken []byte
}

func (s *TaskHandlerTestSuite) SetupTest() {
	svc := handlers.NewMockTaskService(s.T())
	userSvc := handlers.NewMockUserService(s.T())
	patSvc := handlers.NewMockPersonalAccessTokenService(s.T())
	h := handlers.NewTaskHandler(openapi.NewHandler(), svc)
	user := getUser()
	access, _, _ := user.Login()

	s.svc = svc
	s.userSvc = userSvc
	s.server = getServer(userSvc, patSvc, h)
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

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result models.TaskResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusOK, resp.Code)
	s.Assert().Equal(s.user.Id, result.CreatedBy.Id)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_401() {
	testCases := []struct {
		method   string
		endpoint string
	}{
		{http.MethodPost, "/tasks"},
		{http.MethodGet, "/tasks"},
		{http.MethodGet, "/tasks/1"},
		{http.MethodPatch, "/tasks/1"},
		{http.MethodPut, "/tasks/1/transition"},
		{http.MethodDelete, "/tasks/1"},
	}
	for _, tc := range testCases {
		s.T().Run(fmt.Sprintf("%s_%s", tc.method, tc.endpoint), func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.endpoint, nil)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			s.server.ServeHTTP(resp, req)

			var result echo.HTTPError
			_ = json.Unmarshal(resp.Body.Bytes(), &result)

			s.Assert().Equal(http.StatusUnauthorized, resp.Code)
			s.Assert().Equal("token invalid", result.Message)
		})
	}
}

func (s *TaskHandlerTestSuite) TestTaskHandler_404() {
	title := "My Edited Task"
	updateBody := &handlers.UpdateTaskRequest{
		Title: &title,
	}

	t := true
	transitionBody := &handlers.TransitionTaskRequest{
		Completed: &t,
	}

	testCases := []struct {
		method   string
		endpoint string
		body     any
	}{
		{http.MethodGet, "/tasks/1", nil},
		{http.MethodPatch, "/tasks/1", updateBody},
		{http.MethodPut, "/tasks/1/transition", transitionBody},
		{http.MethodDelete, "/tasks/1", nil},
	}
	for _, tc := range testCases {
		s.T().Run(fmt.Sprintf("%s_%s", tc.method, tc.endpoint), func(t *testing.T) {
			b, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(tc.method, tc.endpoint, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
			resp := httptest.NewRecorder()

			// middleware
			s.userSvc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(s.user, nil).Once()

			s.svc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(nil, &services.Error{
					Kind:    services.NotExist,
					Message: services.ErrTaskNotFound.Error(),
				}).Once()

			s.server.ServeHTTP(resp, req)

			var result echo.HTTPError
			_ = json.Unmarshal(resp.Body.Bytes(), &result)

			s.Assert().Equal(http.StatusNotFound, resp.Code)
			s.Assert().Equal(services.ErrTaskNotFound.Error(), result.Message)
		})
	}
}

func (s *TaskHandlerTestSuite) TestTaskHandler_410() {
	title := "My Edited Task"
	updateBody := &handlers.UpdateTaskRequest{
		Title: &title,
	}

	t := true
	transitionBody := &handlers.TransitionTaskRequest{
		Completed: &t,
	}

	testCases := []struct {
		method   string
		endpoint string
		body     any
	}{
		{http.MethodGet, "/tasks/1", nil},
		{http.MethodPatch, "/tasks/1", updateBody},
		{http.MethodPut, "/tasks/1/transition", transitionBody},
		{http.MethodDelete, "/tasks/1", nil},
	}
	for _, tc := range testCases {
		s.T().Run(fmt.Sprintf("%s_%s", tc.method, tc.endpoint), func(t *testing.T) {
			b, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(tc.method, tc.endpoint, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
			resp := httptest.NewRecorder()

			// middleware
			s.userSvc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(s.user, nil).Once()

			s.svc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(nil, &services.Error{
					Kind:    services.Deleted,
					Message: services.ErrTaskDeleted.Error(),
				}).Once()

			s.server.ServeHTTP(resp, req)

			var result echo.HTTPError
			_ = json.Unmarshal(resp.Body.Bytes(), &result)

			s.Assert().Equal(http.StatusGone, resp.Code)
			s.Assert().Equal(services.ErrTaskDeleted.Error(), result.Message)
		})
	}
}

func (s *TaskHandlerTestSuite) TestTaskHandler_422() {
	testCases := []struct {
		method   string
		endpoint string
	}{
		{http.MethodPost, "/tasks"},
		{http.MethodPatch, "/tasks/1"},
		{http.MethodPut, "/tasks/1/transition"},
	}
	for _, tc := range testCases {
		s.T().Run(fmt.Sprintf("%s_%s", tc.method, tc.endpoint), func(t *testing.T) {
			body := &handlers.UpdateTaskRequest{}
			b, _ := json.Marshal(body)
			req := httptest.NewRequest(tc.method, tc.endpoint, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
			resp := httptest.NewRecorder()

			// middleware
			s.userSvc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(s.user, nil).Once()

			s.server.ServeHTTP(resp, req)

			s.Assert().Equal(http.StatusUnprocessableEntity, resp.Code)
		})
	}
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Update_200() {
	title := "My Edited Task"
	payload := &handlers.UpdateTaskRequest{
		Title: &title,
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	task := models.NewTask()
	task.Create(s.user.Id)
	task.CreatedBy = s.user

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil).Once()

	task.Update(s.user.Id)
	task.UpdatedBy = s.user

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(task, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result models.TaskResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusOK, resp.Code)
	s.Assert().Equal(title, result.Title)
	s.Assert().Equal(s.user.Id, result.UpdatedBy.Id)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Update_403() {
	title := "My Edited Task"
	payload := &handlers.UpdateTaskRequest{
		Title: &title,
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	user := &models.User{Model: &models.Model{Id: "2"}}
	task := models.NewTask()
	task.Create(user.Id)
	task.CreatedBy = user

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil).Once()

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusForbidden, resp.Code)
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

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil).Once()

	task.Complete(s.user.Id)
	task.CompletedBy = s.user

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(task, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result models.TaskResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusOK, resp.Code)
	s.Assert().Equal(s.user.Id, result.CompletedBy.Id)
	s.Assert().True(result.Completed)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Delete_200() {
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	task := models.NewTask()
	task.Create(s.user.Id)
	task.CreatedBy = s.user

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil).Once()

	s.svc.EXPECT().
		Delete(mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusNoContent, resp.Code)
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Delete_403() {
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	user := &models.User{Model: &models.Model{Id: "2"}}
	task := models.NewTask()
	task.Create(user.Id)
	task.CreatedBy = user

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(task, nil).Once()

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusForbidden, resp.Code)
}

func createTasks(num int, user *models.User) models.Tasks {
	result := make(models.Tasks, 0)

	for range num {
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

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Find(mock.Anything, mock.Anything).
		Return(int64(num), tasks, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result models.TasksResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	h := resp.Header()
	link := `<http://example.com/tasks?per_page=1&page=3>; rel=next, ` +
		`<http://example.com/tasks?per_page=1&page=10>; rel=last, ` +
		`<http://example.com/tasks?per_page=1&page=1>; rel=first, ` +
		`<http://example.com/tasks?per_page=1&page=1>; rel=prev`

	s.Assert().Equal(http.StatusOK, resp.Code)
	s.Assert().Equal(num, len(result.Tasks))
	s.Assert().Equal("2", h.Get("X-Page"))
	s.Assert().Equal("1", h.Get("X-Per-Page"))
	s.Assert().Equal("10", h.Get("X-Total"))
	s.Assert().Equal("10", h.Get("X-Total-Pages"))
	s.Assert().Equal("3", h.Get("X-Next-Page"))
	s.Assert().Equal("1", h.Get("X-Prev-Page"))
	s.Assert().Equal(link, h.Get("Link"))
}

func (s *TaskHandlerTestSuite) TestTaskHandler_Create_200() {
	payload := &handlers.CreateTaskRequest{Title: "Test"}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Create(mock.Anything, mock.Anything, mock.Anything).
		Return(&models.Task{Model: &models.Model{CreatedBy: s.user}, Title: payload.Title}, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result models.TaskResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusOK, resp.Code)
	s.Assert().Equal(payload.Title, result.Title)
	s.Assert().Equal(s.user.Id, result.CreatedBy.Id)
}
