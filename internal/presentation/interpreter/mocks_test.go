// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package interpreter

import (
	"context"

	"github.com/rdimidov/kvstore/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// newMockhandler creates a new instance of mockhandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockhandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockhandler {
	mock := &mockhandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// mockhandler is an autogenerated mock type for the handler type
type mockhandler struct {
	mock.Mock
}

type mockhandler_Expecter struct {
	mock *mock.Mock
}

func (_m *mockhandler) EXPECT() *mockhandler_Expecter {
	return &mockhandler_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function for the type mockhandler
func (_mock *mockhandler) Delete(ctx context.Context, key domain.Key) error {
	ret := _mock.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, domain.Key) error); ok {
		r0 = returnFunc(ctx, key)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// mockhandler_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockhandler_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx
//   - key
func (_e *mockhandler_Expecter) Delete(ctx interface{}, key interface{}) *mockhandler_Delete_Call {
	return &mockhandler_Delete_Call{Call: _e.mock.On("Delete", ctx, key)}
}

func (_c *mockhandler_Delete_Call) Run(run func(ctx context.Context, key domain.Key)) *mockhandler_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.Key))
	})
	return _c
}

func (_c *mockhandler_Delete_Call) Return(err error) *mockhandler_Delete_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *mockhandler_Delete_Call) RunAndReturn(run func(ctx context.Context, key domain.Key) error) *mockhandler_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function for the type mockhandler
func (_mock *mockhandler) Get(ctx context.Context, key domain.Key) (*domain.Entry, error) {
	ret := _mock.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *domain.Entry
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, domain.Key) (*domain.Entry, error)); ok {
		return returnFunc(ctx, key)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, domain.Key) *domain.Entry); ok {
		r0 = returnFunc(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Entry)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, domain.Key) error); ok {
		r1 = returnFunc(ctx, key)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// mockhandler_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockhandler_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx
//   - key
func (_e *mockhandler_Expecter) Get(ctx interface{}, key interface{}) *mockhandler_Get_Call {
	return &mockhandler_Get_Call{Call: _e.mock.On("Get", ctx, key)}
}

func (_c *mockhandler_Get_Call) Run(run func(ctx context.Context, key domain.Key)) *mockhandler_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.Key))
	})
	return _c
}

func (_c *mockhandler_Get_Call) Return(entry *domain.Entry, err error) *mockhandler_Get_Call {
	_c.Call.Return(entry, err)
	return _c
}

func (_c *mockhandler_Get_Call) RunAndReturn(run func(ctx context.Context, key domain.Key) (*domain.Entry, error)) *mockhandler_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function for the type mockhandler
func (_mock *mockhandler) Set(ctx context.Context, key domain.Key, value domain.Value) error {
	ret := _mock.Called(ctx, key, value)

	if len(ret) == 0 {
		panic("no return value specified for Set")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, domain.Key, domain.Value) error); ok {
		r0 = returnFunc(ctx, key, value)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// mockhandler_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type mockhandler_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - ctx
//   - key
//   - value
func (_e *mockhandler_Expecter) Set(ctx interface{}, key interface{}, value interface{}) *mockhandler_Set_Call {
	return &mockhandler_Set_Call{Call: _e.mock.On("Set", ctx, key, value)}
}

func (_c *mockhandler_Set_Call) Run(run func(ctx context.Context, key domain.Key, value domain.Value)) *mockhandler_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.Key), args[2].(domain.Value))
	})
	return _c
}

func (_c *mockhandler_Set_Call) Return(err error) *mockhandler_Set_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *mockhandler_Set_Call) RunAndReturn(run func(ctx context.Context, key domain.Key, value domain.Value) error) *mockhandler_Set_Call {
	_c.Call.Return(run)
	return _c
}
