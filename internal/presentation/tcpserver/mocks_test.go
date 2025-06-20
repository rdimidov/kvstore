// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package tcpserver

import (
	"context"

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

// Execute provides a mock function for the type mockhandler
func (_mock *mockhandler) Execute(context1 context.Context, bytes []byte) []byte {
	ret := _mock.Called(context1, bytes)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 []byte
	if returnFunc, ok := ret.Get(0).(func(context.Context, []byte) []byte); ok {
		r0 = returnFunc(context1, bytes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}
	return r0
}

// mockhandler_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type mockhandler_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - context1
//   - bytes
func (_e *mockhandler_Expecter) Execute(context1 interface{}, bytes interface{}) *mockhandler_Execute_Call {
	return &mockhandler_Execute_Call{Call: _e.mock.On("Execute", context1, bytes)}
}

func (_c *mockhandler_Execute_Call) Run(run func(context1 context.Context, bytes []byte)) *mockhandler_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]byte))
	})
	return _c
}

func (_c *mockhandler_Execute_Call) Return(bytes1 []byte) *mockhandler_Execute_Call {
	_c.Call.Return(bytes1)
	return _c
}

func (_c *mockhandler_Execute_Call) RunAndReturn(run func(context1 context.Context, bytes []byte) []byte) *mockhandler_Execute_Call {
	_c.Call.Return(run)
	return _c
}
