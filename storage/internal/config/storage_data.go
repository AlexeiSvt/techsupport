package config

import (
	"sync"
	"time"
)

type Config struct {
	Storage StorageConfig `yaml:"postgres"`
}

type StorageConfig struct {
    User     string `yaml:"user" env-default:"alexei"`
    Password string `yaml:"password" env-required:"true"`
    Host     string `yaml:"host" env-default:"localhost"`
    Port     string `yaml:"port" env-default:"5432"`
    Database string `yaml:"database" env-default:"legal_tech_db"`
    SSLMode  string `yaml:"ssl_mode" env-default:"disable"`

    MaxConns        int32         `yaml:"max_conns" env-default:"20"`
    MinConns        int32         `yaml:"min_conns" env-default:"5"`
    MaxConnIdleTime time.Duration `yaml:"max_idle_time" env-default:"30m"`
    MaxConnLifetime time.Duration `yaml:"max_lifetime" env-default:"1h"`
    HealthCheck     time.Duration `yaml:"health_check" env-default:"1m"`
}

var (
	instance *Config
	once     sync.Once
)


