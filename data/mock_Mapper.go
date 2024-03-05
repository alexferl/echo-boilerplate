// Code generated by mockery v2.42.0. DO NOT EDIT.

package data

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	mongo "go.mongodb.org/mongo-driver/mongo"

	options "go.mongodb.org/mongo-driver/mongo/options"
)

// MockMapper is an autogenerated mock type for the Mapper type
type MockMapper struct {
	mock.Mock
}

type MockMapper_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMapper) EXPECT() *MockMapper_Expecter {
	return &MockMapper_Expecter{mock: &_m.Mock}
}

// Aggregate provides a mock function with given fields: ctx, pipeline, results, opts
func (_m *MockMapper) Aggregate(ctx context.Context, pipeline mongo.Pipeline, results interface{}, opts ...*options.AggregateOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, pipeline, results)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Aggregate")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, mongo.Pipeline, interface{}, ...*options.AggregateOptions) (interface{}, error)); ok {
		return rf(ctx, pipeline, results, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, mongo.Pipeline, interface{}, ...*options.AggregateOptions) interface{}); ok {
		r0 = rf(ctx, pipeline, results, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, mongo.Pipeline, interface{}, ...*options.AggregateOptions) error); ok {
		r1 = rf(ctx, pipeline, results, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMapper_Aggregate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Aggregate'
type MockMapper_Aggregate_Call struct {
	*mock.Call
}

// Aggregate is a helper method to define mock.On call
//   - ctx context.Context
//   - pipeline mongo.Pipeline
//   - results interface{}
//   - opts ...*options.AggregateOptions
func (_e *MockMapper_Expecter) Aggregate(ctx interface{}, pipeline interface{}, results interface{}, opts ...interface{}) *MockMapper_Aggregate_Call {
	return &MockMapper_Aggregate_Call{Call: _e.mock.On("Aggregate",
		append([]interface{}{ctx, pipeline, results}, opts...)...)}
}

func (_c *MockMapper_Aggregate_Call) Run(run func(ctx context.Context, pipeline mongo.Pipeline, results interface{}, opts ...*options.AggregateOptions)) *MockMapper_Aggregate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.AggregateOptions, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(*options.AggregateOptions)
			}
		}
		run(args[0].(context.Context), args[1].(mongo.Pipeline), args[2].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockMapper_Aggregate_Call) Return(_a0 interface{}, _a1 error) *MockMapper_Aggregate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMapper_Aggregate_Call) RunAndReturn(run func(context.Context, mongo.Pipeline, interface{}, ...*options.AggregateOptions) (interface{}, error)) *MockMapper_Aggregate_Call {
	_c.Call.Return(run)
	return _c
}

// Collection provides a mock function with given fields: name
func (_m *MockMapper) Collection(name string) Mapper {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for Collection")
	}

	var r0 Mapper
	if rf, ok := ret.Get(0).(func(string) Mapper); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Mapper)
		}
	}

	return r0
}

// MockMapper_Collection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Collection'
type MockMapper_Collection_Call struct {
	*mock.Call
}

// Collection is a helper method to define mock.On call
//   - name string
func (_e *MockMapper_Expecter) Collection(name interface{}) *MockMapper_Collection_Call {
	return &MockMapper_Collection_Call{Call: _e.mock.On("Collection", name)}
}

func (_c *MockMapper_Collection_Call) Run(run func(name string)) *MockMapper_Collection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockMapper_Collection_Call) Return(_a0 Mapper) *MockMapper_Collection_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMapper_Collection_Call) RunAndReturn(run func(string) Mapper) *MockMapper_Collection_Call {
	_c.Call.Return(run)
	return _c
}

