// config.go 
// config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds every configurable setting for the application.
// You can add or remove fields as needed (e.g. Redis settings, API keys, etc.).
type Config struct {
	// SERVER
	Port           string        // HTTP server port (e.g. ":8080" or "8080")
	ReadTimeout    time.Duration // e.g. 5 * time.Second
	WriteTimeout   time.Duration // e.g. 10 * time.Second
	ShutdownPeriod time.Duration // graceful shutdown timeout

	// POSTGRES
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret     string
	JWTExpiryMins int // how long (in minutes) a token remains valid

	// SMTP (Mailer)
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string

	// LOGGING
	LogLevel  string // e.g. "debug" / "info" / "warn" / "error"
	LogFormat string // "text" or "json"

	// Any other integrations you might need, for example:
	// RedisAddress  string
	// RedisPassword string
	// StripeSecret   string
}

// LoadConfig reads from environment variables (or a .env file if you load one manually),
// applies defaults where necessary, and returns a populated Config struct.
// If any required variable is missing, this function will return an error.
func LoadConfig() (*Config, error) {
	// 1) SERVER defaults
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	readTO, err := time.ParseDuration(os.Getenv("SERVER_READ_TIMEOUT"))
	if err != nil || readTO <= 0 {
		readTO = 5 * time.Second
	}
	writeTO, err := time.ParseDuration(os.Getenv("SERVER_WRITE_TIMEOUT"))
	if err != nil || writeTO <= 0 {
		writeTO = 10 * time.Second
	}
	shutdownPeriod, err := time.ParseDuration(os.Getenv("SERVER_SHUTDOWN_PERIOD"))
	if err != nil || shutdownPeriod <= 0 {
		shutdownPeriod = 15 * time.Second
	}

	// 2) POSTGRES (required)
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSL := os.Getenv("DB_SSLMODE")
	if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" {
		return nil, fmt.Errorf("one of required DB env vars missing: DB_HOST, DB_PORT, DB_USER, DB_NAME")
	}
	if dbSSL == "" {
		dbSSL = "disable" // default to disable SSL in dev
	}

	// 3) JWT (required)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET must be set")
	}
	jwtExpiryStr := os.Getenv("JWT_EXPIRES_IN") // in minutes
	jwtExpiry := 60                              // default 60 minutes
	if jwtExpiryStr != "" {
		if m, parseErr := strconv.Atoi(jwtExpiryStr); parseErr == nil && m > 0 {
			jwtExpiry = m
		}
	}

	// 4) SMTP / MAILER (optional, but if you plan to send mail, require these)
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("FROM_EMAIL")
	// If you don’t intend to send email yet, you can choose not to error on missing values.
	// But if sending mail is core, uncomment the following validation:
	/*
		if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" || fromEmail == "" {
			return nil, fmt.Errorf("one of required SMTP env vars missing: SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, FROM_EMAIL")
		}
	*/

	// 5) LOGGING (optional with sensible defaults)
	logLvl := os.Getenv("LOG_LEVEL")
	if logLvl == "" {
		logLvl = "info"
	}
	logFmt := os.Getenv("LOG_FORMAT")
	if logFmt == "" {
		logFmt = "text"
	}

	cfg := &Config{
		Port:           port,
		ReadTimeout:    readTO,
		WriteTimeout:   writeTO,
		ShutdownPeriod: shutdownPeriod,

		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBSSLMode:  dbSSL,

		JWTSecret:     jwtSecret,
		JWTExpiryMins: jwtExpiry,

		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUser,
		SMTPPassword: smtpPass,
		FromEmail:    fromEmail,

		LogLevel:  logLvl,
		LogFormat: logFmt,
	}

	return cfg, nil
}
