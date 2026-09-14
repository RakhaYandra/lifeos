# ADR-001: Tanpa ORM, SQLite pure-Go

Status: diterima.

Konteks: LifeOS pribadi-lokal, 1 file DB, tanpa Docker.
Pola reuse shiftbase (database/sql eksplisit + goose).

Keputusan: `database/sql` + `modernc.org/sqlite` (pure-Go, tanpa CGO),
SQL tangan di `internal/repository`. Migrasi eksplisit via goose dialect sqlite.

Alasan:
- Jalan lokal langsung (`./lifeos`, `lifeos.db`), cocok 1 user.
- Ganti ke Postgres/MySQL nanti hanya ganti driver + DSN (service terisolasi interface).
- Nol magic, mudah direview.

Konsekuensi: mapping manual via Scan; FK via PRAGMA.
