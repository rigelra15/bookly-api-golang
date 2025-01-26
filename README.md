# Bookly API

## Deskripsi
Bookly API adalah sebuah aplikasi backend yang dibuat dengan menggunakan Golang dan PostgreSQL. API ini memungkinkan pengguna untuk mengelola kategori buku, data buku, dan data pengguna dengan sistem autentikasi JWT (JSON Web Token). Proyek ini dilengkapi dengan dokumentasi Swagger untuk mempermudah developer dalam memahami dan menggunakan API ini.

---

## Fitur Utama
1. **Manajemen Kategori Buku**:
   - Tambah, baca, ubah, dan hapus kategori buku.
   - Mendapatkan semua buku berdasarkan kategori.
   - Mencegah penambahan kategori dengan nama yang sama (case-insensitive - mengabaikan huruf besar/kecil).
   - Kategori tidak dapat dihapus jika masih terdapat buku yang terkait.

2. **Manajemen Buku**:
   - Tambah, baca, ubah, dan hapus buku.
   - Mendapatkan semua buku berdasarkan kategori.
   - Mencegah penambahan buku dengan judul yang sama (case-insensitive - mengabaikan huruf besar/kecil).

3. **Manajemen Pengguna**:
   - Tambah, baca, ubah, dan hapus pengguna.
   - Login untuk mendapatkan token JWT.
   - Mencegah penambahan pengguna dengan username yang sama (case-sensitive - memperhatikan huruf besar/kecil).

4. **Autentikasi**:
   - Menggunakan JSON Web Token (JWT) untuk melindungi endpoint tertentu.

5. **Dokumentasi Swagger**:
   - Dokumentasi API yang terintegrasi dan mudah diakses.

---

## Dokumentasi API
Swagger Docs dapat diakses melalui: `http://localhost:8080/swagger/index.html` atau URL Publik yang sudah di-deploy yaitu: `https://bookly-api-golang-production.up.railway.app/swagger/index.html` - [Bookly API (Swagger)](https://bookly-api-golang-production.up.railway.app/swagger/index.html)

### Endpoint Utama

#### 1. **Users**
- `GET /api/users`: Mendapatkan semua pengguna. - **_[Perlu Autentikasi]_**
- `POST /api/users`: Menambahkan pengguna baru. - **_[Perlu Autentikasi]_**
- `GET /api/users/{id}`: Mendapatkan pengguna berdasarkan ID. - **_[Perlu Autentikasi]_**
- `PUT /api/users/{id}`: Memperbarui pengguna berdasarkan ID. - **_[Perlu Autentikasi]_**
- `DELETE /api/users/{id}`: Menghapus pengguna berdasarkan ID. - **_[Perlu Autentikasi]_**
- `POST /api/users/login`: Login pengguna untuk mendapatkan JWT token.

#### 2. **Categories**
- `GET /api/categories`: Mendapatkan semua kategori. - **_[Perlu Autentikasi]_**
- `POST /api/categories`: Menambahkan kategori baru. - **_[Perlu Autentikasi]_**
- `GET /api/categories/{id}`: Mendapatkan kategori berdasarkan ID. - **_[Perlu Autentikasi]_**
- `PUT /api/categories/{id}`: Memperbarui kategori berdasarkan ID. - **_[Perlu Autentikasi]_**
- `DELETE /api/categories/{id}`: Menghapus kategori berdasarkan ID. - **_[Perlu Autentikasi]_**
- `GET /api/categories/{id}/books`: Mendapatkan semua buku berdasarkan kategori. - **_[Perlu Autentikasi]_**

#### 3. **Books**
- `GET /api/books`: Mendapatkan semua buku. - **_[Perlu Autentikasi]_**
- `POST /api/books`: Menambahkan buku baru. - **_[Perlu Autentikasi]_**
- `GET /api/books/{id}`: Mendapatkan buku berdasarkan ID. - **_[Perlu Autentikasi]_**
- `PUT /api/books/{id}`: Memperbarui buku berdasarkan ID. - **_[Perlu Autentikasi]_**
- `DELETE /api/books/{id}`: Menghapus buku berdasarkan ID. - **_[Perlu Autentikasi]_**

---

## Autentikasi
Gunakan header berikut untuk endpoint yang memerlukan autentikasi Bearer Token:
```
Authorization: Bearer <token-jwt>
```
Silahkan membuat pengguna baru terlebih dahulu melalui endpoint `POST /api/users`. Token dapat diperoleh dengan login melalui endpoint `POST /api/users/login`. 

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