package src

import (
	"toy/api/logger"

	"gopkg.in/ini.v1"
)

func GetValue(section, key string) (value string) {
	cfg, err := ini.Load("./setting.ini")
	if err != nil {
		logger.Error("invalid or missing 'setting.ini'")

		return
	}

	return cfg.Section(section).Key(key).MustString("none")
}
