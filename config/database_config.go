// Package config berisi pengaturan-pengaturan dasar aplikasi, seperti koneksi database.
package config

import (
	"database/sql"
	"fmt"

	// Import database driver MySQL. Tanda "_" (blank identifier) berarti kita hanya ingin menjalankan fungsi init() dari package database ini, tanpa memanggil fungsi spesifik lainnya secara sadar.
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseConfig struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func ConnectDB(cfg DatabaseConfig) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}
