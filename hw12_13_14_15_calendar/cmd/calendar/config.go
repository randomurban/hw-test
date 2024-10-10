package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	Store  StoreConf
	DB     DBConf
}

type LoggerConf struct {
	Level string
}

type StoreConf struct {
	StoreType string
}

const (
	StoreTypeSQL    = "SQL"
	StoreTypeMemory = "MEMORY"
)

type DBConf struct {
	DSN string
}

func NewConfig(configFile string) Config {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("Error loading .env file: %s", err)
	}
	viper.SetEnvPrefix("CALENDAR")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigFile(configFile)
	viper.SetConfigType("toml")

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error reading config file: %s", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("error unmarshaling config: ", err)
		os.Exit(1)
	}
	return cfg
}

func readToml(file string, res *Config) error {
	fileToml, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fileToml.Close()

	data, err := io.ReadAll(fileToml)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(data, res)
	if err != nil {
		return err
	}
	fmt.Printf("config: %v\n", res)
	return nil
}