// Count provides a mock function with given fields: ctx, filter, opts
func (_m *MockMapper) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Count")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.CountOptions) (int64, error)); ok {
		return rf(ctx, filter, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.CountOptions) int64); ok {
		r0 = rf(ctx, filter, opts...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, ...*options.CountOptions) error); ok {
		r1 = rf(ctx, filter, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMapper_Count_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Count'
type MockMapper_Count_Call struct {
	*mock.Call
}

// Count is a helper method to define mock.On call
//   - ctx context.Context
//   - filter interface{}
//   - opts ...*options.CountOptions
func (_e *MockMapper_Expecter) Count(ctx interface{}, filter interface{}, opts ...interface{}) *MockMapper_Count_Call {
	return &MockMapper_Count_Call{Call: _e.mock.On("Count",
		append([]interface{}{ctx, filter}, opts...)...)}
}

func (_c *MockMapper_Count_Call) Run(run func(ctx context.Context, filter interface{}, opts ...*options.CountOptions)) *MockMapper_Count_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.CountOptions, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(*options.CountOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockMapper_Count_Call) Return(_a0 int64, _a1 error) *MockMapper_Count_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMapper_Count_Call) RunAndReturn(run func(context.Context, interface{}, ...*options.CountOptions) (int64, error)) *MockMapper_Count_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: ctx, filter, results, opts
func (_m *MockMapper) Find(ctx context.Context, filter interface{}, results interface{}, opts ...*options.FindOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter, results)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.FindOptions) (interface{}, error)); ok {
		return rf(ctx, filter, results, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.FindOptions) interface{}); ok {
		r0 = rf(ctx, filter, results, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}, ...*options.FindOptions) error); ok {
		r1 = rf(ctx, filter, results, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMapper_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type MockMapper_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ctx context.Context
//   - filter interface{}
//   - results interface{}
//   - opts ...*options.FindOptions
func (_e *MockMapper_Expecter) Find(ctx interface{}, filter interface{}, results interface{}, opts ...interface{}) *MockMapper_Find_Call {
	return &MockMapper_Find_Call{Call: _e.mock.On("Find",
		append([]interface{}{ctx, filter, results}, opts...)...)}
}

func (_c *MockMapper_Find_Call) Run(run func(ctx context.Context, filter interface{}, results interface{}, opts ...*options.FindOptions)) *MockMapper_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.FindOptions, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(*options.FindOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), args[2].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockMapper_Find_Call) Return(_a0 interface{}, _a1 error) *MockMapper_Find_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMapper_Find_Call) RunAndReturn(run func(context.Context, interface{}, interface{}, ...*options.FindOptions) (interface{}, error)) *MockMapper_Find_Call {
	_c.Call.Return(run)
	return _c
}

// FindOne provides a mock function with given fields: ctx, filter, result, opts
func (_m *MockMapper) FindOne(ctx context.Context, filter interface{}, result interface{}, opts ...*options.FindOneOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter, result)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindOne")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.FindOneOptions) (interface{}, error)); ok {
		return rf(ctx, filter, result, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.FindOneOptions) interface{}); ok {
		r0 = rf(ctx, filter, result, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}, ...*options.FindOneOptions) error); ok {
		r1 = rf(ctx, filter, result, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMapper_FindOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOne'
type MockMapper_FindOne_Call struct {
	*mock.Call
}

// FindOne is a helper method to define mock.On call
//   - ctx context.Context
//   - filter interface{}
//   - result interface{}
//   - opts ...*options.FindOneOptions
func (_e *MockMapper_Expecter) FindOne(ctx interface{}, filter interface{}, result interface{}, opts ...interface{}) *MockMapper_FindOne_Call {
	return &MockMapper_FindOne_Call{Call: _e.mock.On("FindOne",
		append([]interface{}{ctx, filter, result}, opts...)...)}
}

func (_c *MockMapper_FindOne_Call) Run(run func(ctx context.Context, filter interface{}, result interface{}, opts ...*options.FindOneOptions)) *MockMapper_FindOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.FindOneOptions, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(*options.FindOneOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), args[2].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockMapper_FindOne_Call) Return(_a0 interface{}, _a1 error) *MockMapper_FindOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMapper_FindOne_Call) RunAndReturn(run func(context.Context, interface{}, interface{}, ...*options.FindOneOptions) (interface{}, error)) *MockMapper_FindOne_Call {
	_c.Call.Return(run)
	return _c
}

