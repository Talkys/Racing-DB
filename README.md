# Racing-DB
Racing DB is a database frontend service for racing data compatible with multiple categories. The idea is to
make a very refined codebase to showcase my back end architecture model and serve as a portfolio project as well.

## Running locally

Three processes, each in its own terminal. Front end and back end run separately for now; the Vite dev server
proxies `/api` to the back end.

### 1. Database (PostgreSQL 16)

```sh
docker compose up -d
```

`init.sql` (schema) and `seed.sql` (data) only run on a fresh container. To reset the data:

```sh
docker compose down && docker compose up -d
```

`seed.sql` is generated from `front-end/src/api/mock-data.ts`, so the database and the front end's mock mode hold
the same rows. After changing the mock data, run `bun run seed` in `front-end/` and reset the database.

### 2. Back end (Go, `:5000`)

```sh
cd back-end
go run ./cmd/api
```

| Variable | Default |
|---|---|
| `DATABASE_URL` | `postgres://admin:adminpassword@localhost:5432/racingdb?sslmode=disable` |
| `PORT` | `5000` |

- Tests: `go test ./...`
- Smoke tests against the running API and a fresh seed: `scripts/smoke.sh`
- Swagger docs: regenerate with `go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g cmd/api/main.go -o docs --parseInternal`

### 3. Front end (Vite + bun, `:5173`)

```sh
cd front-end
bun install
bun run dev
```

| Variable | Default | |
|---|---|---|
| `VITE_API_TARGET` | `http://localhost:5000` | Where the dev server proxies `/api` |
| `VITE_API_MOCK` | unset | `true` serves everything from `src/api/mock-data.ts`, no back end needed |

## API

All routes are `GET` under `/api/1`: `categories`, `circuits`, `drivers`, `races`, `teams`, `engines`, `chassis`,
`cars`, `articles`, `article-texts`, each as a list (`/races`) and by id (`/races/{id}`; categories use their slug).

Lists take `page` (25 rows per page; one extra row is returned when a next page exists),
`filter=column,operator,value` (`eq neq gt gte lt lte ilike in`; `in` takes `|`-separated values) and
`sort=column,direction`. Errors are JSON: `{"status": 400, "error": "..."}`. Draft and future-dated articles,
and their texts, are never served.

`/articles/search?q=&lang=&page=` searches article texts: whole words, case-insensitive, no stemming
(`"phrases"`, `or` and `-word` work). One result per article, using the text in `lang` when it matches, with a
Markdown-free preview whose matched words are marked. Without `q` it lists every article, newest first.
