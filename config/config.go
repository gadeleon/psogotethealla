package config

import (
	"log"

	"github.com/go-ini/ini"
)

type Config struct {
	Config *ini.File
}

func New(fname string) (*Config, error) {
	iniFile, err := ini.Load(fname)
	log.Print("loading file: ", fname)
	if err != nil {
		return nil, err
	}
	return &Config{
		Config: iniFile,
	}, nil
}