// FindOneAndUpdate provides a mock function with given fields: ctx, filter, update, result, opts
func (_m *MockMapper) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, result interface{}, opts ...*options.FindOneAndUpdateOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter, update, result)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindOneAndUpdate")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) (interface{}, error)); ok {
		return rf(ctx, filter, update, result, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) interface{}); ok {
		r0 = rf(ctx, filter, update, result, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) error); ok {
		r1 = rf(ctx, filter, update, result, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMapper_FindOneAndUpdate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOneAndUpdate'
type MockMapper_FindOneAndUpdate_Call struct {
	*mock.Call
}

// FindOneAndUpdate is a helper method to define mock.On call
//   - ctx context.Context
//   - filter interface{}
//   - update interface{}
//   - result interface{}
//   - opts ...*options.FindOneAndUpdateOptions
func (_e *MockMapper_Expecter) FindOneAndUpdate(ctx interface{}, filter interface{}, update interface{}, result interface{}, opts ...interface{}) *MockMapper_FindOneAndUpdate_Call {
	return &MockMapper_FindOneAndUpdate_Call{Call: _e.mock.On("FindOneAndUpdate",
		append([]interface{}{ctx, filter, update, result}, opts...)...)}
}

func (_c *MockMapper_FindOneAndUpdate_Call) Run(run func(ctx context.Context, filter interface{}, update interface{}, result interface{}, opts ...*options.FindOneAndUpdateOptions)) *MockMapper_FindOneAndUpdate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.FindOneAndUpdateOptions, len(args)-4)
		for i, a := range args[4:] {
			if a != nil {
				variadicArgs[i] = a.(*options.FindOneAndUpdateOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), args[2].(interface{}), args[3].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockMapper_FindOneAndUpdate_Call) Return(_a0 interface{}, _a1 error) *MockMapper_FindOneAndUpdate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMapper_FindOneAndUpdate_Call) RunAndReturn(run func(context.Context, interface{}, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) (interface{}, error)) *MockMapper_FindOneAndUpdate_Call {
	_c.Call.Return(run)
	return _c
}

// FindOneById provides a mock function with given fields: ctx, id, result, opts
func (_m *MockMapper) FindOneById(ctx context.Context, id string, result interface{}, opts ...*options.FindOneOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, id, result)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindOneById")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, ...*options.FindOneOptions) (interface{}, error)); ok {
		return rf(ctx, id, result, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, ...*options.FindOneOptions) interface{}); ok {
		r0 = rf(ctx, id, result, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, interface{}, ...*options.FindOneOptions) error); ok {
		r1 = rf(ctx, id, result, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMapper_FindOneById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOneById'
type MockMapper_FindOneById_Call struct {
	*mock.Call
}

// FindOneById is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
//   - result interface{}
//   - opts ...*options.FindOneOptions
func (_e *MockMapper_Expecter) FindOneById(ctx interface{}, id interface{}, result interface{}, opts ...interface{}) *MockMapper_FindOneById_Call {
	return &MockMapper_FindOneById_Call{Call: _e.mock.On("FindOneById",
		append([]interface{}{ctx, id, result}, opts...)...)}
}

func (_c *MockMapper_FindOneById_Call) Run(run func(ctx context.Context, id string, result interface{}, opts ...*options.FindOneOptions)) *MockMapper_FindOneById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.FindOneOptions, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(*options.FindOneOptions)
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockMapper_FindOneById_Call) Return(_a0 interface{}, _a1 error) *MockMapper_FindOneById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMapper_FindOneById_Call) RunAndReturn(run func(context.Context, string, interface{}, ...*options.FindOneOptions) (interface{}, error)) *MockMapper_FindOneById_Call {
	_c.Call.Return(run)
	return _c
}

