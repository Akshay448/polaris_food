package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"polaris_food/internal/database"
)

func initViperConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(CONFIGPATH)
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

}

func initDb() *gorm.DB {
	dbType := viper.GetString("database.type")
	connectionString := viper.GetString("database.dsn")
	dbConnector := database.NewDatabase(dbType)
	db := dbConnector.InitDB(connectionString)
	dbConnector.Migrate(db)
	return db
}
