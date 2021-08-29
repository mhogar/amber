package config

import (
	"authserver/common"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/spf13/viper"
)

// Config is a struct with fields needed for configuring the application.
type Config struct {
	TokenConfig            TokenConfig            `yaml:"token"`
	DatabaseConfig         DatabaseConfig         `yaml:"database"`
	PasswordCriteriaConfig PasswordCriteriaConfig `yaml:"password_criteria"`
}

type TokenConfig struct {
	// DefaultIssuer is the the value that will go in the "issuer" field for default tokens.
	DefaultIssuer string `yaml:"default_issuer"`

	// Lifetime is the length of time a token will be valid for.
	Lifetime int64 `yaml:"lifetime"`
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

	//calc the root dir using the provided path and the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return common.ChainError("error getting working directory", err)
	}
	rootDir := path.Join(wd, dir)

	//read the config file
	data, err := ioutil.ReadFile(path.Join(rootDir, fmt.Sprint("config.", viper.Get("env"), ".yml")))
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
	viper.Set("root_dir", rootDir)
	viper.Set("token", cfg.TokenConfig)
	viper.Set("database", cfg.DatabaseConfig)
	viper.Set("password_criteria", cfg.PasswordCriteriaConfig)

	return nil
}

// GetAppRoot gets the app root directory for the application.
func GetAppRoot() string {
	return viper.GetString("root_dir")
}

func GetTokenConfig() TokenConfig {
	return viper.Get("token").(TokenConfig)
}
