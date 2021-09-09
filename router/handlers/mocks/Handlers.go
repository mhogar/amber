// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import (
	data "authserver/data"

	http "net/http"

	httprouter "github.com/julienschmidt/httprouter"

	mock "github.com/stretchr/testify/mock"

	models "authserver/models"
)

// Handlers is an autogenerated mock type for the Handlers type
type Handlers struct {
	mock.Mock
}

// DeleteClient provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Handlers) DeleteClient(_a0 *http.Request, _a1 httprouter.Params, _a2 *models.Session, _a3 data.DataCRUD) (int, interface{}) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) interface{}); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}

// DeleteSession provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Handlers) DeleteSession(_a0 *http.Request, _a1 httprouter.Params, _a2 *models.Session, _a3 data.DataCRUD) (int, interface{}) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) interface{}); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Handlers) DeleteUser(_a0 *http.Request, _a1 httprouter.Params, _a2 *models.Session, _a3 data.DataCRUD) (int, interface{}) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) interface{}); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}

// DeleteUserRole provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Handlers) DeleteUserRole(_a0 *http.Request, _a1 httprouter.Params, _a2 *models.Session, _a3 data.DataCRUD) (int, interface{}) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) interface{}); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}

// PatchUserPassword provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Handlers) PatchUserPassword(_a0 *http.Request, _a1 httprouter.Params, _a2 *models.Session, _a3 data.DataCRUD) (int, interface{}) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) interface{}); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}

// PostClient provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Handlers) PostClient(_a0 *http.Request, _a1 httprouter.Params, _a2 *models.Session, _a3 data.DataCRUD) (int, interface{}) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) interface{}); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}

// PostSession provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Handlers) PostSession(_a0 *http.Request, _a1 httprouter.Params, _a2 *models.Session, _a3 data.DataCRUD) (int, interface{}) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) interface{}); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}

// PostToken provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Handlers) PostToken(_a0 *http.Request, _a1 httprouter.Params, _a2 *models.Session, _a3 data.DataCRUD) (int, interface{}) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) interface{}); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}

// PostUser provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Handlers) PostUser(_a0 *http.Request, _a1 httprouter.Params, _a2 *models.Session, _a3 data.DataCRUD) (int, interface{}) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) interface{}); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}

// PostUserRole provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Handlers) PostUserRole(_a0 *http.Request, _a1 httprouter.Params, _a2 *models.Session, _a3 data.DataCRUD) (int, interface{}) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) interface{}); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}

// PutClient provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Handlers) PutClient(_a0 *http.Request, _a1 httprouter.Params, _a2 *models.Session, _a3 data.DataCRUD) (int, interface{}) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) interface{}); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}

// PutUser provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Handlers) PutUser(_a0 *http.Request, _a1 httprouter.Params, _a2 *models.Session, _a3 data.DataCRUD) (int, interface{}) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) interface{}); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}

// PutUserRole provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Handlers) PutUserRole(_a0 *http.Request, _a1 httprouter.Params, _a2 *models.Session, _a3 data.DataCRUD) (int, interface{}) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) interface{}); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}
