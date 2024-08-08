package status

import (
	"github.com/pkg/errors"
	"net/http"
)

type Code int

type Status struct {
	code    Code
	message string
	err     error
}

var statusMap = map[int]Status{}

func Register(code Code, message string) {
	errCode := Status{
		code:    code,
		message: message,
	}
	if statusMap == nil {
		statusMap = make(map[int]Status)
	}
	statusMap[int(code)] = errCode
}

func New(code Code, message string, err error) *Status {
	return &Status{
		code:    code,
		message: message,
		err:     err,
	}
}

func Error(code Code) error {
	status, ok := statusMap[int(code)]
	if ok {
		return errors.WithStack(status)
	} else {
		return Error(http.StatusInternalServerError)
	}
}

func Wrap(code Code, err error) error {
	status, ok := statusMap[int(code)]
	if ok {
		status.err = err
		return errors.WithStack(status)
	} else {
		return Error(http.StatusInternalServerError)
	}
}

func FromError(err error) *Status {
	err = errors.Cause(err)
	var status Status
	if errors.As(err, &status) {
		return &status
	}
	return New(http.StatusInternalServerError, "内部错误", err)
}

func (e Status) Error() string {
	message := e.message
	if e.err != nil {
		message = message + ": " + e.err.Error()
	}
	return message
}

func (e Status) Code() Code {
	return e.code
}

func (e Status) Message() string {
	return e.message
}

func init() {
	Register(http.StatusInternalServerError, "内部错误")
}
