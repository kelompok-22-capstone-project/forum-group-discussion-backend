package repository

import "errors"

var (
	ErrRecordNotFound      = errors.New("repository: record with given params not found")
	ErrDatabase            = errors.New("repository: something wrong with the database")
	ErrRecordAlreadyExists = errors.New("repository: record already exists")
)
