package config

import "github.com/go-ini/ini"

type Config struct {
	Config *ini.File
}

func (c *Config) New(fname string) (*Config, error) {
	iniFile, err := ini.Load(fname)
	if err != nil {
		return nil, err
	}
	return &Config{
		Config: iniFile,
	}, nil
}
