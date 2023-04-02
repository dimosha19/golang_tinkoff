package homework

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
var IsDigit = regexp.MustCompile(`^[-]*[0-9]+$`).MatchString

type IValidator interface {
	Validate() []error
}

type IntValidator struct {
	val any
	tag string
}

type StrValidator struct {
	val any
	tag string
}

func (v StrValidator) Validate() []error {

	value := v.val.(string)
	errs := []error{}
	elems := strings.Split(v.tag, " ") // TODO: check for multi tag
	for i := range elems {
		pair := strings.Split(elems[i], ":")
		switch pair[0] {
		case "max":
			if err := v.ValidateMax(value, pair[1]); err != nil {
				if err.Error() == ErrInvalidValidatorSyntax.Error() {
					errs = append(errs, err)
				} else {
					s := value + " - " + err.Error()
					errs = append(errs, errors.New(s))
				}
			}
		case "min":
			if err := v.ValidateMin(value, pair[1]); err != nil {
				if err.Error() == ErrInvalidValidatorSyntax.Error() {
					errs = append(errs, err)
				} else {
					s := value + " - " + err.Error()
					errs = append(errs, errors.New(s))
				}
			}
		case "in":
			if err := v.ValidateIn(value, pair[1]); err != nil {
				if err.Error() == ErrInvalidValidatorSyntax.Error() {
					errs = append(errs, err)
				} else {
					s := value + " - " + err.Error()
					errs = append(errs, errors.New(s))
				}
			}
		case "len":
			if err := v.ValidateLen(value, pair[1]); err != nil {
				if err.Error() == ErrInvalidValidatorSyntax.Error() {
					errs = append(errs, err)
				} else {
					s := value + " - " + err.Error()
					errs = append(errs, errors.New(s))
				}
			}
		}
	}
	return errs
}

func (v StrValidator) ValidateIn(val string, tagVal string) error {
	if tagVal == "" {
		return ErrInvalidValidatorSyntax
	}
	tagArr := strings.Split(tagVal, ",") // TODO: [string] вхождение в {string, int, string, ... };
	for i := range tagArr {
		if strings.Contains(val, tagArr[i]) {
			return nil
		}
	}
	return errors.New("field does not contain required value")
}

func (v StrValidator) ValidateMax(val string, tagVal string) error {
	if !IsDigit(tagVal) {
		return ErrInvalidValidatorSyntax
	}
	elem, err := strconv.Atoi(tagVal)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}
	if len(val) <= elem {
		return nil
	}
	return errors.New("field does not fit according to the restriction from above")
}

func (v StrValidator) ValidateMin(val string, tagVal string) error {
	elem, err := strconv.Atoi(tagVal)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}
	if len(val) >= elem {
		return nil
	}
	return errors.New("the field does not fit the limit below")
}

func (v StrValidator) ValidateLen(val string, tagVal string) error {
	elem, err := strconv.Atoi(tagVal)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}
	if len(val) == elem {
		return nil
	}
	return errors.New("field has an invalid length")
}

func (v IntValidator) ValidateIn(val int, tagVal string) error {
	if tagVal == "" {
		return ErrInvalidValidatorSyntax
	}
	tagArr := strings.Split(tagVal, ",") // TODO: [int] вхождение в {int, STRING, int, ... };
	for i := range tagArr {
		elem, err := strconv.Atoi(tagArr[i])
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		if val == elem {
			return nil
		}
	}
	return errors.New("field does not contain required value")
}

func (v IntValidator) ValidateMin(val int, tagVal string) error {
	elem, err := strconv.Atoi(tagVal)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}
	if val >= elem {
		return nil
	}
	return errors.New("the field does not fit the limit below")
}

func (v IntValidator) ValidateMax(val int, tagVal string) error {
	if !IsDigit(tagVal) {
		return ErrInvalidValidatorSyntax
	}
	elem, err := strconv.Atoi(tagVal)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}
	if val <= elem {
		return nil
	}
	return errors.New("field does not fit according to the restriction from above")
}

func (v IntValidator) Validate() []error {
	value := v.val.(int)
	errs := []error{}
	elems := strings.Split(v.tag, " ") // TODO: check for multi tag
	for i := range elems {
		pair := strings.Split(elems[i], ":")
		switch pair[0] {
		case "max":
			if err := v.ValidateMax(value, pair[1]); err != nil {
				if err.Error() == ErrInvalidValidatorSyntax.Error() {
					errs = append(errs, err)
				} else {
					s := strconv.Itoa(value) + " - " + err.Error()
					errs = append(errs, errors.New(s))
				}
			}
		case "min":
			if err := v.ValidateMin(value, pair[1]); err != nil {
				if err.Error() == ErrInvalidValidatorSyntax.Error() {
					errs = append(errs, err)
				} else {
					s := strconv.Itoa(value) + " - " + err.Error()
					errs = append(errs, errors.New(s))
				}
			}
		case "in":
			if err := v.ValidateIn(value, pair[1]); err != nil {
				if err.Error() == ErrInvalidValidatorSyntax.Error() {
					errs = append(errs, err)
				} else {
					s := strconv.Itoa(value) + " - " + err.Error()
					errs = append(errs, errors.New(s))
				}
			}
		}
	}
	return errs
}
