package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	MinIO    MinIOConfig
	Server   ServerConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

type ServerConfig struct {
	Port string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("postgres.sslmode", "disable")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("无法读取配置文件: %v", err)
	}

	// 环境变量覆盖
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
