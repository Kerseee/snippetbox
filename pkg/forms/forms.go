package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX is a compiled pattern for checking email address.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Form embeds a url.Values and Errors field to hold any validation errors for the form data.
type Form struct {
	url.Values
	Errors errors
}

// New initialize a Form instance and return a pointer to it.
func New(data url.Values) *Form {
	return &Form{
		data,
		errors{},
	}
}

// Required check all the given fields in f.url.Values are present and not blank.
// If any fields fail this check, then add an error message into f.Errors.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MaxLength check if the length of character in the value of given field in the form
// not exceed given number d. If it fails then add an error message into f.Errors.
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d character)", d))
	}
}

// PermittedValues check if the value of given field in the form matches on of value in opts.
// If it fails then add an error message into f.Errors.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// MaxLength check if the length of character in the value of given field in the form
// not below given number d. If it fails then add an error message into f.Errors.
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

// MatchesPattern check if the value in the given field matches the given regular expression.
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

// Valid return true if there is no error in the Form.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
