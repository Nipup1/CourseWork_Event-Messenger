package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	SroragePath string `yaml:"storage_path"`
	RedisPath   string `yaml:"redis_path"`
	Secret      string `yaml:"secret"`
	Hostname    string `yaml:"hostname"`
	GRPCAddres  string `yaml:"grpc_addres"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("GONFIG_PATH_MESSENGER")
	}

	return res
}
