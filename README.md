# Taskive Backend

Backend service untuk aplikasi manajemen tugas Taskive menggunakan Go, PostgreSQL, dan Gin framework.

## Teknologi yang Digunakan

- Go 1.20+
- PostgreSQL (via GORM)
- Gin Framework
- JWT untuk autentikasi
- GORM sebagai ORM
- go-playground/validator untuk validasi request

## Struktur Proyek

```
backend/
├── config/         # Konfigurasi aplikasi dan database
├── controllers/    # HTTP request handlers
├── models/         # Model database dan struktur data
├── routes/         # Definisi routing
├── middlewares/    # Middleware (auth, logging, dll)
├── services/       # Business logic
├── utils/          # Helper functions
├── .env           # Environment variables
├── go.mod         # Go modules
└── main.go        # Entry point
```

## Fitur

1. **Manajemen User**
   - Register
   - Login (JWT)

2. **Manajemen Project**
   - CRUD project
   - Invite member
   - Role-based access control

3. **Manajemen Task**
   - CRUD task
   - Update status task
   - Assign task ke user
   - Prioritas task

4. **Komentar**
   - Tambah komentar ke task
   - Lihat komentar per task
   - Hapus komentar

## Setup Development

1. Clone repository
2. Copy `.env.example` ke `.env` dan sesuaikan konfigurasi
3. Pastikan PostgreSQL sudah berjalan
4. Buat database baru sesuai `DB_NAME` di `.env`
5. Install dependencies:
   ```bash
   go mod download
   ```
6. Jalankan aplikasi:
   ```bash
   go run main.go
   ```

## API Endpoints

### Auth
- `POST /auth/register` - Register user baru
- `POST /auth/login` - Login user

### Projects
- `GET /api/projects` - List semua project user
- `POST /api/projects` - Buat project baru
- `GET /api/projects/:id` - Detail project
- `PUT /api/projects/:id` - Update project
- `DELETE /api/projects/:id` - Hapus project
- `POST /api/projects/:id/invite` - Invite member ke project

### Tasks
- `GET /api/projects/:id/tasks` - List task dalam project
- `POST /api/projects/:id/tasks` - Buat task baru
- `GET /api/tasks/:id` - Detail task
- `PUT /api/tasks/:id` - Update task
- `DELETE /api/tasks/:id` - Hapus task
- `PATCH /api/tasks/:id/status` - Update status task

### Comments
- `GET /api/tasks/:id/comments` - List komentar dalam task
- `POST /api/tasks/:id/comments` - Tambah komentar
- `DELETE /api/comments/:id` - Hapus komentar

## Kontribusi

1. Fork repository
2. Buat branch baru (`git checkout -b feature/amazing-feature`)
3. Commit perubahan (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request 