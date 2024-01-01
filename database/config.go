package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config merupakan struktur untuk menyimpan konfigurasi database
type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	SSLMode  string // Contoh: "disable", "require", "verify-full"
}

// NewDB membuat koneksi baru ke database berdasarkan konfigurasi
func NewDB(config Config) (*gorm.DB, error) {
	dsn := "host=" + config.Host + " user=" + config.Username + " password=" + config.Password + " dbname=" + config.Database + " port=" + config.Port + " sslmode=" + config.SSLMode

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
