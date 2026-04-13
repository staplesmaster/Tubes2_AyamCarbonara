# Tubes2_AyamCarbonara

**IF2211 Strategi Algoritma — Semester II 2025/2026**  
Institut Teknologi Bandung

---

## Daftar Isi

- [Deskripsi Singkat](#deskripsi-singkat)
- [Algoritma BFS dan DFS](#algoritma-bfs-dan-dfs)
- [Fitur](#fitur)
- [Requirement](#requirement)
- [Instalasi dan Menjalankan Program](#instalasi-dan-menjalankan-program)
- [Cara Penggunaan](#cara-penggunaan)
- [Author](#author)

---

## Deskripsi Singkat

Aplikasi web untuk menelusuri elemen pada struktur **Document Object Model (DOM)** menggunakan algoritma **Breadth First Search (BFS)** dan **Depth First Search (DFS)** berdasarkan **CSS Selector** yang diberikan pengguna.

Aplikasi ini dapat:
- Melakukan scraping HTML dari URL yang diberikan, atau menerima input HTML langsung
- Melakukan parsing HTML menjadi DOM Tree
- Menelusuri DOM Tree menggunakan BFS atau DFS untuk mencari elemen yang sesuai dengan CSS Selector
- Memvisualisasikan struktur DOM Tree beserta animasi proses traversal
- Menampilkan traversal log, waktu pencarian, dan jumlah node yang dikunjungi
- Mendukung traversal paralel dengan multithreading
- Mencari Lowest Common Ancestor (LCA) dari dua elemen menggunakan Binary Lifting

---

## Algoritma BFS dan DFS

### Breadth First Search (BFS)

BFS menelusuri pohon DOM **per level** dari atas ke bawah. Dimulai dari root, semua node di level pertama dikunjungi sebelum lanjut ke level berikutnya. Implementasi menggunakan **queue (antrian)**.

```
Urutan kunjungan BFS:
        html          → kunjungi: html
       /    \
     head   body      → kunjungi: head, body
            /  \
          div  footer  → kunjungi: div, footer
          /
         p             → kunjungi: p
```

Karakteristik:
- Menemukan elemen di level paling atas terlebih dahulu
- Cocok untuk mencari **top-n** elemen yang paling dekat dengan root
- Kompleksitas waktu: **O(N)**, ruang: **O(N)** (ukuran queue maksimal selebar tree)

Versi paralel: setiap level diproses secara paralel menggunakan goroutine — semua node di level yang sama di-spawn ke goroutine berbeda, kemudian menunggu semua selesai sebelum lanjut ke level berikutnya.

---

### Depth First Search (DFS)

DFS menelusuri pohon DOM dengan **menyelam sedalam mungkin** ke satu cabang sebelum backtrack dan pindah ke cabang berikutnya. Implementasi menggunakan rekursi (implicit stack).

```
Urutan kunjungan DFS (pre-order):
        html          → kunjungi: html
       /    \
     head   body      → kunjungi: head, div, p, body, footer
            /  \
          div  footer
          /
         p
```

Karakteristik:
- Menemukan elemen di kedalaman tree lebih cepat
- Cocok untuk struktur HTML yang dalam dan nested
- Kompleksitas waktu: **O(N)**, ruang: **O(H)** (H = kedalaman tree)

Versi paralel: setiap subtree dari root di-assign ke goroutine berbeda menggunakan `errgroup`, dengan `context.WithCancel` untuk membatalkan goroutine lain saat hasil sudah cukup.

---

### CSS Selector Matching

Matching dilakukan secara **right-to-left** — dimulai dari bagian paling kanan selector, lalu naik ke atas pohon sesuai combinator. Selector yang didukung:

| Selector | Contoh | Keterangan |
|----------|--------|------------|
| Tag | `div` | Semua elemen `<div>` |
| Class | `.box` | Elemen dengan class `box` |
| ID | `#main` | Elemen dengan id `main` |
| Universal | `*` | Semua elemen |
| Child | `div > p` | `<p>` yang merupakan child langsung dari `<div>` |
| Descendant | `div p` | `<p>` yang merupakan keturunan `<div>` |
| Adjacent Sibling | `p + span` | `<span>` tepat setelah `<p>` |
| General Sibling | `p ~ span` | Semua `<span>` setelah `<p>` dalam parent yang sama |

---

### LCA Binary Lifting

Lowest Common Ancestor (LCA) mencari **node leluhur terdalam yang sama** dari dua node di pohon DOM. Implementasi menggunakan teknik **Binary Lifting**:

- Preprocessing `up[v][k]` = leluhur `2^k` langkah ke atas dari node `v`
- Preprocessing: **O(N log N)**
- Query per pasangan node: **O(log N)**


---

## Requirement

### Backend
- [Go](https://go.dev/dl/) versi **1.21** atau lebih baru
- Dependensi (diinstall otomatis via `go mod`):
  - `golang.org/x/net` — HTML parser
  - `golang.org/x/sync` — errgroup untuk parallel traversal

### Frontend
- [Node.js](https://nodejs.org/) versi **18** atau lebih baru
- [npm](https://www.npmjs.com/) versi **9** atau lebih baru

### (Opsional) Docker
- [Docker](https://www.docker.com/) versi **24** atau lebih baru
- [Docker Compose](https://docs.docker.com/compose/) versi **2** atau lebih baru

---

## Instalasi dan Menjalankan Program

### Cara 1: Manual

**Clone repository**
```bash
git clone https://github.com/[username]/Tubes2_[NamaKelompok].git
cd Tubes2_[NamaKelompok]
```

**Jalankan Backend**
```bash
cd backend
go mod tidy
go run src/main.go
```
Backend berjalan di `http://localhost:8080`

**Jalankan Frontend** (terminal baru)
```bash
cd frontend
npm install
npm run dev
```
Frontend berjalan di `http://localhost:3000`

---

### Cara 2: Docker 

```bash
git clone https://github.com/[username]/Tubes2_[NamaKelompok].git
cd Tubes2_[NamaKelompok]
docker compose up --build
```

Aplikasi berjalan di `http://localhost:3000`

Untuk menghentikan:
```bash
docker compose down
```

---

### Cara 3: Akses Online (Azure)

Aplikasi sudah di-deploy dan dapat diakses di:
```

```

---

## Cara Penggunaan


---

## Author

| Nama | NIM | GitHub |
|------|-----|--------|
| [Ahmad Zaky Robbani] | [13524045] | [@SlackingSloth](https://github.com/SlackingSlothh) |
| [Hakam Avicena Mustain] | [13524075] | [@hakamavicena](https://github.com/hakamavicena) |
| [Nicholas Luis Chandra] | [13524105] | [@staplesmaster](https://github.com/staplesmaster) |


---

# Lampiran

| No | Poin | Ya | Tidak |
|----|------|----|-------|
| 1 | Aplikasi berhasil di kompilasi tanpa kesalahan | [ ] | [ ] |
| 2 | Aplikasi berhasil dijalankan | [ ] | [ ] |
| 3 | Aplikasi dapat menerima input URL web, pilihan algoritma, CSS selector, dan jumlah hasil | [ ] | [ ] |
| 4 | Aplikasi dapat melakukan scraping terhadap web pada input | [ ] | [ ] |
| 5 | Aplikasi dapat menampilkan visualisasi pohon DOM | [ ] | [ ] |
| 6 | Aplikasi dapat menelusuri pohon DOM dan menampilkan hasil penelusuran | [ ] | [ ] |
| 7 | Aplikasi dapat menandai jalur tempuh oleh algoritma | [ ] | [ ] |
| 8 | Aplikasi dapat menyimpan jalur yang ditempuh algoritma dalam traversal log | [ ] | [ ] |
| 9 | [Bonus] Membuat video | [ ] | [ ] |
| 10 | [Bonus] Deploy aplikasi | [ ] | [ ] |
| 11 | [Bonus] Implementasi animasi pada penelusuran pohon | [ ] | [ ] |
| 12 | [Bonus] Implementasi multithreading | [ ] | [ ] |
| 13 | [Bonus] Implementasi LCA Binary Lifting | [ ] | [ ] |
