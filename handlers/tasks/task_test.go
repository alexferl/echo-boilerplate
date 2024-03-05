package tasks_test

//func TestHandler_GetTask_200(t *testing.T) {
//	mockModel := tasks.NewMockIModel(t)
//	h := tasks.NewHandler(&mongo.Client{}, openapi.NewHandler(), mockModel)
//	s := app.NewTestServer(h)
//
//	user := users.NewUser("test@example.com", "test")
//	access, _, err := user.Login()
//	assert.NoError(t, err)
//
//	req := httptest.NewRequest(http.MethodGet, "/tasks/id", nil)
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
//	resp := httptest.NewRecorder()
//
//	Model := &tasks.Model{}
//
//	mockModel.EXPECT().
//		Load(mock.Anything, mock.Anything, mock.Anything).
//		Return(Model, nil)
//
//	//mockModel.Mock.
//	//	On("Load", mock.Anything, mock.Anything, mock.Anything).
//	//	Return(&tasks.Model{}, nil).
//	//	On("Read", mock.Anything).
//	//	Return(&tasks.Task{
//	//		Model:       &tasks.Model{Model: &data.Model{Id: "pd"}},
//	//		CreatedBy:   nil,
//	//		DeletedAt:   nil,
//	//		DeletedBy:   nil,
//	//		UpdatedBy:   nil,
//	//		CompletedBy: nil,
//	//	}, nil)
//
//	s.ServeHTTP(resp, req)
//
//	assert.Equal(t, http.StatusOK, resp.Code)
//}
