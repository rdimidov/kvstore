package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	logger *zap.SugaredLogger
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetDefault("DEBUG", true)
	v.SetDefault("LOG_LEVEL", "DEBUG")

	debug := v.GetBool("DEBUG")
	logLevelStr := v.GetString("LOG_LEVEL")

	config := &Config{}

	if err := config.setLogger(debug, logLevelStr); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) setLogger(debug bool, logLevel string) error {
	cfg := zap.NewProductionConfig()
	if !debug {
		cfg = zap.NewProductionConfig()
	}

	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	cfg.Level = zap.NewAtomicLevelAt(level)

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	c.logger = logger.Sugar()
	return nil
}

func (c *Config) Logger() *zap.SugaredLogger {
	return c.logger
}
func (c *Config) Cleanup() {
	_ = c.logger.Sync()
}
