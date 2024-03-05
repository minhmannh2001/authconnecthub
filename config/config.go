package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App    `yaml:"app"`
		Log    `yaml:"logger"`
		PG     `yaml:"postgres"`
		Redis  `yaml:"redis"`
		Authen `yaml:"authen"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		Host    string `env-required:"true" yaml:"host"    env:"APP_HOST"`
		Port    string `env-required:"true" yaml:"port"    env:"APP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}

	PG struct {
		Host     string `env-required:"true" yaml:"host"     env:"PG_HOST"`
		Port     string `env-required:"true" yaml:"port"     env:"PG_PORT"`
		Username string `env-required:"true" yaml:"username" env:"PG_USERNAME"`
		Password string `env-required:"true" yaml:"password" env:"PG_PASSWORD"`
		Dbname   string `env-required:"true" yaml:"dbname"   env:"PG_DBNAME"`
		Sslmode  string `env-required:"true" yaml:"sslmode"  env:"PG_SSLMODE"`
	}

	Redis struct {
		Host     string `env-required:"true" yaml:"host"     env:"REDIS_HOST"`
		Port     string `env-required:"true" yaml:"port"     env:"REDIS_PORT"`
		Password string `env-required:"true" yaml:"password" env:"REDIS_PASSWORD"`
	}

	Authen struct {
		AdminUsername     string `env-required:"true" yaml:"admin_username"       env:"ADMIN_USERNAME"`
		AdminEmail        string `env-required:"true" yaml:"admin_email"          env:"ADMIN_EMAIL"`
		AdminPassword     string `env-required:"true" yaml:"admin_password"       env:"ADMIN_PASSWORD"`
		AccessTokenTtl    int    `env-required:"true" yaml:"access_token_ttl"     env:"ACCESS_TOKEN_TTL"`
		RefreshTokenTtl   int    `env-required:"true" yaml:"refresh_token_ttl"    env:"REFRESH_TOKEN_TTL"`
		JwtPrivateKeyPath string `env-required:"true" yaml:"jwt_private_key_path" env:"JWT_PRIVATE_KEY_PATH"`
		JwtPrivateKey     *rsa.PrivateKey
		SecretKey         string `env-required:"true" yaml:"secret_key"           env:"SECRET_KEY"`
	}
)

// instance stores the singleton instance
var instance *Config

// mutex for thread-safety
var mutex sync.Mutex

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if instance == nil {
		cfg := &Config{}

		err := cleanenv.ReadConfig("./config/config.yml", cfg)
		if err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}

		err = cleanenv.ReadEnv(cfg)
		if err != nil {
			return nil, err
		}

		cfg.Authen.JwtPrivateKey, err = readPrivateKeyFromFile(cfg.Authen.JwtPrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("error while reading private key: %w", err)
		}

		instance = cfg
	}

	return instance, nil
}

func readPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	data, _ := pem.Decode(buffer)
	privateKey, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
