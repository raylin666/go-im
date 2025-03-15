// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: v1/message/service.proto

package message

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on SendC2CMessageRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *SendC2CMessageRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SendC2CMessageRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SendC2CMessageRequestMultiError, or nil if none found.
func (m *SendC2CMessageRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *SendC2CMessageRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Seq

	// no validation rules for ToAccount

	// no validation rules for Message

	if len(errors) > 0 {
		return SendC2CMessageRequestMultiError(errors)
	}

	return nil
}

// SendC2CMessageRequestMultiError is an error wrapping multiple validation
// errors returned by SendC2CMessageRequest.ValidateAll() if the designated
// constraints aren't met.
type SendC2CMessageRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SendC2CMessageRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SendC2CMessageRequestMultiError) AllErrors() []error { return m }

// SendC2CMessageRequestValidationError is the validation error returned by
// SendC2CMessageRequest.Validate if the designated constraints aren't met.
type SendC2CMessageRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SendC2CMessageRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SendC2CMessageRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SendC2CMessageRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SendC2CMessageRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SendC2CMessageRequestValidationError) ErrorName() string {
	return "SendC2CMessageRequestValidationError"
}

// Error satisfies the builtin error interface
func (e SendC2CMessageRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSendC2CMessageRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SendC2CMessageRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SendC2CMessageRequestValidationError{}

// Validate checks the field values on SendC2CMessageResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *SendC2CMessageResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SendC2CMessageResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SendC2CMessageResponseMultiError, or nil if none found.
func (m *SendC2CMessageResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *SendC2CMessageResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return SendC2CMessageResponseMultiError(errors)
	}

	return nil
}

// SendC2CMessageResponseMultiError is an error wrapping multiple validation
// errors returned by SendC2CMessageResponse.ValidateAll() if the designated
// constraints aren't met.
type SendC2CMessageResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SendC2CMessageResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SendC2CMessageResponseMultiError) AllErrors() []error { return m }

// SendC2CMessageResponseValidationError is the validation error returned by
// SendC2CMessageResponse.Validate if the designated constraints aren't met.
type SendC2CMessageResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SendC2CMessageResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SendC2CMessageResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SendC2CMessageResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SendC2CMessageResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SendC2CMessageResponseValidationError) ErrorName() string {
	return "SendC2CMessageResponseValidationError"
}

// Error satisfies the builtin error interface
func (e SendC2CMessageResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSendC2CMessageResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SendC2CMessageResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SendC2CMessageResponseValidationError{}
