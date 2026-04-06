// Package main adalah titik masuk (entry point) utama dari aplikasi Go ini.
package main

import (
	"fmt"
	"golearn-structured/config"
	"golearn-structured/internal/handler"
	"golearn-structured/internal/middleware"
	"golearn-structured/internal/repository"
	"golearn-structured/internal/service"

	"net/http"
)

// jwtSecret adalah kunci rahasia yang digunakan untuk membuat dan memvalidasi token JWT.
// PERHATIAN: Di dunia nyata, ini tidak boleh ditulis langsung di kode, melainkan dari environment variable (misal: file .env)!
var jwtSecret = []byte("secret123")

// Fungsi main akan dijalankan pertama kali secara otomatis saat aplikasi dimulai.
func main() {
	// 1. Menghubungkan ke database MySQL
	db := config.Connect()
	
	// 2. Inisiasi Repository (Komponen yang bertugas langsung berinteraksi dengan database / SQL)
	repo := repository.NewUserRepository(db)
	
	// 3. Inisiasi Service (Komponen yang bertugas menangani logika bisnis, seperti validasi password, hashing)
	service := service.NewUserService(repo)
	
	// 4. Inisiasi Handler (Komponen yang bertugas menerima dan membalas request HTTP dari user)
	handler := handler.NewAuthHandler(service, jwtSecret)

	// 5. Mendaftarkan rute-rute (endpoint) aplikasi kita
	
	// Rute Publik (tidak perlu login)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/register", handler.Register)
	
	// Rute Privat (WAJIB login, maka dibungkus dilindungi dengan func middleware.JWT)
	http.HandleFunc("/profile", middleware.JWT(jwtSecret)(handler.Profile))
	http.HandleFunc("/logout", middleware.JWT(jwtSecret)(handler.Logout))

	// 6. Menjalankan server HTTP pada port 8080
	fmt.Println("Server jalan di :8080")
	// ListenAndServe akan membuat program berjalan terus untuk mendengarkan koneksi yang masuk
	http.ListenAndServe(":8080", nil) 
}
