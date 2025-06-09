package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string    `env:"ENV" yaml:"env" env-required:"true"`
	App AppConfig `yaml:"app"`
}

type AppConfig struct {
	Port         int           `env:"PORT" yaml:"port" env-required:"true"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" yaml:"write_timeout" env-default:"10"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" yaml:"read_timeout" env-default:"10"`
}

// MustLoad загружает конфигурацию из файла, путь к которому указан в флаге `config`
// или переменной окружения `CONFIG_PATH`. Если не указан путь,
// файл не существует или нет прав доступа к файлу, вызывает панику.
func MustLoad() *Config {
	cfgPath := fetchConfigPath()
	if cfgPath == "" {
		panic("config path is not set")
	}

	_, err := os.Stat(cfgPath)
	if err != nil && os.IsPermission(err) {
		panic("no permission to config: " + cfgPath)
	}
	if err != nil && os.IsNotExist(err) {
		panic("there is no config file: " + cfgPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		panic("failed to load config: " + err.Error())
	}

	return cfg
}

// fetchConfigPath получает путь к конфигурационному файлу из флага `config` или
// переменной окружения `CONFIG_PATH`. Если путь не указан, возвращает пустую строку.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
