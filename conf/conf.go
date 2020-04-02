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
		Common   *common   `mapstructure:"common"`
		Log      *log      `mapstructure:"log"`
		Database *database `mapstructure:"database"`
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

	database struct {
		Dialect     string `mapstructure:"dialect" validate:"required"`
		Name        string `mapstructure:"name" validate:"required"`
		Network     string `mapstructure:"network" validate:"required,oneof=tcp unix"`
		Username    string `mapstructure:"username" validate:"required"`
		Password    string `mapstructure:"password" validate:"required"`
		Address     string `mapstructure:"address" validate:"required"`
		Charset     string `mapstructure:"charset" validate:"required"`
		Collation   string `mapstructure:"collation" validate:"required"`
		Loc         string `mapstructure:"loc" validate:"required"`
		ParseTime   bool   `mapstructure:"parseTime"`
		TablePrefix string `mapstructure:"tablePrefix"`
		MaxIdle     int    `mapstructure:"maxIdle"`
		MaxOpen     int    `mapstructure:"maxOpen"`
		MaxLifetime string `mapstructure:"maxLifetime"`
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

	// database
	v.SetDefault("database.dialect", "mysql")
	v.SetDefault("database.network", "tcp")
	v.SetDefault("database.address", "localhost:3306")
	v.SetDefault("database.username", "root")
	v.SetDefault("database.charset", "utf8mb64")
	v.SetDefault("database.collation", "utf8mb4_general_ci")
	v.SetDefault("database.loc", "UTC")
	v.SetDefault("database.parseTime", true)
	v.SetDefault("database.maxIdle", 5)
	v.SetDefault("database.maxOpen", 10)
	v.SetDefault("database.maxLifetime", "5m")

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
