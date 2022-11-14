// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	data "github.com/alexferl/echo-boilerplate/data"
	mock "github.com/stretchr/testify/mock"

	options "go.mongodb.org/mongo-driver/mongo/options"
)

// Mapper is an autogenerated mock type for the Mapper type
type Mapper struct {
	mock.Mock
}

// Aggregate provides a mock function with given fields: ctx, filter, limit, skip, result, opts
func (_m *Mapper) Aggregate(ctx context.Context, filter interface{}, limit int, skip int, result interface{}, opts ...*options.AggregateOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter, limit, skip, result)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, int, int, interface{}, ...*options.AggregateOptions) interface{}); ok {
		r0 = rf(ctx, filter, limit, skip, result, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, int, int, interface{}, ...*options.AggregateOptions) error); ok {
		r1 = rf(ctx, filter, limit, skip, result, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Collection provides a mock function with given fields: name
func (_m *Mapper) Collection(name string) data.Mapper {
	ret := _m.Called(name)

	var r0 data.Mapper
	if rf, ok := ret.Get(0).(func(string) data.Mapper); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(data.Mapper)
		}
	}

	return r0
}

// Count provides a mock function with given fields: ctx, filter, opts
func (_m *Mapper) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.CountOptions) int64); ok {
		r0 = rf(ctx, filter, opts...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, ...*options.CountOptions) error); ok {
		r1 = rf(ctx, filter, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: ctx, filter, result, opts
func (_m *Mapper) Find(ctx context.Context, filter interface{}, result interface{}, opts ...*options.FindOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter, result)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.FindOptions) interface{}); ok {
		r0 = rf(ctx, filter, result, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}, ...*options.FindOptions) error); ok {
		r1 = rf(ctx, filter, result, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOne provides a mock function with given fields: ctx, filter, result, opts
func (_m *Mapper) FindOne(ctx context.Context, filter interface{}, result interface{}, opts ...*options.FindOneOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter, result)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.FindOneOptions) interface{}); ok {
		r0 = rf(ctx, filter, result, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}, ...*options.FindOneOptions) error); ok {
		r1 = rf(ctx, filter, result, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneById provides a mock function with given fields: ctx, id, result, opts
func (_m *Mapper) FindOneById(ctx context.Context, id string, result interface{}, opts ...*options.FindOneOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, id, result)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, ...*options.FindOneOptions) interface{}); ok {
		r0 = rf(ctx, id, result, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, interface{}, ...*options.FindOneOptions) error); ok {
		r1 = rf(ctx, id, result, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: ctx, document, result, opts
func (_m *Mapper) Insert(ctx context.Context, document interface{}, result interface{}, opts ...*options.InsertOneOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, document, result)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.InsertOneOptions) interface{}); ok {
		r0 = rf(ctx, document, result, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}, ...*options.InsertOneOptions) error); ok {
		r1 = rf(ctx, document, result, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, filter, update, result, opts
func (_m *Mapper) Update(ctx context.Context, filter interface{}, update interface{}, result interface{}, opts ...*options.UpdateOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter, update, result)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, interface{}, ...*options.UpdateOptions) interface{}); ok {
		r0 = rf(ctx, filter, update, result, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}, interface{}, ...*options.UpdateOptions) error); ok {
		r1 = rf(ctx, filter, update, result, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateById provides a mock function with given fields: ctx, id, document, result, opts
func (_m *Mapper) UpdateById(ctx context.Context, id string, document interface{}, result interface{}, opts ...*options.UpdateOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, id, document, result)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, interface{}, ...*options.UpdateOptions) interface{}); ok {
		r0 = rf(ctx, id, document, result, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, interface{}, interface{}, ...*options.UpdateOptions) error); ok {
		r1 = rf(ctx, id, document, result, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upsert provides a mock function with given fields: ctx, filter, update, result, opts
func (_m *Mapper) Upsert(ctx context.Context, filter interface{}, update interface{}, result interface{}, opts ...*options.FindOneAndUpdateOptions) (interface{}, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter, update, result)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) interface{}); ok {
		r0 = rf(ctx, filter, update, result, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) error); ok {
		r1 = rf(ctx, filter, update, result, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMapper interface {
	mock.TestingT
	Cleanup(func())
}

// NewMapper creates a new instance of Mapper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMapper(t mockConstructorTestingTNewMapper) *Mapper {
	mock := &Mapper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
