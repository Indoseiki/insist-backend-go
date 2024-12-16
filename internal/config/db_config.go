package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBINSIST *gorm.DB
var DBINFOR *gorm.DB

func ConnectDBINSIST() {
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:      logger.Info,
			Colorful:      true,
			SlowThreshold: time.Second,
		},
	)

	dsnINSIST := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_INSIST_HOST"),
		os.Getenv("DB_INSIST_USER"),
		os.Getenv("DB_INSIST_PASSWORD"),
		os.Getenv("DB_INSIST_NAME"),
		os.Getenv("DB_INSIST_PORT"),
	)

	DBINSIST, err = gorm.Open(postgres.Open(dsnINSIST), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("Error connecting to the INSIST database: %v", err)
	}

	fmt.Println("Connected to the INSIST database!")

	IDBINSIST, err := DBINSIST.DB()
	if err != nil {
		log.Fatalf("Failed to get DB instance from GORM: %v", err)
	}

	migrateDriverINSIST, err := migratePostgres.WithInstance(IDBINSIST, &migratePostgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		migrateDriverINSIST,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate INSIST instance: %v", err)
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatalf("INSIST migration failed: %v", err)
		}
	}

	log.Println("INSIST migration completed successfully!")
}

func ConnectDBINFOR() {
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:      logger.Info,
			Colorful:      true,
			SlowThreshold: time.Second,
		},
	)

	dsnINFOR := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		os.Getenv("DB_INFOR_USER"),
		os.Getenv("DB_INFOR_PASSWORD"),
		os.Getenv("DB_INFOR_HOST"),
		os.Getenv("DB_INFOR_PORT"),
		os.Getenv("DB_INFOR_NAME"),
	)

	DBINFOR, err = gorm.Open(sqlserver.Open(dsnINFOR), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("Error connecting to the INFOR database: %v", err)
	}

	fmt.Println("Connected to the INFOR database!")
}
