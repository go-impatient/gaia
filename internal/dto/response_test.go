package response

import (
	"net/http"
	"testing"

	"github.com/go-impatient/gaia/internal/dto"

	"github.com/stretchr/testify/assert"
)

func TestResponseNew(t *testing.T) {
	type Test struct {
		Id   int
		Name string
	}
	code := 200
	msg := "have some error"
	data := &Test{Id: 1, Name: "test-name"}
	err := New(code, Msg(msg), Data(data)).Err()
	assert.NotNil(t, err)
	assert.Equal(t, "go-error: code = TEST_CODE ,message = have some error ,data = &{%!s(int=1) test-name}", err.Error())
	e, ok := err.(*ResponseBody)
	assert.Equal(t, true, ok)
	assert.Equal(t, code, e.Code)
	assert.Equal(t, msg, e.Message)
	assert.Equal(t, data, e.Data)
	assert.Nil(t, e.Extra)
}

func TestError(t *testing.T) {
	code1 := 500
	msg1 := "unknown error"
	err := dto.NewError(code1, msg1)
	e, ok := dto.FromError(err)
	assert.EqualValues(t, ok, true)
	assert.EqualValues(t, code1, e.Code)
	assert.EqualValues(t, msg1, e.Message)

	err = New(500)
	e, ok = dto.FromError(err)
	assert.EqualValues(t, ok, false)
	assert.EqualValues(t, http.StatusInternalServerError, e.Code)
	assert.EqualValues(t, msg1, e.Message)
}

func TestErrorf(t *testing.T) {
	type Test struct {
		Id   int
		Name string
	}
	type ExtraTest struct {
		TotalCount int
	}
	code1 := 301
	msg1 := "unknown error"
	data1 := &Test{Id: 1, Name: "test-name"}
	extra1 := &ExtraTest{TotalCount: 500}
	err := dto.NewError(code1, msg1, Data(data1), Extra(extra1))
	e, ok := dto.FromError(err)
	assert.EqualValues(t, true, ok)
	assert.EqualValues(t, code1, e.Code)
	assert.EqualValues(t, msg1, e.Message)
	assert.EqualValues(t, data1, e.Data)
	assert.EqualValues(t, extra1, e.Extra)

	msg2 := "test error"
	err = New(500)
	e, ok = dto.FromError(err)
	assert.EqualValues(t, false, ok)
	assert.EqualValues(t, http.StatusInternalServerError, e.Code)
	assert.EqualValues(t, msg2, e.Message)
}
