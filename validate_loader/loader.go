package validate_loader

import (
	"fmt"
	"reflect"

	"github.com/ekharchenko-avito/configLoader/config_loader"
)

// RequiredTag struct fields tag name
const RequiredTag = "required"

// Loader type to provide config verification
type Loader struct {
}

func (l *Loader) walkStruct(
	out reflect.Value,
	path string,
	errCh chan<- error,
) {
	if out.Kind() == reflect.Ptr {
		if out.IsNil() {
			return
		} else {
			out = out.Elem()
		}
	}
	if out.Kind() != reflect.Struct {
		errCh <- config_loader.NewError("config should be a struct for validator to work properly", nil)
		return
	}
	cnt := out.NumField()
	for i := 0; i < cnt; i++ {
		t := out.Type().Field(i)
		v := out.Field(i)
		tag, hasTag := t.Tag.Lookup(RequiredTag)
		// if processing struct - attempt to look inside
		if t.Type.Kind() == reflect.Struct || t.Type.Kind() == reflect.Ptr && t.Type.Elem().Kind() == reflect.Struct {
			if t.Type.Kind() == reflect.Ptr && v.IsNil() {
				if hasTag {
					errCh <- config_loader.NewError(
						fmt.Sprintf("missed required field '%s'", path+"."+t.Name),
						nil,
					)
					continue
				} else {
					continue
				}
			}
			if tag != "false" {
				l.walkStruct(out.Field(i), path+"."+t.Name, errCh)
			}
			continue
		}
		if !hasTag {
			continue
		}
		if v.Interface() == reflect.Zero(v.Type()).Interface() {
			fullName := t.Name
			if path != "" {
				fullName = path + "." + fullName
			}
			errCh <- config_loader.NewError(
				fmt.Sprintf("missing required field '%s', tags: `%s`", fullName, t.Tag),
				nil,
			)
			continue
		}
	}
	return
}

// Load method used by config_loader
func (l *Loader) Load(data interface{}) <-chan error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		l.walkStruct(reflect.ValueOf(data), "", errCh)
	}()
	return errCh
}

// Create create with default unmarshaler
func Create() *Loader {
	return &Loader{}
}
