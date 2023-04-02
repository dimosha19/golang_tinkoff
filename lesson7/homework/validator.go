package homework

import (
	"reflect"

	"github.com/pkg/errors"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

var tag = "validate"

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	res := ""
	for i := range v {
		res += v[i].Err.Error()
	}
	return res
}

func Validate(v any) error {
	errorst := ValidationErrors{}

	typ := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	if typ.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for i := 0; i < typ.NumField(); i++ {
		curr := typ.Field(i)

		if curr.Tag.Get(tag) == "" {
			continue
		}

		if !curr.IsExported() {
			errorst = append(errorst, ValidationError{ErrValidateForUnexportedFields})
			return errorst
		}

		if value.Kind() == reflect.Ptr { // TODO: check for ptr
			value = reflect.Indirect(value)
		}

		var validator IValidator

		if curr.Type.Kind() == reflect.Int {
			vvv := int(value.Field(i).Int())
			tg := curr.Tag.Get(tag)
			validator = IntValidator{vvv, tg}
		}
		if curr.Type.Kind() == reflect.String {
			vvv := string(value.Field(i).String())
			tg := curr.Tag.Get(tag)
			validator = StrValidator{vvv, tg}
		}

		if res := validator.Validate(); res != nil {
			for i := range res {
				errorst = append(errorst, ValidationError{res[i]})
			}
		}
	}
	if len(errorst) != 0 {
		return errorst
	}
	return nil
}
