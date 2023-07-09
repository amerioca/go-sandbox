package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/antoniopapa/go-admin/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	// database, err := gorm.Open(mysql.Open("root:rootroot@/go_admin"), &gorm.Config{})

	// if err != nil {
	// 	panic("Could not connect to the database")
	// }

	var err error
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=America%%2FSao_Paulo",
		"root",
		"root",
		"localhost",
		"admin_go",
	)
	// fmt.Print(dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   newLogger,
	})

	if err != nil {
		panic("Could not connect with the database!")
	}

	// DB = database

	DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Product{}, &models.Order{}, &models.OrderItem{})
}
