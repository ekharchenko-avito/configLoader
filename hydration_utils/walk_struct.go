package hydration_utils

import (
	"fmt"
	"reflect"

	"github.com/ekharchenko-avito/configLoader/config_loader"
)

// ProceedNode interface for function to process each node
type ProceedNode func(value reflect.Value, hasTag bool, tag string, context string) error

// RecurseChecker interface for function to decide on processing embedded struct
type RecurseChecker func(field reflect.StructField, tag string, context string) (recurse bool, newContext string, err error)

// WalkStruct function to walk recursively over config struct and process it using struct field tags
func WalkStruct(
	out reflect.Value,
	context string,
	tagName string,
	proceedNode ProceedNode,
	recurseChecker RecurseChecker,
) error {
	if out.Kind() == reflect.Ptr {
		if !out.IsNil() {
			out = out.Elem()
		} else {
			// don't bother to make new items from pointer. it is a config, just use plain values and defaults
			return nil
		}
	}
	if out.Kind() != reflect.Struct {
		return config_loader.NewError("config should be a struct for loader to work properly", nil)
	}
	cnt := out.NumField()
	for i := 0; i < cnt; i++ {
		t := out.Type().Field(i)
		tag, hasTag := t.Tag.Lookup(tagName)
		// if processing struct - attempt to look inside
		if t.Type.Kind() == reflect.Struct || t.Type.Kind() == reflect.Ptr && t.Type.Elem().Kind() == reflect.Struct {
			recurse, newContext, err := recurseChecker(t, tag, context)
			if err != nil {
				return config_loader.NewError(
					fmt.Sprintf(
						"error during processing from env  at %s, field %s (%s)",
						context,
						t.Name,
						t.Type.Name(),
					),
					err,
				)
			}
			if recurse {
				WalkStruct(out.Field(i), newContext, tagName, proceedNode, recurseChecker)
			}
			continue
		}
		err := proceedNode(out.Field(i), hasTag, tag, context)
		if err != nil {
			return config_loader.NewError(
				fmt.Sprintf(
					"error during processing from env at '%s', field %s (%s)",
					context,
					t.Name,
					t.Type.Name(),
				),
				err,
			)
		}
	}
	return nil
}
