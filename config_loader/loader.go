package config_loader

import (
	"fmt"
)

type LoaderError struct {
	message string
	Err     error
}

func NewError(message string, prev error) *LoaderError {
	return &LoaderError{message: message, Err: prev}
}

func (e *LoaderError) Error() string {
	if e.Err == nil {
		return e.message
	}
	return fmt.Sprintf("%s: %+v", e.message, e.Err)
}

type MultiError struct {
	Message string
	Errors  []error
}

func (m *MultiError) Error() (s string) {
	s = m.Message + ":"
	for _, e := range m.Errors {
		s += "\n" + e.Error()
	}
	return
}

func NewMultiError(message string, sub []error) *MultiError {
	return &MultiError{Message: message, Errors: sub}
}

type LoaderI interface {
	Load(data interface{}) <-chan error
}

func WrapSingleErrLoader(f func() error) func() <-chan error {
	return func() <-chan error {
		errCh := make(chan error)
		go func() {
			defer close(errCh)
			err := f()
			if err != nil {
				errCh <- err
			}
		}()
		return errCh
	}
}

type LoaderPrepareI interface {
	Prepare(cl *ConfLoader) error
}

type ConfLoader struct {
	Data    interface{}
	loaders []LoaderI
	Prepare LoaderPrepareI
}

func NewLoader(config interface{}, loaders ...LoaderI) *ConfLoader {
	if loaders == nil {
		loaders = make([]LoaderI, 0)
	}
	return &ConfLoader{Data: config, loaders: loaders}
}

func NewManagedLoader(config interface{}, prepare LoaderPrepareI) *ConfLoader {
	return &ConfLoader{Data: config, loaders: make([]LoaderI, 0), Prepare: prepare}
}

func (l *ConfLoader) Load() error {
	// prep phase
	if l.Prepare != nil {
		err := l.Prepare.Prepare(l)
		if err != nil {
			return err
		}
	}
	errors := make([]error, 0)
	// load phase
	for _, loader := range l.loaders {
		for err := range loader.Load(l.Data) {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return NewMultiError("error during configuration loading", errors)
	}
	return nil
}

// AddLoader adds new loader to chain
func (l *ConfLoader) AddLoader(loader LoaderI) {
	l.loaders = append(l.loaders, loader)
}
