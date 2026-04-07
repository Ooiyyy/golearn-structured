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
	// 1. Buat cetakan tipe data anonim untuk menampung variabel JSON dari input body
	var req dto.AuthloginRequest

	// 2. Decode atau menerjemahkan nilai payload request "r.Body" bahasa JSON HTTP agar dimengerti menjadi variabel struct Go bernama "&req"
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		// Error 400 Bad Request jika gagal dibaca
		return
	}

	// 3. Masukkan datanya dan jalankan Logika Bisnis (Login) di file Service
	token, err := h.service.Login(req.Username, req.Password, []byte(h.jwtSecret))
	if err != nil {
		// Kalau gagal login akibat password salah dari service, kembalikan 401 (Tidak Diizinkan Akses)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 4. Kalau berhasil dapat token sakti, ubahlah menjadi JSON (Encode) ke `w` agar browser pengguna bisa menerimanya!
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"token": token,
	// })
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

	// Meminta service membanting tulang untuk daftarkan user
	err := h.service.Register(req.Username, req.Password)
	if err != nil {
		// http.Error(w, err.Error(), 400)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Bila mulus tanpa eror, berikan pesan sukses.
	c.JSON(http.StatusOK, gin.H{
		"message": "user berhasil dibuat",
	})
}

// Profile menampilkan "Siapa yang sedang login ini" sesudahnya token dilampirkan.
func (h *AuthHandler) Profile(c *gin.Context) {
	// Ingat variabel r.Context() yang didekorasi pada Middleware JWT kita?
	// Kini bagian profil bisa menyadap `username` dari latar belakang Konteks data tersebut dengan pasti!
	// Kode `.(string)` dinamakan "Type Assertion", memastikan pemaksaan tipenya adalah string.
	id := c.GetInt("id")
	username := c.GetString("username")

	// Ambil keterangan data pribadi dari MySQL lewat jembatan Service -> Repository.
	id, user, err := h.service.GetProfile(id, username)
	if err != nil {
		// Kali ini menggunakan utils dari fungsi utils/response.go yang dibuat spesifik agar lebih estetik format JSON kembaliannya.
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Gunakan alat tolong kembalikan Success dengan memanggil helper utils.
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
