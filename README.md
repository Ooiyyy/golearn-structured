# Panduan Struktur Folder GoLearn-Structured untuk Pemula 🚀

Selamat datang di kode proyek ini! Proyek ini menggunakan pola tata letak file mirip dengan **Clean Architecture** sederhana (atau MVC berlapis). 
Tujuan dari membagi-bagi file ke dalam beberapa folder adalah agar kode kita lebih rapi, gampang dites secara individual, dan lebih mudah dikembangkan jika aplikasinya semakan membesar meski di masa depan jumlah foldernya bertambah ratusan.

Di bawah ini adalah penjelasan sederhana khusus disediakan untuk *Pemula Go* dari setiap folder yang ada di proyek ini:

## 📁 `cmd/`
- **Fungsi:** Titik utama di mana program ini mulai dijalankan `(go run cmd/main.go)`.
- **Isi:** Biasanya berisi file `main.go`. Isinya nyaris cuma deklarasi variabel rahasia, pemanggilan koneksi database, inisiasi modul lain, lalu "dijahit" rute endpointnya bersama-sama.
- **Analogi:** Seperti "Gerbang Pintu Masuk" dari sebuah pabrikk.

## 📁 `config/`
- **Fungsi:** Berisi konfigurasi global aplikasi.
- **Isi:** Pada kode ini, berisi `database_config.go` untuk menghubungi MySQL.
- **Analogi:** "Ruang Panel Mesin Listrik", tempat pengaturan kabel dan aliran dasar ke seluruh pabrik diaktifkan.

## 📁 `internal/`
Ini adalah folder **rahasia** dalam tradisi Go. Di bahasa program Go, folder bernama `internal` itu punya sifat khusus: isinya tidak bisa diimpor sembarangan/dipungut oleh orang luar dan proyek Github lainnya. Inti dari semua rahasia logika perusahaan/aplikasi diletakkan di bawah payung folder `internal`.

### 📂 `internal/handler/` (Biasa juga disebut Controllers)
- **Fungsi:** Tempat bertemunya HTTP Request JSON (dari Postman/Browser) dengan sistem di pabrik kita.
- **Isi:** Cenderung mengecek validitas format input (JSON dibongkar menjadi variabel Struct), dan mencetak output jawaban ke browser. **Tidak boleh** mengerjakan aturan logika bisnis di dalamnya secara berat.
- **Analogi:** "Pelayan Restoran". Dia hanya mencatat pesanan pelanggan pakai kertas lalu menyerahkannya ke dapur (service), setelah makanan jadi pelayan mengantar hasilnya ke pelanggan.

### 📂 `internal/middleware/`
- **Fungsi:** Pasukan pencegat (saringan pos jaga satpam) yang berada di tengah rute.
- **Isi:** Mengecek apakah token authorization itu valid (seperti `jwt.go`). Jika valid, diperbolehkan masuk. Jika ketahuan palsu, dicegat dengan status 401 Unauthorized sebelum rutenya terpanggil di Handler.
- **Analogi:** "Satpam Resepsionis" yang selalu menanyakan KTP atau Kartu Akses sebelum kamu diizinkan menerobos masuk ke Lift area privat.

### 📂 `internal/service/` (Atau biasa disebut Usecase/Business Logic)
- **Fungsi:** "Otak" utama dari aplikasi. Tempat menaruh keputusan rumit.
- **Isi:** Disinilah diputuskan aturan bahwa *"password harus minimal 6 karakter"*, tempat di mana password asli kamu direkatkan dan di-*hash* acak (bcrypt), sampai validasi registrasi lainnya ada di sini.
- **Analogi:** "Koki Spesialis Dapur". Dialah yang tahu resepnya dan benar-benar memproses makanan mentah (data) yang dibikin pelanggan ke matang.

### 📂 `internal/repository/`
- **Fungsi:** Tempat di mana seluruh lalu lintas yang menuju ke Database MySQL tinggal!
- **Isi:** Seluruh perintah SQL mentah dioperasikan di sini seperti `SELECT, UPDATE, INSERT, DELETE`. Repository itu cuek dan lugu, ia sama sekali tidak peduli HTTP internet, tidak pula ia peduli soal batasan password login di Service. Ia cuma petugas catat dan cari.
- **Analogi:** "Petugas Gudang Rak Buku". Koki di service bisa bilang padanya: *"Tolong carikan username bernama 'Budi'"*, lalu Petugas Gudang (Repository) yang akan ke belakang mencari dengan query SQL! 

### 📂 `internal/utils/` (Kadang dinamakan helper/pkg)
- **Fungsi:** File penolong serba-guna. Tempat menaruh kumpulan fungsi mungil yang bisa dipanggil berkali-kali di folder mana saja.
- **Isi:** Misal `response.go` yang gunanya merapikan hasil akhir format JSON agar standar bentuknya selalu seragam tiap respons Error atau Sukses-nya.
- **Analogi:** "Pisau Dapur Multifungsi", selalu tersedia di sudut manapun siap dipakai tanpa mikir panjang.

---

## 📄 File Spesial Lainnya
- **`go.mod` dan `go.sum`**: Kartu identitas utama proyek ini, berisi daftar pustaka modul eksternal (Dependency seperti go-sql-driver dan bcrypt JWT) yang kamu unduh dari internet.
- **`jwt.go` (File terluar)**: Ini adalah versi proyek ini dengan seluruh logika dijadikan SATU FILE SAJA! Pendekatan "Satu file" seperti ini lazim disebut pola monolitik ekstrem tidak terstruktur. Sengaja disediakan buat perbandingan visual para Pemula Go (bandingkan isi jwt.go vs semua folder di /internal!), agar terbiasa bahwa menggunakan folder *(Clean Architecture)* jauh lebih rapi dibanding satu file dijejalkan ribuan kode!
