package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ecommerce-api/domain"
)

type PostgresRepository struct {
	DB *gorm.DB
}

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	AdminUser  string
	AdminPass  string
}

func NewPostgresRepository(cfg Config) (*PostgresRepository, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// AutoMigrate tables (creates tables if they don't exist)
	err = db.AutoMigrate(&domain.User{}, &domain.Product{}, &domain.Cart{}, &domain.CartItem{})
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Database connection successful and migrations complete.")

	// Create admin user if it doesn't exist (Initial setup logic)
	var admin domain.User
	if err := db.Where("username = ?", cfg.AdminUser).First(&admin).Error; err == gorm.ErrRecordNotFound {
		// Hash the admin password
		hashedPassword := hashPassword(cfg.AdminPass)
		admin = domain.User{
			Username: cfg.AdminUser,
			Password: hashedPassword,
			IsAdmin:  true,
		}
		if err := db.Create(&admin).Error; err != nil {
			log.Printf("Failed to create admin user: %v", err)
		} else {
			log.Println("Admin user created.")
		}
	}

	return &PostgresRepository{DB: db}, nil
}

// hashPassword is a simple utility (use bcrypt in production!)
func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
