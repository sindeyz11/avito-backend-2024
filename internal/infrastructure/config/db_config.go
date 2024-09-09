package config

import (
	"database/sql"
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
		Username: os.Getenv("PGSQL_USER"),
		Password: os.Getenv("PGSQL_PASSWORD"),
		Port:     os.Getenv("PGSQL_PORT"),
		Database: os.Getenv("PGSQL_DATABASE"),
	}
}

func NewPostgresConn(c *DatabaseConfig) *sql.DB {
	//connStr := fmt.Sprintf(
	//	"postgresql://%s:%s@database/%s?sslmode=disable",
	//	c.Username,
	//	c.Password,
	//	c.Database,
	//)

	connStr := "postgres://postgres:root@localhost:5432/tendors?sslmode=disable"

	conn, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("Can't open database connection, %v", err)
		return nil
	}
	loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	conn = sqldblogger.OpenDriver(connStr, conn.Driver(), loggerAdapter)

	if err := conn.Ping(); err != nil {
		log.Fatalf("Can't open database connection, %v", err)
		return nil
	}
	return conn
}
