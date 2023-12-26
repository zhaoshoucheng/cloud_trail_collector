package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"reflect"
	"strings"
)

type InputConfig struct {
	Type      string `toml:"type"`
	EndPoint  string `toml:"endpoint" json:"endpoint"`
	AccessKey string `toml:"access_key" json:"access_key"`
	SecretKey string `toml:"secret_key" json:"secret_key"`
	RegionId  string `toml:"region_id" json:"region_id"`
}
type OutPutConfig struct {
	Type      string   `toml:"type" json:"type"`
	Condition string   `toml:"condition" json:"condition"`
	Index     string   `toml:"index" json:"index"`
	EndPoints []string `toml:"end_points" json:"end_points"`
	UserName  string   `toml:"username" json:"userName"`
	PassWord  string   `toml:"password" json:"passWord"`
}

type Config struct {
	Worker  int             `toml:"worker"`
	Inputs  []*InputConfig  `toml:"input"`
	Outputs []*OutPutConfig `toml:"output"`
}

var config *Config

func NewConfig(filePath string) (conf *Config, err error) {
	conf = &Config{}
	if err = LoadTomlCfg(filePath, conf); err != nil {
		return
	}
	config = conf
	return conf, err
}

func GetConfig() *Config {
	return config
}

func LoadTomlCfg(configFile string, cfgStruct interface{}) error {
	var (
		meta toml.MetaData
		err  error
	)
	meta, err = toml.DecodeFile(configFile, cfgStruct)
	if err != nil {
		return err
	}
	if err = structFieldsDefinedInConfig(&meta, reflect.TypeOf(cfgStruct).Elem(), []string{}); err != nil {
		return err
	}
	return nil
}

func structFieldsDefinedInConfig(meta *toml.MetaData, structType reflect.Type, hierachyTomlKeys []string) error {
	hierachyTomlKeys = append(hierachyTomlKeys, "")
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if tomlKey, tomlKeyOK := field.Tag.Lookup("toml"); tomlKeyOK {
			hierachyTomlKeys[len(hierachyTomlKeys)-1] = tomlKey
			if _, ok := field.Tag.Lookup("required"); ok {
				if !meta.IsDefined(hierachyTomlKeys...) {
					return fmt.Errorf("\"%s\" not defined in config file", strings.Join(hierachyTomlKeys, "."))
				}
			}
			// make sure all sub-config-block is mapped to a struct's pointer type
			if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
				if err := structFieldsDefinedInConfig(meta, field.Type.Elem(), hierachyTomlKeys); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
