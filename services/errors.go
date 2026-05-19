package services

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrAlreadyFollowing    = errors.New("already following")
	ErrNotFollowing        = errors.New("not following")
	ErrDuplicateUsername   = errors.New("username already exists")
	ErrCannotFollowSelf    = errors.New("cannot follow yourself")
	ErrAlreadyLiked        = errors.New("already liked")
	ErrNotLiked            = errors.New("not liked")
)
