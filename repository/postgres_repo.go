package repository

import (
	main "ecommerce-api"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ecommerce-api/domain"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func NewInMemoryRepository(cfg main.Config) *PostgresRepository {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	// AutoMigrate tables (creates tables if they don't exist)
	err = db.AutoMigrate(&domain.User{}, &domain.Product{}, &domain.Cart{}, &domain.CartItem{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database connection successful and migrations complete.")

	// Create admin user if it doesn't exist (Initial setup logic)
	var admin domain.User
	if err := db.Where("username = ?", cfg.AdminUser).First(&admin).Error; err == gorm.ErrRecordNotFound {
		admin = domain.User{
			Username: cfg.AdminUser,
			Password: cfg.AdminPass, // NOTE: In real app, this MUST be hashed!
			IsAdmin: true,
		}
		if err := db.Create(&admin).Error; err != nil {
			log.Printf("Failed to create admin user: %v", err)
		} else {
			log.Println("Admin user created.")
		}
	}


	return &PostgresRepository{DB: db}
}