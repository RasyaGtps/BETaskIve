package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
	"taskive/utils"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServerPort string
}

var (
	AppConfig Config
	appLogger = utils.NewLogger()
)

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	AppConfig = Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		ServerPort: os.Getenv("SERVER_PORT"),
	}

	return nil
}

type customDBLogger struct {
	logger *utils.Logger
}

func (l *customDBLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *customDBLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.LogInfo("DB", fmt.Sprintf(msg, data...))
}

func (l *customDBLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.LogWarning("DB", fmt.Sprintf(msg, data...))
}

func (l *customDBLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.LogError("DB", fmt.Sprintf(msg, data...))
}

func (l *customDBLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	
	if err != nil {
		l.logger.LogError("DB", fmt.Sprintf("Error: %v SQL: %v", err, sql))
		return
	}

	l.logger.LogDB("QUERY", sql, rows, elapsed)
}

func createDatabase() error {
	connStr := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable",
		AppConfig.DBHost,
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBPort,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error connecting to postgres: %w", err)
	}
	defer db.Close()

	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s')", AppConfig.DBName)
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking database existence: %w", err)
	}

	if !exists {
		createQuery := fmt.Sprintf("CREATE DATABASE %s", AppConfig.DBName)
		_, err = db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("error creating database: %w", err)
		}
		appLogger.LogSuccess("DB", fmt.Sprintf("Database '%s' created successfully", AppConfig.DBName))
	} else {
		appLogger.LogInfo("DB", fmt.Sprintf("Database '%s' already exists", AppConfig.DBName))
	}

	return nil
}

func InitDB() (*gorm.DB, error) {
	if err := createDatabase(); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		AppConfig.DBHost,
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBName,
		AppConfig.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &customDBLogger{logger: appLogger},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	appLogger.LogSuccess("DB", "Connected to database successfully")
	return db, nil
} 