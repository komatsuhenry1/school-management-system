package config

import (
	userModel "schoolmanagement/internal/user/model"
	"fmt"
	"os"

	"github.com/joho/godotenv"	
	storage_go "github.com/supabase-community/storage-go"
	"github.com/supabase-community/supabase-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB            *gorm.DB
	Client        *supabase.Client
	StorageClient *storage_go.Client
)

func ConnectDB() error {

	// if os.Getenv("ENV") == "local" {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("erro ao carregar o .env: %w", err)
	}
	// }

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL não definido")
	}

	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			PrepareStmt: false,
		},
	)
	if err != nil {
		return fmt.Errorf("erro ao conectar com o banco de dados: %w", err)
	}
	if err := db.AutoMigrate(&userModel.User{}); err != nil {
		return fmt.Errorf("erro ao migrar banco de dados: %w", err)
	}
	DB = db

	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func GetClient() *supabase.Client {
	return Client
}

func GetStorageClient() *storage_go.Client {
	return StorageClient
}
