# lifeos

Personal Life OS — **Go/Gin + SQLite + JWT, single-user lokal**.
`database/sql` tanpa ORM, migrasi goose. Reuse pola shiftbase, tanpa Docker/MySQL.
Dashboard web: [`lifeos-web`](https://github.com/RakhaYandra/lifeos-web) (Vite+React+TS, identitas Nexus).

## Quickstart 5 menit

```bash
git clone https://github.com/RakhaYandra/lifeos.git && cd lifeos
cp .env.example .env                     # ganti JWT_SECRET di prod
go install github.com/pressly/goose/v3/cmd/goose@v3.28.0
rm -f lifeos.db
goose -dir migrations sqlite lifeos.db up
sqlite3 lifeos.db < seed/seed.sql
go run ./cmd/api                        # :8080
curl -s localhost:8080/healthz          # {"status":"ok"}
```

Login seed: `aku@lifeos.local / Rahasia123` (single-user, register hanya akun pertama → 403 setelahnya).

Seed P6 terisi semua tabel (persona Jakarta, IDR, relasi nyambung goals→projects→tasks, trx→budgets, habit 30 hari): buka `npm run dev` di `lifeos-web` (`:5174`) dan masuk — dashboard langsung hidup.

## Contoh curl

```bash
B=localhost:8080
TOKEN=$(curl -s -X POST $B/v1/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"aku@lifeos.local","password":"Rahasia123"}' | cut -d'"' -f4)
A="Authorization: Bearer $TOKEN"
curl -s $B/v1/dashboard -H "$A"                                  # agregat home
curl -s "$B/v1/tasks/week?start=2026-09-08" -H "$A"              # seminggu
curl -s "$B/v1/budgets?year=2026&month=9" -H "$A"                # aktual otomatis
curl -s $B/v1/habits/streaks -H "$A"                             # streak
curl -s -X POST $B/v1/reviews -H "$A" -H 'Content-Type: application/json' \
  -d '{"week_start":"2026-09-14","wins":"coba"}'                 # auto-stat
```

## Endpoint (`/v1`, JWT kecuali auth)

auth register/login, `me`, settings, life-areas, projects, tasks (+today/week),
goals (?level), habits (+log/streaks/logs), transactions (+summary), budgets,
subscriptions (+upcoming), health-logs, workouts, learning, reading, reviews,
reminders (+upcoming), dashboard.

PUT = full replace (field tak dikirim = dihapus). Kontrak: `api/swagger.yaml`, `api/postman_collection.json`.

## Dev & CI

```bash
go build ./... && go vet ./... && gofmt -l .
go test ./... -cover            # service ≥ 60% (saat ini ~80%)
npx newman run api/postman_collection.json --env-var baseUrl=http://localhost:8080  # 41 cek
```

CI (`.github/workflows/ci.yml`): vet → lint → test+coverage gate → goose migrate sqlite → seed → boot → **Newman (41)** → swagger validate.

Struktur: `cmd/api/main.go`, `internal/{config,domain,repository,service,handler,middleware}`,
`migrations/` (goose sqlite, 9 file), `seed/seed.sql`, `api/{swagger.yaml,postman_collection.json}`.

Keputusan arsitektur: [ADR-001 SQLite tanpa ORM](docs/ADR-001-sqlite-no-orm.md).
Post-MVP: quarterly goals, milestones, calendar, travel, decision matrix, aset, CRM, sync Postgres — lihat plan di issue.
