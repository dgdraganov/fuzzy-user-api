package pg

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	config *DbConfig
	pg     *gorm.DB
}

type DbConfig struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     string
	SSLMode  string
	TimeZone string
}

func NewPostgresDb(config *DbConfig) *Database {
	return &Database{
		config: config,
	}
}

func (db *Database) Connect() error {
	// ex. "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Europe/Sofia"
	dsn := db.buildDSN()
	pgdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gorm open connection: %w", err)
	}
	db.pg = pgdb
	return nil
}

func (db *Database) buildDSN() string {
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

/*


dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

*/