// FindOneByIdAndUpdate provides a mock function with given fields: ctx, id, update, result, opts
func (_m *MockMapper) FindOneByIdAndUpdate(ctx context.Context, id string, update interface{}, result interface{}, opts ...*options.FindOneAndUpdateOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, id, update, result)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindOneByIdAndUpdate")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) (interface{}, error)); ok {
		return rf(ctx, id, update, result, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) interface{}); ok {
		r0 = rf(ctx, id, update, result, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) error); ok {
		r1 = rf(ctx, id, update, result, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMapper_FindOneByIdAndUpdate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOneByIdAndUpdate'
type MockMapper_FindOneByIdAndUpdate_Call struct {
	*mock.Call
}

// FindOneByIdAndUpdate is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
//   - update interface{}
//   - result interface{}
//   - opts ...*options.FindOneAndUpdateOptions
func (_e *MockMapper_Expecter) FindOneByIdAndUpdate(ctx interface{}, id interface{}, update interface{}, result interface{}, opts ...interface{}) *MockMapper_FindOneByIdAndUpdate_Call {
	return &MockMapper_FindOneByIdAndUpdate_Call{Call: _e.mock.On("FindOneByIdAndUpdate",
		append([]interface{}{ctx, id, update, result}, opts...)...)}
}

func (_c *MockMapper_FindOneByIdAndUpdate_Call) Run(run func(ctx context.Context, id string, update interface{}, result interface{}, opts ...*options.FindOneAndUpdateOptions)) *MockMapper_FindOneByIdAndUpdate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.FindOneAndUpdateOptions, len(args)-4)
		for i, a := range args[4:] {
			if a != nil {
				variadicArgs[i] = a.(*options.FindOneAndUpdateOptions)
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}), args[3].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockMapper_FindOneByIdAndUpdate_Call) Return(_a0 interface{}, _a1 error) *MockMapper_FindOneByIdAndUpdate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMapper_FindOneByIdAndUpdate_Call) RunAndReturn(run func(context.Context, string, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) (interface{}, error)) *MockMapper_FindOneByIdAndUpdate_Call {
	_c.Call.Return(run)
	return _c
}

// GetNextSequence provides a mock function with given fields: ctx, name
func (_m *MockMapper) GetNextSequence(ctx context.Context, name string) (*Sequence, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for GetNextSequence")
	}

	var r0 *Sequence
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*Sequence, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *Sequence); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Sequence)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMapper_GetNextSequence_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNextSequence'
type MockMapper_GetNextSequence_Call struct {
	*mock.Call
}

// GetNextSequence is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
func (_e *MockMapper_Expecter) GetNextSequence(ctx interface{}, name interface{}) *MockMapper_GetNextSequence_Call {
	return &MockMapper_GetNextSequence_Call{Call: _e.mock.On("GetNextSequence", ctx, name)}
}

func (_c *MockMapper_GetNextSequence_Call) Run(run func(ctx context.Context, name string)) *MockMapper_GetNextSequence_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockMapper_GetNextSequence_Call) Return(_a0 *Sequence, _a1 error) *MockMapper_GetNextSequence_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMapper_GetNextSequence_Call) RunAndReturn(run func(context.Context, string) (*Sequence, error)) *MockMapper_GetNextSequence_Call {
	_c.Call.Return(run)
	return _c
}

