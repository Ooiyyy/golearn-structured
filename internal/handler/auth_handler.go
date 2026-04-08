// Package handler berada di lapisan terdepan paling luar untuk lalu lintas data internet (jembatan Request dan Service).
// Tugasnya mencerna data JSON yang dikirim via Postman/Browser, menyuruh file Service kerja, lalu mengirim kembali tulisan hasil kerja JSON untuk dilihat user.
package handler

import (
	"golearn-structured/internal/dto"
	"golearn-structured/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler mendefinisikan bentuk sebuah kumpulan rute Handler untuk Autentikasi (Pendaftaran dsb).
type AuthHandler struct {
	service   *service.UserService // Logika bisnis mengecek login ada di dalam ruang kerja service ini
	jwtSecret string               // Kata kunci JWT yang harus disetor saat pengaksesan token login.
}

// NewAuthHandler adalah pembuat instansi tipe object AuthHandler secara mudah.
func NewAuthHandler(s *service.UserService, jwtSecret string) *AuthHandler {
	return &AuthHandler{service: s, jwtSecret: jwtSecret}
}

// Login merupakan handler (petugas yang menangani rute) ketika URL rute /login dipanggil.
// Parameter `w http.ResponseWriter` adalah surat balasan dan `r *http.Request` adalah surat pertanyaan dari Browser klien.
func (h *AuthHandler) Login(c *gin.Context) {
	// Request pertama kali diterima di handler lewat gin.Context (c).
	var req dto.AuthloginRequest

	// Bind JSON body ke DTO; pointer (&req) dipakai agar Gin mengisi struct langsung.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Delegasikan business logic ke service: handler fokus pada HTTP I/O.
	token, err := h.service.Login(req.Username, req.Password, []byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Response dikembalikan ke client dari sini dalam format JSON.
	c.JSON(http.StatusOK, gin.H{
		"message": "login berhasil",
		"token":   token,
	})
}

// Register untuk mendaftarkan akun.
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Service mengeksekusi validasi bisnis + simpan data via repository.
	err := h.service.Register(req.Username, req.Password)
	if err != nil {
		// http.Error(w, err.Error(), 400)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Titik akhir lifecycle endpoint register: kirim status sukses.
	c.JSON(http.StatusOK, gin.H{
		"message": "user berhasil dibuat",
	})
}

// Profile menampilkan "Siapa yang sedang login ini" sesudahnya token dilampirkan.
func (h *AuthHandler) Profile(c *gin.Context) {
	// Data user ini disisipkan middleware JWT setelah token valid.
	id := c.GetInt("id")
	username := c.GetString("username")

	// Alur tetap konsisten: handler -> service -> repository.
	id, user, err := h.service.GetProfile(id, username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Response profile dikembalikan sebagai JSON object bersarang.
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":       id,
			"username": user,
		},
	})
}

// Logout digunakan untuk merespon bahwa mematikan sistem
func (h *AuthHandler) Logout(c *gin.Context) {
	// Note: di JWT, logout idealnya dihapus di level browser (Storage), back-end cukup membalas sukses kalau tidak pakai redis!
	c.JSON(http.StatusOK, gin.H{
		"message": "logout berhasil",
	})
}
