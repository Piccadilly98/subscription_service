package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	GoModName = "go.mod"
	GoSumName = "go.sum"
	NameEnv   = ".env"

	EnvNameDbName             = "DB_NAME"
	EnvNameDbSSlMode          = "DB_SSLMODE"
	EnvNameDbPort             = "DB_PORT"
	EnvNameDbHost             = "DB_HOST"
	EnvNameDbPassword         = "DB_PASSWORD"
	EnvNameDbUser             = "DB_USER"
	EnvNameNeededCache        = "CACHE"
	EnvNameCacheTTL           = "CACHE_TTL_IN_SECONDS"
	EnvNameNeededLogUserError = "LOGGING_USER_ERROR"
	EnvNameServerAddr         = "SERVER_ADDR"
	EnvNameServerPort         = "SERVER_PORT"
)

const (
	DefaultDbHost            = "localhost"
	DefaultDbPort            = "5432"
	DefaultDbSSLMode         = "disable"
	DefaultCacheTTLInSeconds = 300

	DefaultServerAddress = "localhost"
	DefaultServerPort    = "8080"
)

type Config struct {
	ConnectionStr string

	ServerAddr       string
	ServerPort       string
	NeededCache      bool
	CacheTTLInSecond int

	LoggingUserError bool
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		path, err := findEnv()
		if err != nil {
			return nil, err
		}

		err = godotenv.Load(path)
		if err != nil {
			return nil, err
		}
	}
	serverAddr := os.Getenv(EnvNameServerAddr)
	if serverAddr == "" {
		serverAddr = DefaultServerAddress
	}
	servPort, err := validationPort(EnvNameServerPort, DefaultServerPort)
	if err != nil {
		return nil, err
	}
	nameDb := os.Getenv(EnvNameDbName)
	if nameDb == "" {
		return nil, fmt.Errorf("db_name cannot be empty")
	}
	dbSsl := os.Getenv(EnvNameDbSSlMode)
	if dbSsl == "" {
		dbSsl = DefaultDbSSLMode
	}

	dbHost := os.Getenv(EnvNameDbHost)
	if dbHost == "" {
		dbHost = DefaultDbHost
	}
	dbPort, err := validationPort(EnvNameDbPort, DefaultDbPort)
	if err != nil {
		return nil, err
	}
	dbUser := os.Getenv(EnvNameDbUser)
	if dbUser == "" {
		return nil, fmt.Errorf("db_user cannot be empty")
	}

	dbPassword := os.Getenv(EnvNameDbPassword)
	if dbPassword == "" {
		return nil, fmt.Errorf("db_password cannot be empty")
	}

	cache := getBoolEnv(EnvNameNeededCache, false)
	ttl := 0

	if cache {
		if ttlEnv := os.Getenv(EnvNameCacheTTL); ttlEnv != "" {
			res, err := strconv.Atoi(ttlEnv)
			if err != nil {
				return nil, fmt.Errorf("value CACHE_TTL_IN_SECONDS not integer")
			}
			if res <= 0 {
				return nil, fmt.Errorf("invalid CACHE_TTL_IN_SECONDS, cannot be <=0")
			}
			ttl = res
		} else {
			ttl = DefaultCacheTTLInSeconds
		}
	}
	loggingUserError := getBoolEnv(EnvNameNeededLogUserError, false)
	conf := &Config{
		ConnectionStr:    fmt.Sprintf("user=%s port=%s password=%s dbname=%s host=%s sslmode=%s", dbUser, dbPort, dbPassword, nameDb, dbHost, dbSsl),
		NeededCache:      cache,
		CacheTTLInSecond: ttl,
		LoggingUserError: loggingUserError,
		ServerAddr:       serverAddr,
		ServerPort:       servPort,
	}

	return conf, nil
}

func findEnv() (string, error) {
	path := "./"
	for {
		files, err := os.ReadDir(path)
		if err != nil {
			return "", err
		}
		for _, enrty := range files {
			if enrty.Name() == GoModName || enrty.Name() == GoSumName {
				res := path + NameEnv
				_, err := os.Open(res)
				if err != nil {
					return "", fmt.Errorf("no .env file")
				}
				return res, nil
			}
		}
		path += "../"
	}
}

func getBoolEnv(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	switch strings.ToLower(val) {
	case "true", "t", "1", "yes", "y":
		return true
	case "false", "f", "0", "no", "n":
		return false
	default:
		return defaultValue
	}
}

func validationPort(key, defaultValue string) (string, error) {
	port := os.Getenv(key)
	if port == "" {
		port = defaultValue
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("invalid port '%s': %w", port, err)
	}

	if portNum <= 0 || portNum > 65535 {
		return "", fmt.Errorf("port %d out of range (1-65535)", portNum)
	}

	return port, nil
}
