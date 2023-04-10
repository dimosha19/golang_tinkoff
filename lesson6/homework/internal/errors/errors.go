package errors

import "fmt"

var ErrBadRequest = fmt.Errorf("bad request")
var ErrForbidden = fmt.Errorf("forbidden")
