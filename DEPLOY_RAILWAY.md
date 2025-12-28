# Deploy ke Railway

## Langkah-langkah Deploy

### 1. Setup Database PostgreSQL di Railway
1. Buka [Railway.app](https://railway.app) dan login
2. Buat project baru
3. Klik "New" → "Database" → "Add PostgreSQL"
4. Railway akan otomatis membuat database dan memberikan environment variables

### 2. Setup Aplikasi Go API
1. Di project yang sama, klik "New" → "GitHub Repo" atau "Empty Service"
2. Jika menggunakan GitHub, hubungkan repository Anda
3. Railway akan otomatis mendeteksi Dockerfile

### 3. Konfigurasi Environment Variables
Di Railway service aplikasi Anda, tambahkan environment variables berikut:

**REQUIRED (akan otomatis dari Railway PostgreSQL):**
- `PGHOST` atau `BLUEPRINT_DB_HOST` → ambil dari Railway PostgreSQL
- `PGPORT` atau `BLUEPRINT_DB_PORT` → ambil dari Railway PostgreSQL  
- `PGDATABASE` atau `BLUEPRINT_DB_DATABASE` → ambil dari Railway PostgreSQL
- `PGUSER` atau `BLUEPRINT_DB_USERNAME` → ambil dari Railway PostgreSQL
- `PGPASSWORD` atau `BLUEPRINT_DB_PASSWORD` → ambil dari Railway PostgreSQL

**CUSTOM:**
```
PORT=5000
APP_ENV=production
BLUEPRINT_DB_SCHEMA=public
ACCESS_TOKEN_KEY=your-secret-access-token-key-change-this-in-production
REFRESH_TOKEN_KEY=your-secret-refresh-token-key-change-this-in-production
ACCESS_TOKEN_AGE=3600
```

#### Cara mudah setup environment variables:
Railway PostgreSQL secara otomatis menyediakan variable `DATABASE_URL`. Jika aplikasi Anda bisa menggunakan `DATABASE_URL`, itu lebih mudah. Jika tidak, gunakan referensi variable dari PostgreSQL service:

```
BLUEPRINT_DB_HOST=${{Postgres.PGHOST}}
BLUEPRINT_DB_PORT=${{Postgres.PGPORT}}
BLUEPRINT_DB_DATABASE=${{Postgres.PGDATABASE}}
BLUEPRINT_DB_USERNAME=${{Postgres.PGUSER}}
BLUEPRINT_DB_PASSWORD=${{Postgres.PGPASSWORD}}
```

### 4. Deploy
1. Push kode Anda ke GitHub (jika menggunakan GitHub integration)
2. Railway akan otomatis build dan deploy
3. Migrasi database akan dijalankan otomatis saat startup

### 5. Akses Aplikasi
- Railway akan memberikan public URL untuk aplikasi Anda
- Format: `https://your-app-name.up.railway.app`

## Troubleshooting

### Aplikasi tidak start
- Cek logs di Railway dashboard
- Pastikan semua environment variables sudah benar
- Pastikan PORT=5000 atau sesuai dengan yang Railway provide

### Database connection error
- Pastikan PostgreSQL service sudah running
- Cek environment variables database
- Pastikan `BLUEPRINT_DB_HOST` menggunakan internal hostname dari Railway

### Migration gagal
- Cek logs saat startup
- Pastikan folder migrations ter-copy dengan benar
- Cek koneksi database

## Tips
1. Gunakan Railway's internal networking untuk koneksi database (lebih cepat dan aman)
2. Set `APP_ENV=production` untuk production environment
3. Ganti `ACCESS_TOKEN_KEY` dan `REFRESH_TOKEN_KEY` dengan nilai yang aman
4. Monitor logs secara rutin di Railway dashboard
5. Setup custom domain jika diperlukan
6. Aktifkan auto-deployment untuk continuous deployment

## Structure File Penting
```
api-stockflow/
├── Dockerfile                 # Config Docker untuk Railway
├── .dockerignore             # File yang diabaikan saat build
├── railway.json              # Railway configuration
├── cmd/
│   ├── api/main.go          # Entry point aplikasi
│   └── migrate/migrate.go   # Migration tool
└── internal/
    └── database/
        └── migrations/       # SQL migrations
```

## Environment Variables Reference

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `PORT` | Port aplikasi berjalan | Yes | 5000 |
| `APP_ENV` | Environment (local/production) | Yes | local |
| `BLUEPRINT_DB_HOST` | PostgreSQL host | Yes | - |
| `BLUEPRINT_DB_PORT` | PostgreSQL port | Yes | 5432 |
| `BLUEPRINT_DB_DATABASE` | Database name | Yes | - |
| `BLUEPRINT_DB_USERNAME` | Database user | Yes | - |
| `BLUEPRINT_DB_PASSWORD` | Database password | Yes | - |
| `BLUEPRINT_DB_SCHEMA` | Database schema | No | public |
| `ACCESS_TOKEN_KEY` | JWT access token secret | Yes | - |
| `REFRESH_TOKEN_KEY` | JWT refresh token secret | Yes | - |
| `ACCESS_TOKEN_AGE` | Token expiry (seconds) | Yes | 3600 |
