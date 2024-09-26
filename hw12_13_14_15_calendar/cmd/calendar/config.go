package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

func NewConfig() Config {
	var res Config
	err := readToml(configFile, &res)
	if err != nil {
		panic(err)
	}
	return res
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
