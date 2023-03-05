package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New creates and returns a new Form type object
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks whether this request has particular field or not
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)

	if x != "" {
		return true
	}

	f.Errors.Add(field, "This is a required field")
	return false
}

// Required does the same as HasMany but in an efficient manner
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		isPresent := f.Get(field)
		if strings.TrimSpace(isPresent) == "" {
			f.Errors.Add(field, "This is a required field")
		}
	}
}

// HasMany checks whether this request has all field or not
func (f *Form) HasMany(fields []string, r *http.Request) bool {
	for _, field := range fields {
		fs := f.Has(field, r)
		if !fs {
			return false
		}
	}

	return true
}

// Valid returns true if there's no error in the form object
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)

	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field should have a min lenght of %d", length))
		return false
	}

	return true
}

func (f *Form) IsValidEmail(field string, r *http.Request) bool {
	x := r.Form.Get(field)

	if x != "" && !govalidator.IsEmail(x) {
		f.Errors.Add(field, "This is not a valid email")
		return false
	} else if x == "" {
		f.Errors.Add(field, "This is an empty field")
		return false
	}
	return true
}
