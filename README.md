# lifeos

Personal Life OS — **Go/Gin + SQLite + JWT, single-user lokal**.
`database/sql` tanpa ORM, migrasi goose. Reuse pola shiftbase, tanpa Docker/MySQL.

## Quickstart 5 menit (P0)

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

Login seed: `aku@lifeos.local / Rahasia123` (single-user, register hanya akun pertama).

## Contoh curl P0

```bash
B=localhost:8080
TOKEN=$(curl -s -X POST $B/v1/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"aku@lifeos.local","password":"Rahasia123"}' | cut -d'"' -f4)
A="Authorization: Bearer $TOKEN"
curl -s $B/v1/me -H "$A"
curl -s $B/v1/settings -H "$A"
curl -s $B/v1/life-areas -H "$A"
```

## Dev & CI

```bash
go build ./... && go vet ./... && gofmt -l .
go test ./... -cover
```

CI: vet → lint → test+coverage → goose migrate sqlite → seed → boot → Newman → swagger validate.
Struktur: `cmd/api/main.go`, `internal/{config,domain,repository,service,handler,middleware}`,
`migrations/` (goose sqlite), `seed/seed.sql`, `api/{swagger.yaml,postman_collection.json}`.

Roadmap: P1 tasks/projects → P2 goals/habits → P3 finance → P4 health/learning/reviews → P5 web → P6 seed penuh+porto.
