package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	Store  StoreConf
	DB     DBConf
	HTTP   HTTPConf
	GRPC   GRPCConfig
}

type GRPCConfig struct {
	Addr string
}

type HTTPConf struct {
	Addr              string
	ReadHeaderTimeout time.Duration
}

type LoggerConf struct {
	Level string
	Type  string
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
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			fmt.Printf("Error loading .env file: %s", err)
		}
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
		fmt.Printf("error unmarshaling config: %s", err)
		os.Exit(1)
	}
	return cfg
}
