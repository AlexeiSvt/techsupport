package config

import (
	"log"
	"techsupport/storage/constants"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)


func GetConfig() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("Note: .env file not found, using system env")
		}

		instance = &Config{}

		if err := cleanenv.ReadConfig(constants.ConfigPath, instance); err != nil {
			log.Fatalf("Fatal: check config file at %s: %v", constants.ConfigPath, err)
		}
	})
	return instance
}
