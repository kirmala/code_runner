package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

func Load(cfgPath string, cfg any) error {
	if cfgPath == "" {
		return errors.New("config is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist by this path: %s", cfgPath)
	}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		return fmt.Errorf("error reading config: %s", err)
	}
	return nil
}
