package database

import (
	"gorm.io/gorm"
)

// DBConnector defines the interface for database connections and migrations
type DBConnector interface {
	InitDB(connection string) *gorm.DB
	Migrate(db *gorm.DB)
}

// NewDatabase initializes the database based on the type specified in the config
func NewDatabase(dbType string) DBConnector {
	switch dbType {
	case "sqlite":
		return &SQLiteConnector{}
	//case "postgres":
	//	return &PostgresConnector{}
	// Add more cases for other databases
	default:
		panic("Unsupported database type")
	}
}
