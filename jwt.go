// file ini adalah versi MONOLITIK alias "Unstructured", sengaja dibuat untuk perbandingan.
// Semua logic, koneksi database, logic HTTP, verifikasi token, disatukan dalam satu file Go.
// Ini BURUK untuk proyek besar karena menyusahkan pembacaan tim dan testing code, namun kadang dirasa LEBIH CEPAT dimengerti di level sederhana untuk pemula sebelum mereka belajar "Layered Clean Architecture".

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Menyiapkan cetakan pembantu struktur login. Di model terstruktur, ditaruh di folder model/
type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UpdateUserReq struct {
	Password string `json:"password"`
}

// Titik utama masuk aplikasi
func main() {
	connectDB() // Menancapkan basis data MySQL

	// Endpoint rute kita yang langsung didaftarkan tanpa file tambahan
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/profile", jwtMiddleware(profileHandler)) // Pakai middleware fungsi
	http.HandleFunc("/update", jwtMiddleware(updateUserHandler))
	
	fmt.Println("Server jalan di :8080")
	http.ListenAndServe(":8080", nil) // Server Nyala 
}

// Token Rahasia Kunci Masuk Server (Secret Key)
var jwtSecret = []byte("secret123")

// Fungsi kecil pembantu mencetak pesan eror seragam berbentuk JSON.
func jsonError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// Menangani URL pendaftaran di /login.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Pastikan hanya request POST yang diperbolehkan
	if r.Method != "POST" {
		jsonError(w, http.StatusMethodNotAllowed, "Method harus POST")
		return
	}

	var req LoginReq
	// Terjemahkan byte HTTP menjadi Golang STRUCT data JSON.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// ===========================
	// VALIDASI DULU (SEBELUM QUERY DB). Layaknya file utils/service validation
	// ===========================
	if req.Username == "" {
		jsonError(w, http.StatusBadRequest, "Username wajib diisi")
		return
	}
	if req.Password == "" {
		jsonError(w, http.StatusBadRequest, "Password wajib diisi")
		return
	}

	// ---------------------------
	// Bagian memanggil DB (Kalo yg layer, ini masuknya repository database query)
	// ---------------------------
	username, hashedPassword, err := getUserByUsername(req.Username)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "User tidak ditemukan")
		return
	}

	// Cek Password Hash. Apakah password sama?
	if err := checkPassword(hashedPassword, req.Password); err != nil {
		jsonError(w, http.StatusUnauthorized, "Password salah")
		return
	}

	// Bikin Token-nya karena aman bisa lolos sejauh ini
	token := generateToken(username)

	// Kembalikan Tokennya kepada End User agar menempel di browser Postman.
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// Menghubungi database buat dapetin profil utuh
func getUserByUsername(username string) (string, string, error) {
	var user, password string

	query := "SELECT username, password FROM users WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&user, &password)

	return user, password, err
}

// Menggunakan komparator standar sandi Bcrypt.
func checkPassword(hashed string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

// Meracik bahan resep map claim JWT buat ditandatangani.
func generateToken(username string) string {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)

	return tokenString
}

// Middleware pelindung. Mirip kode middleware struktur yang kita buat untuk menyaring Bearer Header
func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsonError(w, http.StatusUnauthorized, "Token tidak ada")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			jsonError(w, http.StatusUnauthorized, "Token tidak valid")
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Nitip inject nama profil dari tokennya ke konteks Request
		ctx := context.WithValue(r.Context(), "username", username)
		r = r.WithContext(ctx)
		
		// Lanjutkan request jalan!
		next(w, r)
	}
}

// Rute buat nampilin profil si token yang valid
func profileHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string) // Ambil lagi nama profil tadi
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "ini data rahasia",
		"username": username,
	})
}

// Variabel global db memegang kunci akses pool query
var db *sql.DB

func connectDB() {
	dsn := "root:Gridy20@tcp(127.0.0.1:3306)/testdb"

	database, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = database.Ping()
	if err != nil {
		panic(err)
	}
	db = database
}

// Untuk daftarin warga
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		jsonError(w, http.StatusMethodNotAllowed, "Method harus POST")
		return
	}

	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if err := validateRegister(req); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	if isUserExist(req.Username) {
		jsonError(w, http.StatusBadRequest, "username sudah digunakan")
		return
	}

	hashed, err := hashPassword(req.Password)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "gagal hash password")
		return
	}

	if err := insertUser(req.Username, hashed); err != nil {
		jsonError(w, http.StatusInternalServerError, "gagal insert user")
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "user berhasil dibuat",
	})
}

func validateRegister(req LoginReq) error {
	if req.Username == "" {
		return fmt.Errorf("username wajib diisi")
	}
	if req.Password == "" {
		return fmt.Errorf("password wajib diisi")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password minimal 6 karakter")
	}
	return nil
}

func isUserExist(username string) bool {
	var existing string
	err := db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existing)
	return err == nil
}

func insertUser(username, password string) error {
	query := "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err := db.Exec(query, username, password)
	return err
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		jsonError(w, http.StatusMethodNotAllowed, "Method harus PUT")
		return
	}

	username := r.Context().Value("username").(string)

	var req UpdateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Password == "" {
		jsonError(w, http.StatusBadRequest, "password wajib diisi")
		return
	}

	hashed, err := hashPassword(req.Password)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "gagal hash password")
		return
	}

	if err := updatePassword(username, hashed); err != nil {
		jsonError(w, http.StatusInternalServerError, "gagal update user")
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "password berhasil diupdate",
	})
}

func updatePassword(username, password string) error {
	query := "UPDATE users SET password = ? WHERE username = ?"
	_, err := db.Exec(query, password, username)
	return err
}
