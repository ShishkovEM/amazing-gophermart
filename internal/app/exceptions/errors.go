package exceptions

import "errors"

var (
	ErrNoDatabaseDSN = errors.New("empty database dsn")
	ErrDuplicatePK   = errors.New("duplicate primary key")
	ErrNoValues      = errors.New("no values from select")
	ErrNoAuth        = errors.New("no Bearer token")
	ErrNoCookie      = errors.New("no cookie")
	ErrNotValidSign  = errors.New("sign is not valid")
)
