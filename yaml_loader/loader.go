package yaml_loader

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	. "github.com/ekharchenko-avito/configLoader/config_loader"
)

type Loader struct {
	path string
}

func (l *Loader) Load(data interface{}) <-chan error {
	return WrapSingleErrLoader(func() error {
		if l.path == "" {
			return NewError("no path is set for yaml loader", nil)
		}
		file, err := ioutil.ReadFile(l.path)
		if err != nil {
			return NewError("failure on config file read", err)
		}
		err = yaml.Unmarshal(file, data)
		if err != nil {
			return NewError("failure on config unmarshal", err)
		}
		return nil
	})()
}

func ByPath(path string) *Loader {
	return &Loader{path}
}