// InsertOne provides a mock function with given fields: ctx, document, opts
func (_m *MockMapper) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, document)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for InsertOne")
	}

	var r0 *mongo.InsertOneResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)); ok {
		return rf(ctx, document, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.InsertOneOptions) *mongo.InsertOneResult); ok {
		r0 = rf(ctx, document, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.InsertOneResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, ...*options.InsertOneOptions) error); ok {
		r1 = rf(ctx, document, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMapper_InsertOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertOne'
type MockMapper_InsertOne_Call struct {
	*mock.Call
}

// InsertOne is a helper method to define mock.On call
//   - ctx context.Context
//   - document interface{}
//   - opts ...*options.InsertOneOptions
func (_e *MockMapper_Expecter) InsertOne(ctx interface{}, document interface{}, opts ...interface{}) *MockMapper_InsertOne_Call {
	return &MockMapper_InsertOne_Call{Call: _e.mock.On("InsertOne",
		append([]interface{}{ctx, document}, opts...)...)}
}

func (_c *MockMapper_InsertOne_Call) Run(run func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions)) *MockMapper_InsertOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.InsertOneOptions, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(*options.InsertOneOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockMapper_InsertOne_Call) Return(_a0 *mongo.InsertOneResult, _a1 error) *MockMapper_InsertOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMapper_InsertOne_Call) RunAndReturn(run func(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)) *MockMapper_InsertOne_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateOne provides a mock function with given fields: ctx, filter, update, opts
func (_m *MockMapper) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter, update)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOne")
	}

	var r0 *mongo.UpdateResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)); ok {
		return rf(ctx, filter, update, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) *mongo.UpdateResult); ok {
		r0 = rf(ctx, filter, update, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.UpdateResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) error); ok {
		r1 = rf(ctx, filter, update, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMapper_UpdateOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateOne'
type MockMapper_UpdateOne_Call struct {
	*mock.Call
}

// UpdateOne is a helper method to define mock.On call
//   - ctx context.Context
//   - filter interface{}
//   - update interface{}
//   - opts ...*options.UpdateOptions
func (_e *MockMapper_Expecter) UpdateOne(ctx interface{}, filter interface{}, update interface{}, opts ...interface{}) *MockMapper_UpdateOne_Call {
	return &MockMapper_UpdateOne_Call{Call: _e.mock.On("UpdateOne",
		append([]interface{}{ctx, filter, update}, opts...)...)}
}

func (_c *MockMapper_UpdateOne_Call) Run(run func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions)) *MockMapper_UpdateOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.UpdateOptions, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(*options.UpdateOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), args[2].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockMapper_UpdateOne_Call) Return(_a0 *mongo.UpdateResult, _a1 error) *MockMapper_UpdateOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMapper_UpdateOne_Call) RunAndReturn(run func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)) *MockMapper_UpdateOne_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateOneById provides a mock function with given fields: ctx, id, document, opts
func (_m *MockMapper) UpdateOneById(ctx context.Context, id string, document interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, id, document)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOneById")
	}

	var r0 *mongo.UpdateResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)); ok {
		return rf(ctx, id, document, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, ...*options.UpdateOptions) *mongo.UpdateResult); ok {
		r0 = rf(ctx, id, document, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.UpdateResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, interface{}, ...*options.UpdateOptions) error); ok {
		r1 = rf(ctx, id, document, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMapper_UpdateOneById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateOneById'
type MockMapper_UpdateOneById_Call struct {
	*mock.Call
}

// UpdateOneById is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
//   - document interface{}
//   - opts ...*options.UpdateOptions
func (_e *MockMapper_Expecter) UpdateOneById(ctx interface{}, id interface{}, document interface{}, opts ...interface{}) *MockMapper_UpdateOneById_Call {
	return &MockMapper_UpdateOneById_Call{Call: _e.mock.On("UpdateOneById",
		append([]interface{}{ctx, id, document}, opts...)...)}
}

func (_c *MockMapper_UpdateOneById_Call) Run(run func(ctx context.Context, id string, document interface{}, opts ...*options.UpdateOptions)) *MockMapper_UpdateOneById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.UpdateOptions, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(*options.UpdateOptions)
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockMapper_UpdateOneById_Call) Return(_a0 *mongo.UpdateResult, _a1 error) *MockMapper_UpdateOneById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMapper_UpdateOneById_Call) RunAndReturn(run func(context.Context, string, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)) *MockMapper_UpdateOneById_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockMapper creates a new instance of MockMapper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMapper(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMapper {
	mock := &MockMapper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}