package types

import "net/http"

type ErrorSt struct {
	Msg    string
	Status int
	Err    error
}

// func (e ErrorSt) Msg() string {
// 	return e.msg
// }

// func (e ErrorSt) Status() int {
// 	return e.status
// }

func (e ErrorSt) Error() string {
	return e.Err.Error()
}

func ErrInvalidID(e error) ErrorSt {
	return ErrorSt{
		Msg:    "invalid ID",
		Status: http.StatusBadRequest,
		Err:    e,
	}
}

func ErrInvalidParams(e error) ErrorSt {
	return ErrorSt{
		Msg:    "invalid parameters",
		Status: http.StatusBadRequest,
		Err:    e,
	}
}
func ErrUnauthorized(e error) ErrorSt {
	return ErrorSt{
		Msg:    "unauthorized",
		Status: http.StatusUnauthorized,
		Err:    e,
	}
}

func ErrNotFound(e error) ErrorSt {
	return ErrorSt{
		Msg:    "resource not found",
		Status: http.StatusNotFound,
		Err:    e,
	}
}
func ErrInternal(e error) ErrorSt {
	return ErrorSt{
		Msg:    "internal server error",
		Status: http.StatusInternalServerError,
		Err:    e,
	}
}
func ErrUnavailableDate(e error) ErrorSt {
	return ErrorSt{
		Msg:    "unavailable date",
		Status: http.StatusUnprocessableEntity,
		Err:    e,
	}
}
