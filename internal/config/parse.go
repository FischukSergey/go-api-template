package config

import (
	"errors"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

// ParseAndValidate декодирует .yaml файл в конфиг и валидирует его.
func ParseAndValidate(filename string) (Config, error) {
	if filename == "" {
		return Config{}, errors.New("filename is required")
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return Config{}, errors.New("file does not exist")
	}

	cfg := Config{}
	// 1) Декодим .yaml файл в конфиг.
	if err := cleanenv.ReadConfig(filename, &cfg); err != nil {
		return Config{}, err
	}
	// 2) Валидируем валидатором из internal/validator.
	if err := validator.New().Struct(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
