package network

import (
	"fmt"
	"strings"
)

type netError struct {
	wrapped error
	error   errType
}

func (e netError) Error() string { return string(e.error) }
func (e netError) Unwrap() error { return e.wrapped }
func (e netError) Is(target error) bool {
	if targetErr, ok := target.(errType); ok {
		return strings.HasPrefix(e.Error(), targetErr.Error())
	}
	return false
}

type errType string

func (e errType) Error() string { return string(e) }

const ErrTCPWriteFail errType = "failed to write to TCP"
const ErrTCPReadFail errType = "failed to read from TCP"

func newNetworkError(topErr error, wrappedErr error) netError {
	wrappedMessage := fmt.Sprintf("%s: %s", topErr.Error(), wrappedErr.Error())
	return netError{error: errType(wrappedMessage), wrapped: wrappedErr}
}
