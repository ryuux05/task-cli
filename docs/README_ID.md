# Task CLI

Aplikasi manajemen tugas berbasis command-line yang powerful, dibangun dengan Go. Kelola tugas-tugas Anda dengan mudah menggunakan antarmuka command-line dan tampilan HTML yang indah.

## Fitur

- Buat, lihat, perbarui, dan hapus tugas
- Menandai tugas sebagai selesai
- Daftar semua tugas dengan opsi penyaringan
- Lihat tugas secara individual atau semua tugas dalam format HTML di browser default
- Mode CLI interaktif dengan riwayat perintah
- UI responsif modern untuk tampilan HTML dengan kemampuan penyaringan dan pengurutan

## Instalasi

### Prasyarat

- Go 1.16 atau lebih tinggi
- SQLite3

### Membangun dari Sumber

1. Clone repositori:
   ```
   git clone https://github.com/ryuux05/task-cli.git
   cd task-cli
   ```

2. Bangun aplikasi:
   ```
   make build
   ```
   
   Ini akan membuat binary di direktori `bin`.

3. Jalankan migrasi database:
   ```
   ./bin/migrate up
   ```

## Penggunaan

Task CLI mendukung eksekusi perintah langsung dan mode interaktif.

### Eksekusi Perintah Langsung

Jalankan perintah secara langsung:

```
./bin/task <perintah> [argumen]
```

### Mode Interaktif

Mulai mode CLI interaktif:

```
./bin/task
```

Dalam mode interaktif, Anda akan melihat prompt `>` di mana Anda dapat memasukkan perintah.

## Perintah

### Menambahkan Tugas

Tambahkan tugas baru:

```
task add "Selesaikan dokumentasi proyek"
```

Dalam mode interaktif:
```
> add Selesaikan dokumentasi proyek
```

### Melihat Daftar Tugas

Lihat semua tugas:

```
task list
```

Lihat hanya tugas yang selesai:
```
task list -c
```

Lihat semua tugas (termasuk yang selesai):
```
task list -a
```

### Menyelesaikan Tugas

Tandai tugas sebagai selesai:

```
task done <id_tugas>
```

Contoh:
```
task done 1
```

### Memperbarui Tugas

Perbarui deskripsi tugas:

```
task update <id_tugas> "Deskripsi tugas baru"
```

Perbarui tugas dan tandai sebagai selesai:
```
task update <id_tugas> -c "Deskripsi tugas baru"
```

### Melihat Tugas

Lihat tugas tunggal dalam format HTML (buka di browser):

```
task view <id_tugas>
```

Lihat tugas dalam format teks:
```
task view <id_tugas> --format text
```

Lihat semua tugas dalam format HTML (buka di browser):
```
task view-all
```

Lihat semua tugas dalam format teks:
```
task view-all --format text
```

### Menghapus Tugas

Hapus tugas:

```
task delete <id_tugas>
```

## Fitur Tampilan HTML

Saat melihat tugas dalam format HTML, Anda dapat:

- Menyaring tugas berdasarkan nama menggunakan kotak pencarian
- Menyaring tugas berdasarkan status (Semua, Tertunda, Selesai)
- Beralih status tugas antara Tertunda dan Selesai
- Mengedit detail tugas
- Menghapus tugas
- Melihat jumlah total tugas/tugas yang difilter

## Struktur Proyek

```
task-cli/
├── bin/                 # Binary executable
├── cmd/                 # Aplikasi command-line
│   ├── migrate/         # Tool migrasi database
│   └── task/            # Aplikasi utama task CLI
├── db/                  # Kode terkait database dan migrasi
├── public/              # Aset publik
│   ├── assets/          # Aset statis (CSS, JS)
│   └── templates/       # Template HTML
├── storage/             # Penyimpanan persisten
└── task/                # Logika manajemen tugas inti
```

## Pengembangan

### Menjalankan Test

```
make test
```

### Membangun untuk Pengembangan

```
make dev
```

### Menjalankan Linter

```
make lint
```

## Lisensi

Proyek ini dilisensikan di bawah Lisensi MIT - lihat file LICENSE untuk detailnya.

## Kontribusi

Kontribusi selalu disambut! Silakan ajukan Pull Request.

1. Fork repositori
2. Buat branch fitur Anda (`git checkout -b feature/fitur-keren`)
3. Commit perubahan Anda (`git commit -m 'Menambahkan fitur keren'`)
4. Push ke branch (`git push origin feature/fitur-keren`)
5. Buka Pull Request

## Ucapan Terima Kasih

- Dibangun dengan [Go](https://golang.org/)
- UI dibangun dengan [Bootstrap](https://getbootstrap.com/) 