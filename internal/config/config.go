package config

import (
	"fmt"
	"os"
)

type Config struct {
	// адрес http сервера, чтобы быстро менять порт/хост
	HTTPAddr string
	// блок с настройками базы
	DB DBConfig
}

type DBConfig struct {
	// хост постгры
	Host string
	// порт постгры
	Port string
	// юзер для коннекта
	User string
	// пароль для коннекта
	Password string
	// имя базы
	Name string
	// режим ssl
	SSLMode string
}

func Load() Config {
	// дефолты под docker-compose, чтобы все встало с одной кнопки
	return Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "db"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "app"),
			Password: getEnv("DB_PASSWORD", "app"),
			Name:     getEnv("DB_NAME", "numbers"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

func (c DBConfig) DSN() string {
	// собираем dsn для lib/pq
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
		c.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	// маленький хелпер: если env нет, юзаем дефолт
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
