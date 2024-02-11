package settings

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPaht = "settings/local.yaml"
)

type Settings struct {
	Env    	 string    `yaml:"env" env-required: "true"`
	Secret 	 string    `yaml:"secret" env-required: "true"`
	Port   	 string    `yaml:"port" env-required: "true" envDefault: "8080"`
	Telegram *Telegram `yaml:"telegram"`
}

type Telegram struct {
	Token  string `yaml:"bot_token"`
	ChatID int 	  `yaml:"chat_id"` 
}

func MustLoad() *Settings {
	if _, err := os.Stat(configPaht); os.IsNotExist(err) {
		log.Fatal("config file %s doesn't exist", configPaht)
	}

	var stg Settings

	if err := cleanenv.ReadConfig(configPaht, &stg); err != nil {
		log.Fatal("cannot read config: %s", err)
	}

	return &stg
}