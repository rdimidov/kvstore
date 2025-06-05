package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultServerAddr = "0.0.0.0:8080"
	defaultLogLevel   = "info"
)

type Config struct {
	Network struct {
		Address        string        `mapstructure:"address"`
		MaxMessageSize int           `mapstructure:"max_message_size"`
		ReadTimeout    time.Duration `mapstructure:"read_timeout"`
		WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	} `mapstructure:"network"`
	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`
	WAL struct {
		Enabled      bool          `mapstructure:"enabled"`
		BatchSize    int           `mapstructure:"batchSize"`
		FlushTimeout time.Duration `mapstructure:"flushTimeout"`
		Dir          string        `mapstructure:"directory"`
		MSS          int           `mapstructure:"maxSegmentSizeMB"`
	} `mapstructure:"wal"`

	logger *zap.SugaredLogger
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetDefault("network.address", defaultServerAddr)
	v.SetDefault("logging.level", defaultLogLevel)

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	// Try to read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}
	if err := cfg.setLogger(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) setLogger() error {
	cfg := zap.NewProductionConfig()

	level, err := zapcore.ParseLevel(c.Logging.Level)
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

func (c *Config) Cleanup() {
	_ = c.logger.Sync()
}

func (c *Config) Logger() *zap.SugaredLogger { return c.logger }

func (c *Config) WALEnabled() bool                    { return c.WAL.Enabled }
func (c *Config) WALBatchSize() int                   { return c.WAL.BatchSize }
func (c *Config) WALBatchFlushTimeout() time.Duration { return c.WAL.FlushTimeout }
func (c *Config) WALDirName() string                  { return c.WAL.Dir }
func (c *Config) WALMaxSegmentSize() int              { return c.WAL.MSS }
