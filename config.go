package main

import (
	"log"
	"os"
)

type Config struct {
	DBHost		string
	DBUser		string
	DBPassword	string
	DBName		string
	DBPort		string
	JWTSecret	string
	StripeKey	string
	Port		string
	AdminUser	string
	AdminPass	string
}

func LoadConfig() Config {
	
	cfg := Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		StripeKey:  os.Getenv("STRIPE_SECRET_KEY"),
		Port:       os.Getenv("PORT"),
		AdminUser:  os.Getenv("ADMIN_USER"),
		AdminPass:  os.Getenv("ADMIN_PASS"),
	}


	if cfg.DBHost == "" { cfg.DBHost = "localhost" }
	if cfg.DBUser == "" { cfg.DBUser = "postgres" }
	if cfg.DBPassword == "" { cfg.DBPassword = "mysecretpassword" }
	if cfg.DBName == "" { cfg.DBName = "ecommerce_db" }
	if cfg.DBPort == "" { cfg.DBPort = "5432" }
	if cfg.JWTSecret == "" { cfg.JWTSecret = "a_highly_secured_secret_for_jwt_signing_1234567890" }
	if cfg.StripeKey == "" { 
		cfg.StripeKey = "sk_test_mocked_for_development_12345" 
	}
	if cfg.Port == "" { cfg.Port = "8080" }
	if cfg.AdminUser == "" { cfg.AdminUser = "ecommerce_admin" }
	if cfg.AdminPass == "" { cfg.AdminPass = "SuperSecureAdminPass123" }

	log.Println("Configuration loaded.")
	return cfg
}