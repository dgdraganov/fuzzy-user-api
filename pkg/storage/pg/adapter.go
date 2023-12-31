package pg

import (
	"fmt"

	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
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

func (db *database) Migrate(dst ...interface{}) error {
	if err := db.pg.AutoMigrate(dst...); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}

// Connect initializes a gorm.DB object and connects to the postgres db
func (db *database) Connect() error {
	dsn := db.buildDSN()
	pgdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
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

func (db *database) GetUser(email string) (model.User, error) {
	var user model.User
	res := db.pg.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return model.User{}, fmt.Errorf("db query: %w", res.Error)
	}
	return user, nil
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
