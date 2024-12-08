package main

import (
	"errors"
)

var (
	ErrEmptyMessage       = errors.New("empty message")
	ErrIllegalByte        = errors.New("illegal byte")
	ErrTagsTooLong        = errors.New("tags too long")
	ErrInvalidTagContent  = errors.New("invalid tag content")
	ErrBodyTooLong        = errors.New("body too long")
	ErrNotConnectedServer = errors.New("not connected server")
	ErrSendToSelf         = errors.New("attempt to send message to self")
	ErrUIDBusy            = errors.New("too many busy uids")
)
