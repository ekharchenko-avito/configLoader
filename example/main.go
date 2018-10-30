package main

import (
	"fmt"
	"github.com/ekharchenko-avito/configLoader/common_suite"
	"os"
)

func main() {
	dir, _ := os.Getwd()
	fmt.Println(dir)
	type ConfigSub struct {
		InnerParam string `yaml:"inner" json:"inner" env:"INNER" arg:"inner" required:""`
	}
	type Config struct {
		ParamStr     string    `json:"str" yaml:"str" env:"STR" arg:"str,s,string param" required:""`
		ParamBool    bool      `json:"bool" yaml:"bool" env:"BOOL" arg:"bool,b,boolean param"`
		ParamInt     int       `json:"int" yaml:"int" env:"INT" arg:"int,i,integer param"`
		ParamStruct  ConfigSub `json:"struct" yaml:"struct" env:"STRUCT_" arg:"struct"`
		Param2Struct ConfigSub `json:"struct2" yaml:"struct2" env:"STRUCT2_" arg:"struct2" required:"false"`
		ParamEnvEx   struct {
			SubParam string `env:"!FORCE_SUB"`
		} `env:"ENV_EX_"`
	}

	conf := &Config{
		ParamStr:    "default value",
		ParamBool:   false,
		ParamInt:    0,
		ParamStruct: ConfigSub{"default value"},
	}

	/* Begin things that are code? Idk. */
	err := common_suite.LoadConfig(conf, "config.json", false)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", conf)
}
