package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger   LoggerConf   `toml:"logger"`   // Конфигурация логгера
	Storage  StorageConf  `toml:"storage"`  // Конфигурация хранилища
	Database DatabaseConf `toml:"database"` // Конфигурация БД
	HTTP     HTTPConf     `toml:"httpsrv"`  // Конфигурация HTTP-сервера
	GRPC     GRPCConf     `toml:"grpcsrv"`  // Конфигурация gRPC-сервера
}

type LoggerConf struct {
	Level    string `toml:"level"`    // Уровень логирования
	FilePath string `toml:"filePath"` // Путь к файлу для записи логов (если пусто, пишет в Stderr)
}

type StorageConf struct {
	Type string `toml:"type"` // Тип хранилища: "sql" или "in-memory"
}

type DatabaseConf struct {
	Host   string `toml:"host"`   // Хост базы данных
	Port   string `toml:"port"`   // Порт базы данных
	User   string `toml:"user"`   // Пользователь базы данных
	Pwd    string `toml:"pwd"`    // Пароль базы данных
	DBName string `toml:"dbname"` // Имя базы данных
}

type HTTPConf struct {
	Host string `toml:"host"` // Хост HTTP-сервера
	Port string `toml:"port"` // Порт HTTP-сервера
}

type GRPCConf struct {
	Host string `toml:"host"` // Хост gRPC-сервера
	Port string `toml:"port"` // Порт gRPC-сервера
}

func NewConfig() *Config {
	return &Config{
		Logger: LoggerConf{
			Level:    "INFO",
			FilePath: "",
		},
		Storage: StorageConf{
			Type: "in-memory",
		},
		HTTP: HTTPConf{
			Host: "0.0.0.0",
			Port: "50050",
		},
		GRPC: GRPCConf{
			Host: "0.0.0.0",
			Port: "50051",
		},
	}
}

func LoadConfig(path string) (*Config, error) {
	var cfg Config

	// Чтение файла конфигурации
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Декодирование TOML
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
