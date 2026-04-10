// Package service berisi logika bisnis (business logic) aplikasi.
// Layer ini bertindak sebagai jembatan antara Handler (pengurus HTTP) dan Repository (pengurus Database).
// Segala aturan seperti "password harus minimal 6 karakter" ditempatkan di sini.
package service

import (
	"fmt"
	"golearn-structured/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService menyimpan pointer ke repository agar UserService bisa meminta repository menyimpan/mengambil data dari DB.
// UserService adalah lapisan business logic; tidak tahu detail HTTP.
type UserService struct {
	repo repository.UserRepository
}

// NewUserService adalah constructor (fungsi pembuat object utama) untuk UserService.
func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

// Login memvalidasi username dan password, serta mengembalikan Token JWT jika berhasil.
func (s *UserService) Login(username, password string, secret []byte) (string, error) {
	// 1. Ambil data user dari tabel berdasarkan nama.
	// Variabel 'user' di sini sekarang berwujud POINTER (hanya menyimpan alamat aslinya)
	// Pointer mengizinkan variabel user menjadi 'nil' jika database mengembalikan kosong.
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", fmt.Errorf("username atau password salah")
	}

	// 2. Bandingkan password asli yang direquest dengan password hashing yang ada di database.
	// Hebatnya di Golang, jika `user` itu Pointer, Anda tetap bisa menuliskan `user.Password` langsung!
	// Go akan diam-diam mengunjungi alamat memorinya dan menarik nilainya untuk Anda.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("username atau password salah")
	}

	// 3. Jika berhasil login, atur klaim (isi konten) untuk Token JWT ini (JWT claims).
	// Token diatur expired atau kadaluwarsa dalam 1 jam ke depan untuk keamanan.
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	// 4. Buat objek token-nya dan tanda tangani menggunakan kunci rahasia (secret key).
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("gagal membuat token login")
	}

	return tokenString, nil
}

// Register menangani logika pembuatan pengguna baru.
func (s *UserService) Register(username, password string) error {
	// Validasi awal agar data yang masuk tidak kosong.
	if username == "" {
		return fmt.Errorf("Username wajib diisi")
	}
	if password == "" {
		return fmt.Errorf("Password wajib diisi")
	}

	// Cek panjang password demi keamanan dasar (business rule logic)
	if len(password) < 6 {
		return fmt.Errorf("Password minimal 6 karakter")
	}

	// Cek apakah di database sebelumnya sudah ada orang pakai username ini
	if s.repo.IsUserExist(username) {
		return fmt.Errorf("Username sudah digunakan")
	}

	// Gunakan bcrypt untuk mengenkripsi password aslinya sebelum disimpan.
	// Kita tidak pernah boleh menyimpan password dalam mode teks biasa di database!
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal hash password")
	}

	// Panggil repository untuk benar-benar memasukkan data tersebut ke tabel MySQL
	return s.repo.Insert(username, string(hash))
}

// UpdatePassword berisi aturan bisnis kalau kamu ingin mengubah password
func (s *UserService) UpdatePassword(username, password string) error {
	if password == "" {
		return fmt.Errorf("Password wajib diisi")
	}

	// Hash password baru sebelum update ke tabel
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal hash password")
	}

	// Perintahkan repository untuk update data SQL!
	return s.repo.UpdatePassword(username, string(hash))
}

// GetProfile menarik data ringkas tentang pengguna
func (s *UserService) GetProfile(id int, username string) (int, string, error) {
	// Kita hanya perlu menyedot profil dari tabel saat ini.
	id, user, err := s.repo.GetProfile(id, username)
	if err != nil {
		return 0, "", fmt.Errorf("user tidak ditemukan")
	}
	return id, user, nil
}
