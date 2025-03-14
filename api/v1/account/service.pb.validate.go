// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: v1/account/service.proto

package account

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

// Validate checks the field values on CreateRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CreateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CreateRequestMultiError, or
// nil if none found.
func (m *CreateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetAccountId()); l < 1 || l > 30 {
		err := CreateRequestValidationError{
			field:  "AccountId",
			reason: "value length must be between 1 and 30 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_CreateRequest_AccountId_Pattern.MatchString(m.GetAccountId()) {
		err := CreateRequestValidationError{
			field:  "AccountId",
			reason: "value does not match regex pattern \"^[a-zA-Z0-9]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetNickname()); l < 1 || l > 30 {
		err := CreateRequestValidationError{
			field:  "Nickname",
			reason: "value length must be between 1 and 30 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if uri, err := url.Parse(m.GetAvatar()); err != nil {
		err = CreateRequestValidationError{
			field:  "Avatar",
			reason: "value must be a valid URI",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	} else if !uri.IsAbs() {
		err := CreateRequestValidationError{
			field:  "Avatar",
			reason: "value must be absolute",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for IsAdmin

	if len(errors) > 0 {
		return CreateRequestMultiError(errors)
	}

	return nil
}

// CreateRequestMultiError is an error wrapping multiple validation errors
// returned by CreateRequest.ValidateAll() if the designated constraints
// aren't met.
type CreateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateRequestMultiError) AllErrors() []error { return m }

// CreateRequestValidationError is the validation error returned by
// CreateRequest.Validate if the designated constraints aren't met.
type CreateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateRequestValidationError) ErrorName() string { return "CreateRequestValidationError" }

// Error satisfies the builtin error interface
func (e CreateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateRequestValidationError{}

var _CreateRequest_AccountId_Pattern = regexp.MustCompile("^[a-zA-Z0-9]+$")

// Validate checks the field values on CreateResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CreateResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CreateResponseMultiError,
// or nil if none found.
func (m *CreateResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccountId

	// no validation rules for Nickname

	// no validation rules for Avatar

	// no validation rules for IsAdmin

	// no validation rules for CreatedAt

	if len(errors) > 0 {
		return CreateResponseMultiError(errors)
	}

	return nil
}

// CreateResponseMultiError is an error wrapping multiple validation errors
// returned by CreateResponse.ValidateAll() if the designated constraints
// aren't met.
type CreateResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateResponseMultiError) AllErrors() []error { return m }

// CreateResponseValidationError is the validation error returned by
// CreateResponse.Validate if the designated constraints aren't met.
type CreateResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateResponseValidationError) ErrorName() string { return "CreateResponseValidationError" }

// Error satisfies the builtin error interface
func (e CreateResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateResponseValidationError{}

// Validate checks the field values on UpdateRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UpdateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UpdateRequestMultiError, or
// nil if none found.
func (m *UpdateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccountId

	if l := utf8.RuneCountInString(m.GetNickname()); l < 1 || l > 30 {
		err := UpdateRequestValidationError{
			field:  "Nickname",
			reason: "value length must be between 1 and 30 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if uri, err := url.Parse(m.GetAvatar()); err != nil {
		err = UpdateRequestValidationError{
			field:  "Avatar",
			reason: "value must be a valid URI",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	} else if !uri.IsAbs() {
		err := UpdateRequestValidationError{
			field:  "Avatar",
			reason: "value must be absolute",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for IsAdmin

	if len(errors) > 0 {
		return UpdateRequestMultiError(errors)
	}

	return nil
}

// UpdateRequestMultiError is an error wrapping multiple validation errors
// returned by UpdateRequest.ValidateAll() if the designated constraints
// aren't met.
type UpdateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateRequestMultiError) AllErrors() []error { return m }

// UpdateRequestValidationError is the validation error returned by
// UpdateRequest.Validate if the designated constraints aren't met.
type UpdateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateRequestValidationError) ErrorName() string { return "UpdateRequestValidationError" }

// Error satisfies the builtin error interface
func (e UpdateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateRequestValidationError{}

// Validate checks the field values on UpdateResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UpdateResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UpdateResponseMultiError,
// or nil if none found.
func (m *UpdateResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccountId

	// no validation rules for Nickname

	// no validation rules for Avatar

	// no validation rules for IsAdmin

	// no validation rules for CreatedAt

	if len(errors) > 0 {
		return UpdateResponseMultiError(errors)
	}

	return nil
}

// UpdateResponseMultiError is an error wrapping multiple validation errors
// returned by UpdateResponse.ValidateAll() if the designated constraints
// aren't met.
type UpdateResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateResponseMultiError) AllErrors() []error { return m }

// UpdateResponseValidationError is the validation error returned by
// UpdateResponse.Validate if the designated constraints aren't met.
type UpdateResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateResponseValidationError) ErrorName() string { return "UpdateResponseValidationError" }

// Error satisfies the builtin error interface
func (e UpdateResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateResponseValidationError{}

// Validate checks the field values on DeleteRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *DeleteRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DeleteRequestMultiError, or
// nil if none found.
func (m *DeleteRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccountId

	if len(errors) > 0 {
		return DeleteRequestMultiError(errors)
	}

	return nil
}

// DeleteRequestMultiError is an error wrapping multiple validation errors
// returned by DeleteRequest.ValidateAll() if the designated constraints
// aren't met.
type DeleteRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteRequestMultiError) AllErrors() []error { return m }

// DeleteRequestValidationError is the validation error returned by
// DeleteRequest.Validate if the designated constraints aren't met.
type DeleteRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteRequestValidationError) ErrorName() string { return "DeleteRequestValidationError" }

// Error satisfies the builtin error interface
func (e DeleteRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteRequestValidationError{}

// Validate checks the field values on GetInfoRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *GetInfoRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetInfoRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in GetInfoRequestMultiError,
// or nil if none found.
func (m *GetInfoRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetInfoRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccountId

	if len(errors) > 0 {
		return GetInfoRequestMultiError(errors)
	}

	return nil
}

// GetInfoRequestMultiError is an error wrapping multiple validation errors
// returned by GetInfoRequest.ValidateAll() if the designated constraints
// aren't met.
type GetInfoRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetInfoRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetInfoRequestMultiError) AllErrors() []error { return m }

// GetInfoRequestValidationError is the validation error returned by
// GetInfoRequest.Validate if the designated constraints aren't met.
type GetInfoRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetInfoRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetInfoRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetInfoRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetInfoRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetInfoRequestValidationError) ErrorName() string { return "GetInfoRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetInfoRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetInfoRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetInfoRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetInfoRequestValidationError{}

// Validate checks the field values on GetInfoResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetInfoResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetInfoResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetInfoResponseMultiError, or nil if none found.
func (m *GetInfoResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetInfoResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccountId

	// no validation rules for Nickname

	// no validation rules for Avatar

	// no validation rules for IsAdmin

	// no validation rules for IsOnline

	// no validation rules for LastLoginIp

	// no validation rules for FirstLoginTime

	// no validation rules for LastLoginTime

	// no validation rules for CreatedAt

	// no validation rules for UpdatedAt

	// no validation rules for DeletedAt

	if len(errors) > 0 {
		return GetInfoResponseMultiError(errors)
	}

	return nil
}

// GetInfoResponseMultiError is an error wrapping multiple validation errors
// returned by GetInfoResponse.ValidateAll() if the designated constraints
// aren't met.
type GetInfoResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetInfoResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetInfoResponseMultiError) AllErrors() []error { return m }

// GetInfoResponseValidationError is the validation error returned by
// GetInfoResponse.Validate if the designated constraints aren't met.
type GetInfoResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetInfoResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetInfoResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetInfoResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetInfoResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetInfoResponseValidationError) ErrorName() string { return "GetInfoResponseValidationError" }

// Error satisfies the builtin error interface
func (e GetInfoResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetInfoResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetInfoResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetInfoResponseValidationError{}

// Validate checks the field values on LoginRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LoginRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoginRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LoginRequestMultiError, or
// nil if none found.
func (m *LoginRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *LoginRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccountId

	if ip := net.ParseIP(m.GetClientIp()); ip == nil {
		err := LoginRequestValidationError{
			field:  "ClientIp",
			reason: "value must be a valid IP address",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for ClientAddr

	// no validation rules for ServerAddr

	// no validation rules for DeviceId

	if _, ok := _LoginRequest_Os_InLookup[string(m.GetOs())]; !ok {
		err := LoginRequestValidationError{
			field:  "Os",
			reason: "value must be in list [[119 101 98] [97 110 100 114 111 105 100] [105 111 115]]",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for System

	if len(errors) > 0 {
		return LoginRequestMultiError(errors)
	}

	return nil
}

// LoginRequestMultiError is an error wrapping multiple validation errors
// returned by LoginRequest.ValidateAll() if the designated constraints aren't met.
type LoginRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoginRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoginRequestMultiError) AllErrors() []error { return m }

// LoginRequestValidationError is the validation error returned by
// LoginRequest.Validate if the designated constraints aren't met.
type LoginRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoginRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoginRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoginRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoginRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoginRequestValidationError) ErrorName() string { return "LoginRequestValidationError" }

// Error satisfies the builtin error interface
func (e LoginRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoginRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoginRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoginRequestValidationError{}

var _LoginRequest_Os_InLookup = map[string]struct{}{
	"\x77\x65\x62":                 {},
	"\x61\x6E\x64\x72\x6F\x69\x64": {},
	"\x69\x6F\x73":                 {},
}

// Validate checks the field values on LoginResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LoginResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoginResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LoginResponseMultiError, or
// nil if none found.
func (m *LoginResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *LoginResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccountId

	// no validation rules for Nickname

	// no validation rules for Avatar

	// no validation rules for IsAdmin

	// no validation rules for IsOnline

	// no validation rules for LastLoginIp

	// no validation rules for FirstLoginTime

	// no validation rules for LastLoginTime

	// no validation rules for OnlineId

	if len(errors) > 0 {
		return LoginResponseMultiError(errors)
	}

	return nil
}

// LoginResponseMultiError is an error wrapping multiple validation errors
// returned by LoginResponse.ValidateAll() if the designated constraints
// aren't met.
type LoginResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoginResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoginResponseMultiError) AllErrors() []error { return m }

// LoginResponseValidationError is the validation error returned by
// LoginResponse.Validate if the designated constraints aren't met.
type LoginResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoginResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoginResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoginResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoginResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoginResponseValidationError) ErrorName() string { return "LoginResponseValidationError" }

// Error satisfies the builtin error interface
func (e LoginResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoginResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoginResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoginResponseValidationError{}

// Validate checks the field values on LogoutRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LogoutRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LogoutRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LogoutRequestMultiError, or
// nil if none found.
func (m *LogoutRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *LogoutRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccountId

	// no validation rules for OnlineId

	if m.ClientIp != nil {

		if ip := net.ParseIP(m.GetClientIp()); ip == nil {
			err := LogoutRequestValidationError{
				field:  "ClientIp",
				reason: "value must be a valid IP address",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.State != nil {

		if _, ok := _LogoutRequest_State_InLookup[m.GetState()]; !ok {
			err := LogoutRequestValidationError{
				field:  "State",
				reason: "value must be in list [0 1 2]",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return LogoutRequestMultiError(errors)
	}

	return nil
}

// LogoutRequestMultiError is an error wrapping multiple validation errors
// returned by LogoutRequest.ValidateAll() if the designated constraints
// aren't met.
type LogoutRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LogoutRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LogoutRequestMultiError) AllErrors() []error { return m }

// LogoutRequestValidationError is the validation error returned by
// LogoutRequest.Validate if the designated constraints aren't met.
type LogoutRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LogoutRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LogoutRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LogoutRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LogoutRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LogoutRequestValidationError) ErrorName() string { return "LogoutRequestValidationError" }

// Error satisfies the builtin error interface
func (e LogoutRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLogoutRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LogoutRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LogoutRequestValidationError{}

var _LogoutRequest_State_InLookup = map[int32]struct{}{
	0: {},
	1: {},
	2: {},
}

// Validate checks the field values on GenerateTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GenerateTokenRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GenerateTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GenerateTokenRequestMultiError, or nil if none found.
func (m *GenerateTokenRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GenerateTokenRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccountId

	// no validation rules for Ttl

	if len(errors) > 0 {
		return GenerateTokenRequestMultiError(errors)
	}

	return nil
}

// GenerateTokenRequestMultiError is an error wrapping multiple validation
// errors returned by GenerateTokenRequest.ValidateAll() if the designated
// constraints aren't met.
type GenerateTokenRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GenerateTokenRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GenerateTokenRequestMultiError) AllErrors() []error { return m }

// GenerateTokenRequestValidationError is the validation error returned by
// GenerateTokenRequest.Validate if the designated constraints aren't met.
type GenerateTokenRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GenerateTokenRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GenerateTokenRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GenerateTokenRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GenerateTokenRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GenerateTokenRequestValidationError) ErrorName() string {
	return "GenerateTokenRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GenerateTokenRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGenerateTokenRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GenerateTokenRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GenerateTokenRequestValidationError{}

// Validate checks the field values on GenerateTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GenerateTokenResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GenerateTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GenerateTokenResponseMultiError, or nil if none found.
func (m *GenerateTokenResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GenerateTokenResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccountId

	// no validation rules for Token

	// no validation rules for TokenExpire

	if len(errors) > 0 {
		return GenerateTokenResponseMultiError(errors)
	}

	return nil
}

// GenerateTokenResponseMultiError is an error wrapping multiple validation
// errors returned by GenerateTokenResponse.ValidateAll() if the designated
// constraints aren't met.
type GenerateTokenResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GenerateTokenResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GenerateTokenResponseMultiError) AllErrors() []error { return m }

// GenerateTokenResponseValidationError is the validation error returned by
// GenerateTokenResponse.Validate if the designated constraints aren't met.
type GenerateTokenResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GenerateTokenResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GenerateTokenResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GenerateTokenResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GenerateTokenResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GenerateTokenResponseValidationError) ErrorName() string {
	return "GenerateTokenResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GenerateTokenResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGenerateTokenResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GenerateTokenResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GenerateTokenResponseValidationError{}
