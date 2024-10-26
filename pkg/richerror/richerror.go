package richerror

import "errors"

type Kind int

type RichError struct {
	operation string
	wrapError error
	message   string
	kind      Kind
	meta      map[string]interface{}
}

func (r RichError) Error() string {
	if r.message == "" && r.wrapError != nil {
		return r.wrapError.Error()
	}

	return r.message
}

func New(op string) RichError {
	return RichError{operation: op}
}

func (r RichError) WithMessage(message string) RichError {
	r.message = message
	return r
}

func (r RichError) WithWrapError(err error) RichError {
	r.wrapError = err
	return r
}

func (r RichError) WithKind(kind Kind) RichError {
	r.kind = kind
	return r
}

func (r RichError) WithMeta(meta map[string]interface{}) RichError {
	r.meta = meta
	return r
}

func (r RichError) Kind() Kind {
	if r.kind != 0 {
		return r.kind
	}

	var err RichError
	if ok := errors.As(r.wrapError, &err); ok {
		return err.Kind()
	}

	return err.Kind()
}

func (r RichError) Message() string {
	if r.message != "" {
		return r.message
	}

	var err RichError
	if ok := errors.As(r, &err); ok {
		return err.Message()
	}

	return err.Message()
}
