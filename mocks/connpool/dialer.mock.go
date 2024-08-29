// Code generated by mockery v2.45.0. DO NOT EDIT.

package mockconnpool

import (
	context "context"
	net "net"

	mock "github.com/stretchr/testify/mock"
)

// MockDialer is an autogenerated mock type for the Dialer type
type MockDialer struct {
	mock.Mock
}

type MockDialer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDialer) EXPECT() *MockDialer_Expecter {
	return &MockDialer_Expecter{mock: &_m.Mock}
}

// DialContext provides a mock function with given fields: ctx
func (_m *MockDialer) DialContext(ctx context.Context) (net.Conn, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DialContext")
	}

	var r0 net.Conn
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (net.Conn, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) net.Conn); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(net.Conn)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDialer_DialContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DialContext'
type MockDialer_DialContext_Call struct {
	*mock.Call
}

// DialContext is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockDialer_Expecter) DialContext(ctx interface{}) *MockDialer_DialContext_Call {
	return &MockDialer_DialContext_Call{Call: _e.mock.On("DialContext", ctx)}
}

func (_c *MockDialer_DialContext_Call) Run(run func(ctx context.Context)) *MockDialer_DialContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockDialer_DialContext_Call) Return(_a0 net.Conn, _a1 error) *MockDialer_DialContext_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDialer_DialContext_Call) RunAndReturn(run func(context.Context) (net.Conn, error)) *MockDialer_DialContext_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDialer creates a new instance of MockDialer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDialer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDialer {
	mock := &MockDialer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
