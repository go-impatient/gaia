package response

import (
	"net/http"
)

func NewError(code int, msg string, opts ...Option) error {
	opts = append(opts, func(e *ResponseBody) {
		e.Message = msg
	})
	return New(code, opts...).Err()
}

func Success(opts ...Option) *ResponseBody {
	return New(http.StatusOK, opts...)
}

func Error(err error) *ResponseBody {
	return Convert(err)
}

func Errorf(code int, err error) *ResponseBody {
	e := Convert(err)
	e.Code = code
	return e
}

func Convert(err error) *ResponseBody {
	s, _ := FromError(err)
	return s
}

func FromError(err error) (e *ResponseBody, ok bool) {
	if err == nil {
		return nil, false
	}
	if e, ok := err.(*ResponseBody); ok {
		return e, true
	}
	return New(http.StatusInternalServerError, Msg(err.Error())), false
}

