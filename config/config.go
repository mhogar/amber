package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mhogar/amber/common"

	"gopkg.in/yaml.v3"

	"github.com/spf13/viper"
)

type Config struct {
	// AppName is a client facing name to refer to the app as.
	AppName string `yaml:"app_name"`

	// DataAdapter is the name of the data adapter the app will use.
	DataAdapter string `yaml:"data_adapter,omitempty"`

	TokenConfig            TokenConfig            `yaml:"token"`
	PermissionConfig       PermissionConfig       `yaml:"permissions"`
	DatabaseConfig         DatabaseConfig         `yaml:"database"`
	FirestoreConfig        FirestoreConfig        `yaml:"firestore"`
	PasswordCriteriaConfig PasswordCriteriaConfig `yaml:"password_criteria"`
}

type TokenConfig struct {
	// DefaultIssuer is the the value that will go in the "issuer" field for default tokens.
	DefaultIssuer string `yaml:"default_issuer"`

	// Lifetime is the length of time a token will be valid for.
	Lifetime int64 `yaml:"lifetime"`
}

type PermissionConfig struct {
	// MinClientRank is the minimum rank a user must have to manage clients.
	MinClientRank int `yaml:"min_client_rank"`
}

type DatabaseConfig struct {
	// Driver is the database driver to use.
	Driver string `yaml:"driver"`

	// ConnectionStrings is a string map that maps db keys to the connection string of the database.
	ConnectionStrings map[string]string `yaml:"connection_strings"`

	// Timeout is the default timeout all database requests should use.
	Timeout int `yaml:"timeout"`
}

type FirestoreConfig struct {
	// ServiceFile is the file location for the firebase service account json.
	ServiceFile string `yaml:"service_file,omitempty"`

	// Timeout is the default timeout all firestore requests should use.
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

// InitConfig sets the default config values and binds environment variables.
// Should be called at the start of the application.
func InitConfig(dir string) error {
	//set defaults
	viper.SetDefault("db_key", "core")
	viper.SetDefault("env", "local")

	//bind environment variables
	viper.SetEnvPrefix("cfg")
	viper.BindEnv("env")
	viper.BindEnv("data_adapter")

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
	viper.Set("app_name", cfg.AppName)
	viper.SetDefault("data_adapter", cfg.DataAdapter)
	viper.Set("token", cfg.TokenConfig)
	viper.Set("permission", cfg.PermissionConfig)
	viper.Set("database", cfg.DatabaseConfig)
	viper.Set("firestore", cfg.FirestoreConfig)
	viper.Set("password_criteria", cfg.PasswordCriteriaConfig)

	return nil
}

// GetAppRoot gets the app root directory for the application.
func GetAppRoot() string {
	return viper.GetString("root_dir")
}

// GetAppName gets the name for the app.
func GetAppName() string {
	return viper.GetString("app_name")
}

// GetDataAdapter gets the data adapter currently selected for the app.
func GetDataAdapter() string {
	return viper.GetString("data_adapter")
}

// GetTokenConfig gets the token config object.
func GetTokenConfig() TokenConfig {
	return viper.Get("token").(TokenConfig)
}

// GetPermissionConfig gets the permissions config object.
func GetPermissionConfig() PermissionConfig {
	return viper.Get("permission").(PermissionConfig)
}

// GetDatabaseConfig gets the database config object.
func GetDatabaseConfig() DatabaseConfig {
	return viper.Get("database").(DatabaseConfig)
}

// GetFirestoreConfig gets the firestore config object.
func GetFirestoreConfig() FirestoreConfig {
	return viper.Get("firestore").(FirestoreConfig)
}

// GetPasswordCriteriaConfig gets the password criteria config object.
func GetPasswordCriteriaConfig() PasswordCriteriaConfig {
	return viper.Get("password_criteria").(PasswordCriteriaConfig)
}
