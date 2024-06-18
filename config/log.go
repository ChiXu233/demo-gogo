package config

import (
	"github.com/wonderivan/logger"
)

func InitLog() error {
	logConfFile := "./conf/log.json"
	err := logger.SetLogger(logConfFile)
	if err != nil {
		return err
	}
	return nil
}
