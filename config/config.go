package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	Server   ServerConfig   `mapstructure:"server"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

var Global *Config

func InitConfig(cfgFile string) error {
	setDefaultConfig()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("DEX")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		log.Warn("No config file found, using defaults")
	}

	decoderConfig := &mapstructure.DecoderConfig{
		Result:           &Global,
		TagName:          "mapstructure",
		WeaklyTypedInput: true,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(viper.AllSettings()); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}

	if err := Global.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

// setDefaultConfig
func setDefaultConfig() {
	defaultValues := structToMap(DefaultConfig)
	for key, value := range defaultValues {
		viper.SetDefault(key, value)
	}
}

func structToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	fillMap(reflect.ValueOf(obj), "", result)
	return result
}

func fillMap(v reflect.Value, prefix string, result map[string]interface{}) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}
		key := tag
		if prefix != "" {
			key = prefix + "." + tag
		}

		if value.Kind() == reflect.Struct {
			fillMap(value, key, result)
		} else {
			result[key] = value.Interface()
		}
	}
}

func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Port <= 0 {
		return fmt.Errorf("invalid database port")
	}
	if c.Server.Port <= 0 {
		return fmt.Errorf("invalid server port")
	}
	return nil
}

func (c *Config) GetRootDir() string {
	return "./"
}
