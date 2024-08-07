package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type Config struct {
	Env       string     `yaml:"env" env-default:"local"`
	Server    HttpServer `yaml:"server"`
	Postrgres Postgres   `yaml:"postgres"`
}

type HttpServer struct {
	AppVersion   string        `yaml:"appVersion"`
	Port         string        `yaml:"port" env-required:"true"`
	PProfPort    string        `yaml:"pProfPort" env-required:"true"`
	JwtSecretKey string        `yaml:"jwtSecretKey" env-required:"true"`
	Mode         string        `yaml:"mode" env-default:"Development"`
	Timeout      time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout  time.Duration `yaml:"ideTimeout" env-default:"60s"`
	Debug        bool          `yaml:"debug" env-default:"false"`
}

type Postgres struct {
	Host     string `yaml:"postgresqlHost" env-required:"true"`
	Port     string `yaml:"postgresqlPort" env-required:"true"`
	User     string `yaml:"postgresqlUser" env-required:"true"`
	Password string `yaml:"postgresqlPassword" env-required:"true"`
	Dbname   string `yaml:"postgresqlDbname" env-required:"true"`
	SSLMode  bool   `yaml:"postgresqlSSLMode" env-default:"false"`
	Driver   string `yaml:"pgDriver" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
