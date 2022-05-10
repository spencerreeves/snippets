package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Debug       bool
	DbUrl       string
	ChannelName string
}

func LoadConfig() (*Config, error) {
	c := Config{Debug: true}

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return &c, errors.Wrap(err, "bad config")
	}

	viper.SetDefault("DEBUG", true)
	c.Debug = viper.GetBool("DEBUG")

	c.DbUrl = viper.GetString("DATABASE_URL")
	if c.DbUrl == "" {
		return &c, errors.New("database config not found")
	}

	c.ChannelName = viper.GetString("PUB_SUB_CHANNEL")
	if c.ChannelName == "" {
		return &c, errors.New("channel name not found")
	}

	return &c, nil
}
