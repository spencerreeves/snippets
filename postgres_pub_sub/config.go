package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Debug       bool
	DbUrl       string
	ChannelName string
	TestType    string
}

func LoadConfig() (*Config, error) {
	c := Config{Debug: true}

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return &c, errors.Wrap(err, "bad config")
	}

	flags := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	flags.StringP("test", "t", "load", "Defines type of test to run")
	if err := flags.Parse(os.Args[1:]); err != nil {
		return &c, errors.Wrap(err, "invalid flags")
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

	var err error
	if c.TestType, err = flags.GetString("test"); err != nil {
		return nil, errors.Wrap(err, "failed to get test type")
	}

	return &c, nil
}
