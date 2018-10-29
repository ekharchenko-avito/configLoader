package args_loader

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"os"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/ekharchenko-avito/configLoader/hydration_utils"
)

// EnvTag struct fields tag name
const EnvArgs = "arg"

// Loader type to provide config loading from env variables
type Loader struct {
	handleHelp bool
}

const formatError = "invalid args loader tag format: '%s', should be 'full_param[,<short char>][,<description>]'"

func (l *Loader) extractTagParts(tag string) (string, rune, string, error) {
	p := strings.Split(tag, ",")
	if len(p) > 3 {
		return "", 0, "", fmt.Errorf(formatError, tag)
	}
	long := p[0]
	short := ""
	description := ""
	if len(p) > 1 {
		short = p[1]
		if len(p) > 2 {
			description = p[2]
		}
	}
	var r rune = 0
	if len(short) > 0 {
		r, _ = utf8.DecodeRuneInString(short)
	}
	return long, r, description, nil
}

func (l *Loader) registerParam(target reflect.Value, long string, r rune, description string) {
	if target.Kind() == reflect.Ptr && !target.IsNil() {
		l.registerParam(target.Elem(), long, r, description)
		return
	}
	getopt.FlagLong(target.Addr().Interface(), long, r, description)
}

// function to proceed a ever single node on config object
func (l *Loader) makeProceedVal() hydration_utils.ProceedNode {
	return func(target reflect.Value, hasTag bool, tag string, context string) error {
		if tag == "" {
			return nil
		}
		long, r, description, err := l.extractTagParts(tag)
		if context != "" { // disable short name nesting
			r = 0
		}
		if err != nil {
			return err
		}
		l.registerParam(target, context+long, r, description)
		return nil
	}
}

// function to decide whenever to recurse into embedded config struct
func (l *Loader) recurse(field reflect.StructField, tag string, context string) (recurse bool, newContext string, err error) {
	long, _, _, err := l.extractTagParts(tag)
	newContext = long + "."
	if context != "" {
		newContext = context + newContext
	}
	return true, newContext, err
}

// Load method used by config_loader
func (l *Loader) Load(data interface{}) <-chan error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		err := hydration_utils.WalkStruct(reflect.ValueOf(data), "", EnvArgs, l.makeProceedVal(), l.recurse)
		if err != nil {
			errCh <- err
			return
		}

		err = getopt.Getopt(nil)
		if err != nil {
			errCh <- err
			return
		}
		if l.handleHelp && getopt.GetValue("help") == "true" {
			getopt.Usage()
			os.Exit(100)
		}
	}()
	return errCh
}

func (l *Loader) EnableHelp() {
	getopt.BoolLong("help", 'h', "show help")
	l.handleHelp = true
}

// Create create with default unmarshaler
func Create() *Loader {
	return &Loader{}
}
