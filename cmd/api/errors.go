package main

import "errors"

var (
	ErrMissingDatabaseURL = errors.New("DATABASE_URL é obrigatório")
	ErrMissingJWTSecret   = errors.New("JWT_SECRET é obrigatório")
)
