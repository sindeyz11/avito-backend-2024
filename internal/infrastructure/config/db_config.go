package config

import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Driver   string
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type ParseConfig interface {
	PostgresConfig() *DatabaseConfig
}

type Config struct{}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) PostgresConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Username: os.Getenv("POSTGRES_USERNAME"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Database: os.Getenv("POSTGRES_DATABASE"),
		Host:     os.Getenv("POSTGRES_HOST"),
	}
}

func NewPostgresConn(c *DatabaseConfig) *sql.DB {
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)

	//connStr := "postgres://postgres:root@localhost:5432/tendors?sslmode=disable"

	conn, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("Can't open database connection, %v", err)
		return nil
	}
	loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	conn = sqldblogger.OpenDriver(connStr, conn.Driver(), loggerAdapter)

	if err = conn.Ping(); err != nil {
		log.Fatalf("Can't open database connection, %v", err)
		return nil
	}
	return conn
}
