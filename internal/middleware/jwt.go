// Package middleware berada di tengah-tengah perjalanan lalu lintas web.
// File jwt ini berguna untuk menyetop akses orang yang tidak berhak lanjut ke rute yang dilindungi.
package middleware

import (
	"context"
	"golearn-structured/internal/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWT merupakan fungsi dekorator pemblokir untuk memastikan rute memiliki token aktif.
// Middleware menerima Handler asli (next) lalu merapikannya dalam fungsi HTTP lain baru.
func JWT(secret []byte) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		// Inilah handler yang kita kembalikan secara "dibungkus"
		return func(w http.ResponseWriter, r *http.Request) {
			// 1. Cek isi text header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// JIKA dikosongkan (tanpa token), kita tolak akses.
				// http.Error(w, "Token tidak ada", http.StatusUnauthorized)
				utils.Error(w, http.StatusUnauthorized, "Token tidak ada")
				return // STOP, rute aslinya tidak dipanggil
			}

			// Format yang benar adalah "Bearer <tokennya>". Hapus tulisan Bearernya agar bersih.
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// 2. Parse mendeteksi isi token, dan memastikan kata sandinya (secret) sama persis!
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return secret, nil
			})

			// 3. Jika tokennya palsu, sudah kadaluwarsa, atau error syntax, maka diblokir
			if err != nil || !token.Valid {
				// http.Error(w, "Token tidak valid", http.StatusUnauthorized)
				utils.Error(w, http.StatusUnauthorized, "Token tidak valid")
				return
			}

			// 4. Baca catatan klaim (isinya saat dibuat di fungsi Login Service)
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				// http.Error(w, "Token tidak Valid", http.StatusUnauthorized)
				utils.Error(w, http.StatusUnauthorized, "Token tidak valid")
				return
			}

			// Ambil bagian nilai klaim `username` lalu dikonversi paksa jadi string.
			username, ok := claims["username"].(string)
			if !ok {
				utils.Error(w, http.StatusUnauthorized, "Token tidak valid")
				// http.Error(w, "Token tidak valid", http.StatusUnauthorized)
				return
			}

			// 5. Setelah lulus semua filter, kita selipkan nilai variabel `username` ke dalam "Konteks Latar" (Context).
			// Supaya handler rute asli selanjutnya tahu username siapa yang sedang membuka URL ini.
			ctx := context.WithValue(r.Context(), "username", username)

			// 6. Jalankan "next" yang berarti Handler inti (misal file handler.Profile) dipersilahkan berjalan dengan context lengkap.
			next(w, r.WithContext(ctx))
		}
	}
}
