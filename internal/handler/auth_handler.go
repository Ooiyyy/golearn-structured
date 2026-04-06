// Package handler berada di lapisan terdepan paling luar untuk lalu lintas data internet (jembatan Request dan Service).
// Tugasnya mencerna data JSON yang dikirim via Postman/Browser, menyuruh file Service kerja, lalu mengirim kembali tulisan hasil kerja JSON untuk dilihat user.
package handler

import (
	"encoding/json"
	"golearn-structured/internal/service"
	"golearn-structured/internal/utils"
	"net/http"
)

// AuthHandler mendefinisikan bentuk sebuah kumpulan rute Handler untuk Autentikasi (Pendaftaran dsb).
type AuthHandler struct {
	service *service.UserService // Logika bisnis mengecek login ada di dalam ruang kerja service ini
	secret  []byte               // Kata kunci JWT yang harus disetor saat pengaksesan token login.
}

// NewAuthHandler adalah pembuat instansi tipe object AuthHandler secara mudah.
func NewAuthHandler(s *service.UserService, secret []byte) *AuthHandler {
	return &AuthHandler{service: s, secret: secret}
}

// Login merupakan handler (petugas yang menangani rute) ketika URL rute /login dipanggil.
// Parameter `w http.ResponseWriter` adalah surat balasan dan `r *http.Request` adalah surat pertanyaan dari Browser klien.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// 1. Buat cetakan tipe data anonim untuk menampung variabel JSON dari input body
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 2. Decode atau menerjemahkan nilai payload request "r.Body" bahasa JSON HTTP agar dimengerti menjadi variabel struct Go bernama "&req"
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// http.Error(w, err.Error(), 400)
		utils.Error(w, http.StatusBadRequest, err.Error())
		// Error 400 Bad Request jika gagal dibaca
		return
	}

	// 3. Masukkan datanya dan jalankan Logika Bisnis (Login) di file Service
	token, err := h.service.Login(req.Username, req.Password, h.secret)
	if err != nil {
		// Kalau gagal login akibat password salah dari service, kembalikan 401 (Tidak Diizinkan Akses)
		// http.Error(w, err.Error(), 401)
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	// 4. Kalau berhasil dapat token sakti, ubahlah menjadi JSON (Encode) ke `w` agar browser pengguna bisa menerimanya!
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"token": token,
	// })
	utils.Success(w, map[string]string{
		"message": "login berhasil",
		"token":   token,
	})
}

// Register untuk mendaftarkan akun.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// http.Error(w, err.Error(), 400)
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Meminta service membanting tulang untuk daftarkan user
	err := h.service.Register(req.Username, req.Password)
	if err != nil {
		// http.Error(w, err.Error(), 400)
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Bila mulus tanpa eror, berikan pesan sukses.
	utils.Success(w, map[string]string{
		"message": "user berhasil dibuat",
	})
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"message": "user berhasil dibuat",
	// })
}

// Profile menampilkan "Siapa yang sedang login ini" sesudahnya token dilampirkan.
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	// Ingat variabel r.Context() yang didekorasi pada Middleware JWT kita?
	// Kini bagian profil bisa menyadap `username` dari latar belakang Konteks data tersebut dengan pasti!
	// Kode `.(string)` dinamakan "Type Assertion", memastikan pemaksaan tipenya adalah string.
	username := r.Context().Value("username").(string)

	// Ambil keterangan data pribadi dari MySQL lewat jembatan Service -> Repository.
	user, err := h.service.GetProfile(username)
	if err != nil {
		// Kali ini menggunakan utils dari fungsi utils/response.go yang dibuat spesifik agar lebih estetik format JSON kembaliannya.
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Gunakan alat tolong kembalikan Success dengan memanggil helper utils.
	utils.Success(w, map[string]string{
		"username": user, // Balas {"username": "nama yang disisipkan"}
	})
}

// Logout digunakan untuk merespon bahwa mematikan sistem
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Note: di JWT, logout idealnya dihapus di level browser (Storage), back-end cukup membalas sukses kalau tidak pakai redis!
	utils.Success(w, map[string]string{
		"message": "logout berhasil",
	})
}
