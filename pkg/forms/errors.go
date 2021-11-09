package forms

// errors store the validation error messages in the forms.
// Key is field in the form, and values are errors in the field.
type errors map[string][]string

// Add adds an error message for a given field to the errors.
func (e errors) Add(field, msg string) {
	e[field] = append(e[field], msg)
}

// Get return the first error message for a given field in the errors.
// If there is no error messages for the given field, then return blank string.
func (e errors) Get(field string) string {
	msgs := e[field]
	if len(msgs) == 0 {
		return ""
	}
	return msgs[0]
}


