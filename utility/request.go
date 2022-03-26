package utility

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	maso "github.com/myOmikron/masochism"
	"io/ioutil"
	"reflect"
	"strings"
)

var json = jsoniter.Config{
	EscapeHTML:    true,
	CaseSensitive: true,
}.Froze()

const tagName = "echotools"

//ValidateJsonForm use this method to validate a json request. Annotate your struct with `echotools:"required"` to
// mark the field as required. As
func ValidateJsonForm(c echo.Context, form interface{}) error {
	t := reflect.TypeOf(form)
	e := reflect.ValueOf(form)

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return errors.New("error while reading body")
	}

	err = json.Unmarshal(b, form)
	if err != nil {
		return errors.New("error while decoding json")
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		e = e.Elem()
	}

	var missing []string
	var notEmptyViolated []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldType := field.Type
		fieldElem := e.Field(i)
		tags := strings.Split(field.Tag.Get(tagName), ";")
		cleaned := maso.Map(func(elem string) string { return strings.TrimSpace(elem) })(tags)
		required := maso.Any(func(elem string) bool { return elem == "required" })(cleaned)
		notEmpty := maso.Any(func(elem string) bool { return elem == "not empty" })(cleaned)
		jsonName := field.Tag.Get("json")
		if s := strings.Split(jsonName, ","); len(s) > 1 {
			jsonName = s[0]
		}
		isPointer := fieldType.Kind() == reflect.Ptr

		// Required validation -> Can only be done on pointer
		if required {
			if isPointer {
				if fieldElem.IsNil() {
					missing = append(missing, jsonName)
				}
			} else {
				// As this was probably not intended, output warnings
				c.Logger().Warnf("echotools required tag set on a non-pointer field: %s", jsonName)
			}
		}

		// Not empty validation -> Can only be done on string and *string
		if notEmpty {
			// Dereference pointer of needed
			if isPointer {
				fieldElem = fieldElem.Elem()
				fieldType = fieldElem.Type()
			}

			// Check if field type is string
			if fieldType.Kind() == reflect.String {
				if fieldElem.String() == "" {
					notEmptyViolated = append(notEmptyViolated, jsonName)
				}
			} else {
				// As this was probably not intended, output warnings
				c.Logger().Warnf("echotools not empty tag set on a non-string field: %s", jsonName)
			}
		}
	}
	if len(missing) == 1 {
		return errors.New(fmt.Sprintf("parameter %s is missing but required", missing[0]))
	} else if len(missing) > 1 {
		return errors.New(fmt.Sprintf("parameter %s are missing but required", strings.Join(missing, ", ")))
	}

	if len(notEmptyViolated) > 0 {
		name := strings.Join(notEmptyViolated, ", ")
		return errors.New(fmt.Sprintf("parameter %s must not be empty", name))
	}

	return nil
}
