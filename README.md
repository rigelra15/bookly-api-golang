# Bookly API

## Deskripsi
Bookly API adalah sebuah aplikasi backend yang dibuat dengan menggunakan Golang dan PostgreSQL. API ini memungkinkan pengguna untuk mengelola kategori buku, data buku, dan data pengguna dengan sistem autentikasi JWT (JSON Web Token). Proyek ini dilengkapi dengan dokumentasi Swagger untuk mempermudah developer dalam memahami dan menggunakan API ini.

---

## Fitur Utama
1. **Manajemen Kategori Buku**:
   - Tambah, baca, ubah, dan hapus kategori buku.
   - Mendapatkan semua buku berdasarkan kategori.

2. **Manajemen Buku**:
   - Tambah, baca, ubah, dan hapus buku.

3. **Manajemen Pengguna**:
   - Tambah, baca, ubah, dan hapus pengguna.
   - Login untuk mendapatkan token JWT.

4. **Autentikasi**:
   - Menggunakan JSON Web Token (JWT) untuk melindungi endpoint tertentu.

5. **Dokumentasi Swagger**:
   - Dokumentasi API yang terintegrasi dan mudah diakses.

---

## Teknologi yang Digunakan
- **Golang**: Bahasa pemrograman untuk pengembangan aplikasi backend.
- **PostgreSQL**: Database relasional.
- **Gin**: Framework web untuk Golang.
- **Swagger**: Dokumentasi API interaktif.
- **bcrypt**: Untuk hashing password pengguna.

---

## Instalasi dan Penggunaan

### 1. Clone Repository
```bash
git clone <repository-url>
cd bookly-api-golang
```

### 2. Konfigurasi Database
Buat file `.env` di folder `config` dengan isi seperti berikut:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=bookly_go
```

### 3. Instalasi Dependensi
```bash
go mod tidy
```

### 4. Jalankan Migrasi Database
```bash
go run main.go
```

### 5. Menjalankan Server
```bash
go run main.go
```
Server akan berjalan di `http://localhost:8080`

---

## Dokumentasi API
Swagger Docs dapat diakses melalui: `http://localhost:8080/swagger/index.html` atau URL Publik yang sudah di-deploy yaitu: `https://bookly-api-golang-production.up.railway.app/swagger/index.html` - [Bookly API (Swagger)](https://bookly-api-golang-production.up.railway.app/swagger/index.html)

### Contoh Endpoint Utama

#### 1. **Users**
- `GET /api/users`: Mendapatkan semua pengguna.
- `GET /api/users/{id}`: Mendapatkan pengguna berdasarkan ID.
- `PUT /api/users/{id}`: Memperbarui pengguna berdasarkan ID.
- `DELETE /api/users/{id}`: Menghapus pengguna berdasarkan ID.
- `POST /api/users`: Menambahkan pengguna baru.
- `POST /api/users/login`: Login pengguna untuk mendapatkan JWT token.

#### 2. **Categories**
- `GET /api/categories`: Mendapatkan semua kategori.
- `POST /api/categories`: Menambahkan kategori baru.
- `PUT /api/categories/{id}`: Memperbarui kategori berdasarkan ID.
- `DELETE /api/categories/{id}`: Menghapus kategori berdasarkan ID.

#### 3. **Books**
- `GET /api/books`: Mendapatkan semua buku.
- `POST /api/books`: Menambahkan buku baru.
- `PUT /api/books/{id}`: Memperbarui buku berdasarkan ID.
- `DELETE /api/books/{id}`: Menghapus buku berdasarkan ID.

---

## Autentikasi
Gunakan header berikut untuk endpoint yang memerlukan autentikasi Bearer Token:
```
Authorization: Bearer <token-jwt>
```
Token dapat diperoleh dengan login melalui endpoint `POST /api/users/login`.

---