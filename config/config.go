package config

import (
	"authserver/common"
	"fmt"
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// Config is a struct with fields needed for configuring the application.
type Config struct {
	RootDir                string                 `yaml:"root_dir"`
	AppID                  string                 `yaml:"app_id"`
	DatabaseConfig         DatabaseConfig         `yaml:"database"`
	PasswordCriteriaConfig PasswordCriteriaConfig `yaml:"password_criteria"`
}

type DatabaseConfig struct {
	// ConnectionStrings is a string map that maps db keys to the connection string of the database.
	ConnectionStrings map[string]string `yaml:"connection_strings"`

	// Timeout is the default timeout all database requests should use.
	Timeout int `yaml:"timeout"`
}

type PasswordCriteriaConfig struct {
	// MinLength is the minimum length the password must be.
	MinLength int `yaml:"min_length"`

	// RequireLowerCase determines if at least one lower case letter must be present.
	RequireLowerCase bool `yaml:"require_lower_case"`

	// RequireUpperCase determines if at least one upper case letter must be present.
	RequireUpperCase bool `yaml:"require_upper_case"`

	// RequireDigit determines if at least one digit must be present.
	RequireDigit bool `yaml:"require_digit"`

	// RequireSymbol determines if at least one symbol must be present.
	RequireSymbol bool `yaml:"require_symbol"`
}

// InitConfig sets the default config values and binds environment variables. Should be called at the start of the application.
func InitConfig(dir string) error {
	//set defaults
	viper.SetDefault("db_key", "core")
	viper.SetDefault("env", "local")

	//bind environment variables
	viper.SetEnvPrefix("cfg")
	viper.BindEnv("env")

	//read the config file
	data, err := ioutil.ReadFile(path.Join(dir, fmt.Sprint("config.", viper.Get("env"), ".yml")))
	if err != nil {
		return common.ChainError("error loading config file", err)
	}

	//parse the yaml
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return common.ChainError("error parsing config file", err)
	}

	//set the config
	viper.Set("root_dir", cfg.RootDir)
	viper.Set("app_id", cfg.AppID)
	viper.Set("password_criteria", cfg.PasswordCriteriaConfig)
	viper.Set("database", cfg.DatabaseConfig)

	return nil
}

// GetAppId gets the app id for the application.
func GetAppId() uuid.UUID {
	return uuid.MustParse(viper.Get("app_id").(string))
}

// GetAppRoot gets the app root directory for the application.
func GetAppRoot() string {
	return viper.GetString("root_dir")
}
