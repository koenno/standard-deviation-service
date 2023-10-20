// Code generated by mockery v2.35.2. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// Handler is an autogenerated mock type for the Handler type
type Handler struct {
	mock.Mock
}

type Handler_Expecter struct {
	mock *mock.Mock
}

func (_m *Handler) EXPECT() *Handler_Expecter {
	return &Handler_Expecter{mock: &_m.Mock}
}

// ServeHTTP provides a mock function with given fields: _a0, _a1
func (_m *Handler) ServeHTTP(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// Handler_ServeHTTP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ServeHTTP'
type Handler_ServeHTTP_Call struct {
	*mock.Call
}

// ServeHTTP is a helper method to define mock.On call
//   - _a0 http.ResponseWriter
//   - _a1 *http.Request
func (_e *Handler_Expecter) ServeHTTP(_a0 interface{}, _a1 interface{}) *Handler_ServeHTTP_Call {
	return &Handler_ServeHTTP_Call{Call: _e.mock.On("ServeHTTP", _a0, _a1)}
}

func (_c *Handler_ServeHTTP_Call) Run(run func(_a0 http.ResponseWriter, _a1 *http.Request)) *Handler_ServeHTTP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(*http.Request))
	})
	return _c
}

func (_c *Handler_ServeHTTP_Call) Return() *Handler_ServeHTTP_Call {
	_c.Call.Return()
	return _c
}

func (_c *Handler_ServeHTTP_Call) RunAndReturn(run func(http.ResponseWriter, *http.Request)) *Handler_ServeHTTP_Call {
	_c.Call.Return(run)
	return _c
}

// NewHandler creates a new instance of Handler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *Handler {
	mock := &Handler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}