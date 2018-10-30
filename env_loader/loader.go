package env_loader

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/ekharchenko-avito/configLoader/config_loader"
	"github.com/ekharchenko-avito/configLoader/hydration_utils"
)

// EnvTag struct fields tag name
const EnvTag = "env"

// UnmarshalStr unmarshaler function
type UnmarshalStr func(out reflect.Value, val string) error

// Loader type to provide config loading from env variables
type Loader struct {
	ums UnmarshalStr
}

// function to proceed a ever single node on config object
func (l *Loader) makeProceedVal() hydration_utils.ProceedNode {
	return func(target reflect.Value, hasTag bool, tag string, context string) error {
		if tag == "" {
			return nil
		}
		envVar := context + tag
		if strings.HasPrefix(tag, "!") {
			envVar = tag[1:]
		}
		val, present := os.LookupEnv(envVar)
		if !present {
			return nil
		}
		err := l.ums(target, val)
		if err != nil {
			return config_loader.NewError(fmt.Sprintf("unmarshaling value '%s' error", val), err)
		}
		return nil
	}
}

// function to decide whenever to recurse into embedded config struct
func (_ *Loader) recurse(
	field reflect.StructField,
	tag string,
	context string,
) (recurse bool, newContext string, err error) {
	return true, context + tag, nil
}

// Load method used by config_loader
func (l *Loader) Load(data interface{}) <-chan error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		err := hydration_utils.WalkStruct(reflect.ValueOf(data), "", EnvTag, l.makeProceedVal(), l.recurse)
		if err != nil {
			errCh <- err
		}
	}()
	return errCh
}

// Create create with default unmarshaler
func Create() *Loader {
	return &Loader{UnmarshalSimple}
}

// CreateCustom create with custom unmarshaler
func CreateCustom(ums UnmarshalStr) *Loader {
	return &Loader{ums}
}
