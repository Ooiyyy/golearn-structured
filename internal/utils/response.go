// Package utils berisi fungsi-fungsi kecil penolong (utility) yang bisa dipakai berulang kali.
// File ini khusus memuat cara mengirimkan response (jawaban HTTP JSON) agar standar.
package utils

import (
	"encoding/json"
	"net/http"
)

// JSON adalah fungsi utama untuk memformat output web sebagai tipe JSON.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	// Memberitahukan browser atau Postman bahwa data yang akan didapatkan adalah aplikasi JSON
	w.Header().Set("Content-Type", "application/json")
	
	// Memberitahukan Status Code HTTP (Misal 200 = OK, 400 = Bad Request, dsb)
	w.WriteHeader(status)
	
	// Mengubah variabel 'data' bertipe apa saja (interface{}) menjadi format string JSON.
	json.NewEncoder(w).Encode(data)
}

// Error adalah fungsi penolong untuk membalas JSON error cepat.
func Error(w http.ResponseWriter, status int, message string) {
	// Akan menjawab dalam format: {"error": "pesan error"}
	JSON(w, status, map[string]string{
		"error": message,
	})
}

// Success adalah fungsi penolong saat tidak ada halangan dan memberikan status 200 OK standar.
func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}
