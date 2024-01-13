package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/an-halim/golang-api-product/config"
	"github.com/an-halim/golang-api-product/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func Connect() {
	fmt.Println("Connecting to database...")
	p  := config.Load("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		fmt.Println("Error parsing port")
		os.Exit(2)
	}


	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", 
		config.Load("DB_HOST"), config.Load("DB_USER"), config.Load("DB_PASS"), config.Load("DB_NAME"), port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		fmt.Println("Error connecting to database")
		os.Exit(3)
	}

	log.Println("Connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migrations")
	db.AutoMigrate(&model.Product{})

 DB = Dbinstance{
  Db: db,
 }

}