package conf

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type AppMode string

const (
	AppModeDevelopment AppMode = "development"
	AppModeStaging     AppMode = "staging"
	AppModeProduction  AppMode = "production"
)

var (
	Opts *Options
)

type (
	Options struct {
		Common *common `mapstructure:"common"`
		Log    *log    `mapstructure:"log"`
	}

	common struct {
		Mode AppMode `mapstructure:"mode" validate:"required,oneof=development staging production"`
	}

	log struct {
		Level      string   `mapstructure:"level" validate:"required,oneof=debug info warn error panic fatal"`
		Format     string   `mapstructure:"format" validate:"required,oneof=json console"`
		Output     []string `mapstructure:"output" validate:"required,dive,min=1"`
		ErrOutput  []string `mapstructure:"err_output" validate:"required,dive,min=1"`
		MaxAge     int      `mapstructure:"max_age"`
		MaxBackups int      `mapstructure:"max_backups"`
		MaxSize    int      `mapstructure:"max_size"`
	}
)

func Init(filePath, fileType string) (string, error) {
	v := viper.New()

	v.SetConfigFile(filePath)
	v.SetConfigType(fileType)

	// common
	v.SetDefault("common.mode", AppModeDevelopment)

	// log
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "console")
	v.SetDefault("log.output", []string{"stdout"})
	v.SetDefault("log.err_output", []string{"stderr"})
	v.SetDefault("log.max_age", 7)
	v.SetDefault("log.max_backups", 3)
	v.SetDefault("log.max_size", 500)

	if err := v.ReadInConfig(); err != nil {
		return "", fmt.Errorf("read confing failed, error: %v", err)
	}

	o := new(Options)

	if err := v.Unmarshal(o); err != nil {
		return "", fmt.Errorf("unmarshal config failed, error: %v", err)
	}

	logFlag := viper.GetString("log")
	if logFlag != "" {
		logFolder := filepath.ToSlash(filepath.Clean(logFlag))
		o.Log.Output = []string{fmt.Sprintf("%s/info.log", logFolder)}
		o.Log.ErrOutput = []string{fmt.Sprintf("%s/error.log", logFolder)}
	}

	Opts = o

	return v.ConfigFileUsed(), nil
}

func (o *Options) Validate() error {
	validate := validator.New()

	return validate.StructExcept(o)
}

func (o *Options) RunMode() AppMode {
	return o.Common.Mode
}

func (o *Options) IsProdMode() bool {
	return o.RunMode() == AppModeProduction
}

func (o *Options) IsDevMode() bool {
	return o.RunMode() == AppModeDevelopment
}
