package config

import (
	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "app"

// DB contains database settings
type DB struct {
	Host          string `required:"true"`
	User          string `required:"true"`
	Password      string `required:"true"`
	Port          int    `required:"true" default:"5432"`
	Name          string `required:"true"`
	SSLMode       string `required:"true" envconfig:"SSL_MODE" default:"disable" desc:"disable,allow,prefer,require,verify-ca,verify-full"`
	MigrationsDir string `required:"true" envconfig:"MIGRATIONS_DIR" default:"migrations"`
}

// Log contains log settings
type Log struct {
	Level string `required:"true" default:"error" desc:"debug,info,warn,error,off"`
	Path  string `required:"true" default:"/dev/console"`
}

// JWT contains JWT settings
type JWT struct {
	Secret  string `required:"true"`
	TTL     int    `required:"true" default:"60"`
	Issuer  string `required:"true" default:"Shortener"`
	Subject string `required:"true" default:"SHORTENER"`
}

// Config contains app settings
type Config struct {
	Port          int    `required:"true"`
	Host          string `required:"true"`
	ServerTimeout int    `required:"true" default:"60" envconfig:"SERVER_TIMEOUT"`
	IsShowEnvInfo bool   `default:"true" envconfig:"IS_SHOW_ENV_INFO"`
	DB            DB
	Log           Log
	JWT           JWT
	HashSalt      string `envconfig:"HASH_SALT"`
	HashMinLength int    `envconfig:"HASH_MIN_LENGTH"`
}

// BuildConfig creates a configuration structure and
// populates the specified structure based on environment variables
func BuildConfig() (*Config, error) {
	c := &Config{}

	err := envconfig.Process(envPrefix, c)

	if err != nil {
		return nil, err
	}

	if c.IsShowEnvInfo {
		if err = envconfig.Usage(envPrefix, c); err != nil {
			return nil, err
		}
	}

	return c, nil
}
