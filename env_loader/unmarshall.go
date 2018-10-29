package env_loader

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func UnmarshalSimple(target reflect.Value, val string) error {
	switch target.Kind() {
	case reflect.String:
		target.SetString(val)
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		target.SetInt(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		target.SetFloat(f)
	case reflect.Bool:
		val = strings.ToUpper(val)
		if val == "1" || val == "TRUE" || val == "T" || val == "YES" || val == "Y" || val == "OK" || val == "ENABLE" || val == "ENABLED" {
			target.SetBool(true)
		} else if val == "" || val == "0" || val == "FALSE" || val == "F" || val == "NO" || val == "N" || val == "DISABLE" || val == "DISABLED" {
			target.SetBool(false)
		} else {
			return fmt.Errorf(
				"invalid value for bool target, valid are 1|t[rue]|y[es]|ok|enable[d] or <empty>|0|f[alse]|n[o]|disable[d] case-insensitive",
			)
		}
	case reflect.Ptr:
		return UnmarshalSimple(target.Elem(), val)
	default:
		return fmt.Errorf(
			"don't know how to unmarshal value into target",
		)
	}
	return nil
}
