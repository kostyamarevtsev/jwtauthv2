package jwtauthv2

import (
	"bytes"
	"fmt"
)

const InternalMsg = "An internal error has occurred. Please contact technical support."

const (
	EPERMISSION = "permission"
	ECONFLICT   = "conflict"
	EINTERNAL   = "internal"
	EINVALID    = "invalid"
	ENOTFOUND   = "not_found"
)

type Error struct {
	Code    string
	Message string
	Op      string
	Err     error
}

func ErrorCode(err error) string {
	if err == nil {
		return ""
	}

	e, ok := err.(*Error)

	if ok && e.Code != "" {
		return e.Code
	}

	if ok && e.Err != nil {
		return ErrorCode(e.Err)
	}

	return EINTERNAL
}

func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	e, ok := err.(*Error)

	if ok && e.Message != "" {
		return e.Message
	}

	if ok && e.Err != nil {
		return ErrorMessage(e.Err)
	}

	return InternalMsg
}

func (e *Error) Error() string {
	var buf bytes.Buffer

	if e.Op != "" {
		fmt.Fprintf(&buf, "%s: ", e.Op)
	}

	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			fmt.Fprintf(&buf, "<%s> ", e.Code)
		}
		buf.WriteString(e.Message)
	}
	return buf.String()
}
