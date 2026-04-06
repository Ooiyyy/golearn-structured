// Package config berisi pengaturan-pengaturan dasar aplikasi, seperti koneksi database.
package config

import (
	"database/sql"

	// Import database driver MySQL. Tanda "_" (blank identifier) berarti kita hanya ingin menjalankan fungsi init() dari package database ini, tanpa memanggil fungsi spesifik lainnya secara sadar.
	_ "github.com/go-sql-driver/mysql"
)

// Connect adalah fungsi untuk membuka koneksi ke database MySQL.
// Fungsi ini mengembalikan pointer ke sql.DB yang merepresentasikan pool koneksi ke database.
func Connect() *sql.DB {
	// DSN (Data Source Name) menyimpan informasi login database.
	// Format MySQL: username:password@tcp(host:port)/nama_database
	dsn := "root:Gridy20@tcp(127.0.0.1:3306)/testdb"

	// Membuka koneksi awal ke database menggunakan driver "mysql" dan DSN di atas.
	// Perlu diingat bahwa fungsi ini belum tentu langsung mencoba terhubung ke database.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		// Jika terjadi error fatal dari syntax, aplikasi dihentikan perlahan dengan panic.
		panic(err)
	}
	
	// db.Ping() digunakan untuk memastikan koneksi benar-benar berhasil terjalin ke dalam server database.
	err = db.Ping()
	if err != nil {
		// Jika gagal terhubung (database mati atau password salah), hentikan program.
		panic(err)
	}

	// Mengembalikan objek koneksi database yang siap digunakan di seluruh program.
	return db
}
