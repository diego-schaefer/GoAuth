package domain

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("usuário já existe no sistema")
	ErrInternal           = errors.New("erro interno do servidor")
	ErrInvalidCredentials = errors.New("email ou senha incorretos")
)
