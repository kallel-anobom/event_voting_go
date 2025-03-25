package repository

import "errors"

var (
    ErrNotFoundVotes = errors.New("votes not found")
		ErrInternal = errors.New("internal error")
		ErrInvalidArgument = errors.New("invalid argument")
		ErrConflict = errors.New("conflict")
		ErrUnauthorized = errors.New("unauthorized")
		ErrForbidden = errors.New("forbidden")
		ErrUnavailable = errors.New("service unavailable")
		ErrFailedPrecondition = errors.New("failed precondition")
		ErrOutOfRange = errors.New("out of range")
		ErrUnimplemented = errors.New("unimplemented")
		ErrDataLoss = errors.New("data loss")
		ErrUnknown = errors.New("unknown error")
		ErrAlreadyExists = errors.New("already exists")
	)