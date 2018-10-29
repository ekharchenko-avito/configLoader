package common_suite

import (
	"fmt"
	"path"

	"github.com/ekharchenko-avito/configLoader/args_loader"
	"github.com/ekharchenko-avito/configLoader/config_loader"
	"github.com/ekharchenko-avito/configLoader/env_loader"
	"github.com/ekharchenko-avito/configLoader/json_loader"
	"github.com/ekharchenko-avito/configLoader/validate_loader"
	"github.com/ekharchenko-avito/configLoader/yaml_loader"
)

type firstPassConf struct {
	ConfigPath string `env:"CONFIG_PATH" arg:"config_path,c,configuration file path"`
}

type CommonSuite struct {
	fpc       firstPassConf
	ignoreEnv bool
}

func NewCommonSuite(configPath string, ignoreEnv bool) *CommonSuite {
	return &CommonSuite{fpc: firstPassConf{ConfigPath: configPath}, ignoreEnv: ignoreEnv}
}

func (s *CommonSuite) Prepare(cl *config_loader.ConfLoader) error {
	env := env_loader.Create()
	args := args_loader.Create()

	// preload configuration
	ch := env.Load(&s.fpc)
	for err := range ch {
		println(err)
	}
	ch = args.Load(&s.fpc) // todo: fix "unknown option" bug on this stage
	for err := range ch {
		println(err)
	}

	args.EnableHelp()

	// setup
	if s.fpc.ConfigPath != "" {
		switch path.Ext(s.fpc.ConfigPath) {
		case ".json":
			cl.AddLoader(json_loader.ByPath(s.fpc.ConfigPath))
		case ".yaml", ".yml":
			cl.AddLoader(yaml_loader.ByPath(s.fpc.ConfigPath))
		default:
			return config_loader.NewError(
				fmt.Sprintf("unknown extension type for config file %s", s.fpc.ConfigPath),
				nil,
			)
		}
	}

	if !s.ignoreEnv {
		cl.AddLoader(env)
	}
	cl.AddLoader(args)
	cl.AddLoader(validate_loader.Create())
	return nil
}

// LoadConfig default shorthand for config loading
func LoadConfig(config interface{}, defaultConfigPath string, ignoreEnv bool) error {
	prepare := NewCommonSuite(defaultConfigPath, ignoreEnv)
	loader := config_loader.NewManagedLoader(config, prepare)
	return loader.Load()
}
