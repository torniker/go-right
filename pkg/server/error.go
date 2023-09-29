package server

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/torniker/go-right/pkg/logger"
)

type customError struct {
	err     error
	Message string `json:"message"`
	stack   []string
	status  int
}

type badRequest struct {
	customError
}
type internal struct {
	customError
}

type notFound struct {
	customError
}

func (e badRequest) Error() string {
	return e.Message
}

func (e internal) Error() string {
	return e.Message
}

func (e notFound) Error() string {
	return e.Message
}

func ErrBadRequest(msg string) badRequest {
	return badRequest{
		customError{
			err:     fmt.Errorf(msg),
			Message: msg,
			stack:   stack(),
			status:  http.StatusBadRequest,
		},
	}
}

func ErrInternal(err error) internal {
	return internal{
		customError{
			err:     err,
			Message: "Something went wrong",
			stack:   stack(),
			status:  http.StatusInternalServerError,
		},
	}
}

func ErrNotFound(err error) notFound {
	return notFound{
		customError{
			err:     err,
			Message: "Not Found",
			stack:   stack(),
			status:  http.StatusNotFound,
		},
	}
}

func stack() []string {
	var ret []string
	pc := make([]uintptr, 10)
	n := runtime.Callers(3, pc)
	if n == 0 {
		return ret
	}
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		ret = append(ret, fmt.Sprintf("%s:%d", frame.File, frame.Line))
		if !more {
			break
		}
	}
	return ret
}

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	switch e := err.(type) {
	case badRequest:
		logger.ErrorWithCaller(e.stack[0], e.err)
		c.Response().WriteHeader(e.status)
		c.Response().Writer.Write([]byte(e.Message))
	case *echo.HTTPError:
		switch e.Code {
		case http.StatusNotFound:
			notfoundErr := ErrNotFound(err)
			logger.ErrorWithCaller(notfoundErr.stack[0], notfoundErr.err)
			c.Response().WriteHeader(notfoundErr.status)
			c.Response().Writer.Write([]byte(notfoundErr.Message))
		}
	default:
		internalErr := ErrInternal(err)
		logger.ErrorWithCaller(internalErr.stack[0], internalErr.err)
		c.Response().WriteHeader(internalErr.status)
		c.Response().Writer.Write([]byte(internalErr.Message))
	}
}
