// Package repository difokuskan HANYA untuk berkomunikasi dengan database.
// Di dunia Clean Architecture, lapisan ini adalah lapisan yang langsung memanipulasi penyimpanan fisik (SQL).
package repository

import "database/sql"

// UserRepository menyimpannya koneksi database agar bisa digunakan oleh fungsi-fungsi (method) di dalamnya.
type UserRepository struct {
	DB *sql.DB // Object koneksi database bawaan Go
}

// NewUserRepository adalah sebuah constructor function untuk membuat object dari UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GetByUsername mengambil data username dan password berdasarkan nama user.
func (r *UserRepository) GetByUsername(username string) (string, string, error) {
	var user, password string

	// QueryRow digunakan untuk mengambil tepat SATU baris data (karena hasil username pasti satu).
	// Scan() menyalin data dari kolom hasil database ke alamat memori variabel 'user' dan 'password'.
	err := r.DB.QueryRow("SELECT username, password FROM users WHERE username = ?", username).Scan(&user, &password)

	return user, password, err // Return nilai sekaligus dengan tipe yang dijanjikan
}

// IsUserExist mengecek apakah username tertentu sudah pernah di daftarkan.
func (r *UserRepository) IsUserExist(username string) bool {
	var existing string
	// Coba cari data tersebut
	err := r.DB.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existing)

	// Jika tidak ada error (err == nil), maka baris dapat ditemukan (berarti true)
	return err == nil
}

// Insert menambahkan pendaftaran profil dan password baru (telah di-hash) ke dalam tabel.
func (r *UserRepository) Insert(username, password string) error {
	// db.Exec biasa digunakan untuk query yang MENGUBAH data dan tidak mengembalikan baris (INSERT, UPDATE, DELETE).
	// Tanda "?" adalah parameter query (preparedStatement) yang aman dari peretasan SQL Injection.
	_, err := r.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	return err
}

// UpdatePassword memperbarui password untuk username terkait.
func (r *UserRepository) UpdatePassword(username, password string) error {
	_, err := r.DB.Exec(
		"UPDATE users SET password = ? WHERE username = ?",
		password,
		username,
	)
	return err
}

// GetProfile mengambil data username dan dibungkus di kembalian saja.
func (r *UserRepository) GetProfile(username string) (string, error) {
	var user string

	err := r.DB.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&user)

	return user, err
}
