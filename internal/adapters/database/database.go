package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ivmello/kakebo-go-api/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	Connect() Connection
}

type Connection interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Close()
}

type DB struct {
	config *config.Config
}

func New(config *config.Config) Database {
	return &DB{
		config,
	}
}

func (d *DB) Config() *pgxpool.Config {
	dbConfig, err := pgxpool.ParseConfig(d.config.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = 4
	dbConfig.MinConns = 0
	dbConfig.MaxConnLifetime = time.Hour
	dbConfig.MaxConnIdleTime = time.Minute * 30
	dbConfig.HealthCheckPeriod = time.Minute
	dbConfig.ConnConfig.ConnectTimeout = time.Second * 5

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database")
	}

	return dbConfig
}

func (d *DB) Connect() Connection {
	connPool, err := pgxpool.NewWithConfig(context.Background(), d.Config())
	if err != nil {
		log.Fatal("Error while creating connection to the database")
	}

	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error while acquiring connection from the database pool")
	}
	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("Could not ping database")
	}

	return connPool
}
