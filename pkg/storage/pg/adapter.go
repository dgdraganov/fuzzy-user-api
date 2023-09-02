package pg

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DbConfig is a config type for the database construction
type DbConfig struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     string
	SSLMode  string
	TimeZone string
}

type database struct {
	config *DbConfig
	pg     *gorm.DB
}

// NewDatabase is a constructor function for the database type
func NewDatabase(config *DbConfig) *database {
	return &database{
		config: config,
	}
}

// Connect initializes a gorm.DB object and connects to the postgres db
func (db *database) Connect() error {
	dsn := db.buildDSN()
	pgdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gorm open connection: %w", err)
	}
	db.pg = pgdb
	return nil
}

func (db *database) Create(obj any) error {
	result := db.pg.Create(obj)
	return result.Error
}

func (db *database) buildDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		db.config.Host,
		db.config.User,
		db.config.Password,
		db.config.DbName,
		db.config.Port,
		db.config.SSLMode,
		db.config.TimeZone,
	)
}
