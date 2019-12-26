package config

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JamesClonk/compose-broker/env"
)

var (
	config Config
	once   sync.Once
)

type Config struct {
	SkipSSL         bool
	LogLevel        string
	LogTimestamp    bool
	Username        string
	Password        string
	CatalogFilename string
	API             API
}
type API struct {
	URL               string
	Token             string
	DefaultDatacenter string
	DefaultAccountID  string
	Retries           int
	RetryInterval     time.Duration
}

func loadConfig() {
	skipSSL, _ := strconv.ParseBool(env.Get("BROKER_SKIP_SSL_VALIDATION", "false"))
	logTimestamp, _ := strconv.ParseBool(env.Get("BROKER_LOG_TIMESTAMP", "false"))
	config = Config{
		SkipSSL:         skipSSL,
		LogLevel:        env.Get("BROKER_LOG_LEVEL", "info"),
		LogTimestamp:    logTimestamp,
		Username:        env.MustGet("BROKER_AUTH_USERNAME"),
		Password:        env.MustGet("BROKER_AUTH_PASSWORD"),
		CatalogFilename: env.Get("BROKER_CATALOG_FILENAME", "catalog.yml"),
		API: API{
			URL:               strings.TrimSuffix(env.Get("COMPOSE_API_URL", "https://api.compose.io/2016-07"), "/"),
			Token:             env.MustGet("COMPOSE_API_TOKEN"),
			DefaultDatacenter: env.Get("COMPOSE_API_DEFAULT_DATACENTER", "aws:eu-central-1"),
			DefaultAccountID:  env.Get("COMPOSE_API_DEFAULT_ACCOUNT_ID", ""),
			Retries:           3,
			RetryInterval:     3 * time.Second,
		},
	}
}

func Get() *Config {
	once.Do(func() {
		loadConfig()
	})
	return &config
}
