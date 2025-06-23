package store

import "errors"

var ErrViewAlreadyExist = errors.New("view already exists")
var ErrViewDoesNotExist = errors.New("view does not exists")
