# Request Lifecycle (Handler -> Service -> Repository)

Dokumen ini menjawab dua hal utama:
- request diterima di mana,
- response dikembalikan di mana dan bentuknya seperti apa.

## 1. Request diterima di mana?

Request HTTP pertama kali masuk lewat router Gin di `internal/routes`.
Router memetakan URL ke handler, contoh:
- `POST /login` -> `AuthHandler.Login`
- `GET /profile` -> JWT middleware -> `AuthHandler.Profile`

Jadi, titik masuk request ada di:
- `internal/routes/*.go` (pendaftaran endpoint),
- `internal/handler/*.go` (fungsi handler sebagai penerima request aktual).

## 2. Handler melakukan apa?

Handler (`internal/handler`) adalah layer HTTP:
- membaca input (`ShouldBindJSON`, `PostForm`, `Query`, `Param`),
- validasi format request level transport,
- memanggil service untuk logika bisnis,
- mengirim response JSON via `c.JSON(statusCode, payload)`.

Contoh bentuk response:
- sukses login:
  - status: `200`
  - body: `{"message":"login berhasil","token":"..."}`
- gagal login:
  - status: `401`
  - body: `{"error":"username atau password salah"}`

## 3. Service melakukan apa?

Service (`internal/service`) adalah layer aturan bisnis:
- validasi rule domain (mis. password minimal 6 karakter),
- enkripsi/verifikasi password (`bcrypt`),
- bikin JWT token,
- aturan update todo (fallback nilai lama jika field kosong).

Service tidak tahu HTTP, jadi return data/error ke handler.

## 4. Repository melakukan apa?

Repository (`internal/repository`) adalah layer persistence:
- menjalankan query SQL (`QueryRow`, `Query`, `Exec`),
- mengubah hasil SQL jadi data Go (`Scan`),
- tidak mengurus request/response HTTP.

## 5. Contoh alur end-to-end: `POST /login`

1. Client kirim JSON ke `/login`.
2. Router arahkan ke `AuthHandler.Login`.
3. Handler bind JSON ke DTO.
4. Handler panggil `UserService.Login`.
5. Service panggil `UserRepository.GetByUsername`.
6. Repository query DB, return user + hash password.
7. Service verifikasi bcrypt, buat JWT token, return token.
8. Handler kirim JSON response ke client (`200` + token).

## 6. Contoh alur end-to-end: `POST /api/todos`

1. Client kirim multipart/form-data ke `/api/todos`.
2. Router arahkan ke `TodoHandler.Create`.
3. Handler baca form field (`title`, `note`) dan file (`image`) bila ada.
4. Handler simpan file ke `uploads/`, lalu bentuk `imageURL`.
5. Handler panggil `TodoService.CreateTodo`.
6. Service validasi rule bisnis (judul wajib), lalu bentuk struct `model.Todo`.
7. Service panggil `TodoRepository.Create`.
8. Repository menjalankan `INSERT` ke tabel `todos`.
9. Handler kirim JSON response ke client:
   - sukses: `200` + `{"message":"todo berhasil dibuat"}`
   - gagal validasi: `400` + `{"error":"..."}`.

## 7. Konsep dasar Go yang dipakai

- `struct`: kontrak data (`model`, `dto`) dan object layer (handler/service/repository).
- `pointer`: dipakai untuk dependency injection (`*UserService`, `*sql.DB`) agar efisien dan konsisten.
- `error` as value: setiap langkah bisa mengembalikan error, lalu ditangani eksplisit.
- `closure`: middleware JWT mengembalikan fungsi yang "membawa" `jwtSecret`.
- `interface`: pola arsitekturnya siap untuk interface, tapi saat ini masih concrete type.
