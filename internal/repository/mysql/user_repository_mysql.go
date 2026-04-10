// Package mysql fokus pada detail SQL dan persistence.
package mysql

import (
	"database/sql"
	"fmt"
	"golearn-structured/internal/model"
	"golearn-structured/internal/repository"
)

// userRepositoryImpl menyimpan koneksi database agar bisa digunakan oleh fungsi-fungsi (method) di dalamnya.
type userRepositoryImpl struct {
	DB *sql.DB // Object koneksi database bawaan Go
}

// NewUserRepository adalah sebuah constructor function untuk membuat object dari UserRepository.
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepositoryImpl{DB: db}
}

// GetByUsername mengambil data username dan password berdasarkan nama user.
func (r *userRepositoryImpl) GetByUsername(username string) (*model.Users, error) {
	var user model.Users

	// QueryRow digunakan untuk mengambil tepat SATU baris data (karena hasil username pasti satu).
	// Scan() menyalin data dari kolom hasil database ke alamat memori variabel 'user' dan 'password'.
	// QueryRow cocok untuk single-record query; Scan butuh pointer sebagai target.
	err := r.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password)

	// Inilah cara membuat Pointer menjadi 'nil'!
	// Jika ada error (termasuk jika data tidak ditemukan), kita cegah kode lanjut ke bawah.
	if err != nil {
		if err == sql.ErrNoRows {
			// Jika error-nya spesifik karena baris tidak ketemu, kembailkan "nil" sebagai usernya
			return nil, fmt.Errorf("user tidak ditemukan")
		}
		// Apabila error lainnya (misal koneksi terputus)
		return nil, err
	}

	// Operator '&' digunakan untuk mengambil "ALAMAT MEMORI" lokasi data struct user yang baru diisi datanya di atas.
	// Alamat memori inilah (Pointer) yang kita serahkan kembali untuk dikonsumsi layer Service.
	return &user, nil
}

// IsUserExist dipakai service untuk rule "username unik".
func (r *userRepositoryImpl) IsUserExist(username string) bool {
	var existing string
	// Coba cari data tersebut
	err := r.DB.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existing)

	// Jika tidak ada error (err == nil), maka baris dapat ditemukan (berarti true)
	return err == nil
}

// Insert menambahkan pendaftaran profil dan password baru (telah di-hash) ke dalam tabel.
func (r *userRepositoryImpl) Insert(username, password string) error {
	// db.Exec biasa digunakan untuk query yang MENGUBAH data dan tidak mengembalikan baris (INSERT, UPDATE, DELETE).
	// Tanda "?" adalah parameter query (preparedStatement) yang aman dari peretasan SQL Injection.
	_, err := r.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	return err
}

// UpdatePassword memperbarui password untuk username terkait.
func (r *userRepositoryImpl) UpdatePassword(username, password string) error {
	_, err := r.DB.Exec(
		"UPDATE users SET password = ? WHERE username = ?",
		password,
		username,
	)
	return err
}

// GetProfile mengambil data username dan dibungkus di kembalian saja.
func (r *userRepositoryImpl) GetProfile(id int, username string) (int, string, error) {
	var userID int
	var user string

	err := r.DB.QueryRow("SELECT id, username FROM users WHERE id = ?", id).Scan(&userID, &user)

	return userID, user, err
}
