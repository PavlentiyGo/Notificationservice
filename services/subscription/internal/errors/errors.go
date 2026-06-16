package subscription_errors

import "errors"

var ErrInvalidArgument = errors.New("invalid argument in request")
var ErrNotFound = errors.New("not found")
