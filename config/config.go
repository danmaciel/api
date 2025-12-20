package config

import (
	"fmt"
	"os"
	"strconv"
)

// configuração principal da aplicação
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// configuração do servidor
type ServerConfig struct {
	Port int
	Host string
}

// Configuração do banco de dados
type DatabaseConfig struct {
	Driver   string
	FilePath string
}

// carrega as configurações do ambiente ou usa valores padrão
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("SERVER_PORT", 8080),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "sqlite"),
			FilePath: getEnv("DB_FILE_PATH", "./database/api.db"),
		},
	}
}

// helper que ajuda a retornar valores de ambiente com default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// helper que ajuda a retornar valores inteirosde ambiente com default
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Helper que retorna um print com informações do servidor
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
