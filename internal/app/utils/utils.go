package utils

import "github.com/sirupsen/logrus"

type Engineer struct {
	config *Config
	logger *logrus.Logger
}

func New(config *Config) *Engineer {
	return &Engineer{
		config: config,
		logger: logrus.New(),
	}
}

func (input *Engineer) Start() error {
	if err := input.configureLogger(); err != nil {
		return err
	}

	input.logger.Info("engineer starting work")
	return nil
}

func (input *Engineer) configureLogger() error {
	level, err := logrus.ParseLevel(input.config.LogLevel)
	if err != nil {
		return err
	}

	input.logger.SetLevel(level)
	return nil
}
