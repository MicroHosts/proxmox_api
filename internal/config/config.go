package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	MySQL   MySQL   `yaml:"mysql"`
	Proxmox Proxmox `yaml:"proxmox"`
}

type Proxmox struct {
	PmUser   string `yaml:"PM_USER"`
	PmToken  string `yaml:"PM_TOKEN"`
	PmApiUrl string `yaml:"PM_API_URL"`
}

type MySQL struct {
	Host     string `yaml:"host" env:"MYSQL_HOST"`
	Port     string `yaml:"port" env:"MYSQL_PORT"`
	User     string `yaml:"user" env:"MYSQL_USER"`
	Password string `yaml:"pass" env:"MYSQL_PASSWORD"`
	DB       string `yaml:"db" env:"MYSQL_DB"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		//logger := logging.GetLogger()
		//logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			//help, _ := cleanenv.GetDescription(instance, nil)
			//logger.Info(help)
			//logger.Fatal(err)
		}
	})
	return instance
}
