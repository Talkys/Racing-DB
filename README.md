# Racing-DB
Racing DB is a database frontend service for racing data compatible with multiple categories. The idea is to
make a very refined codebase to showcase my back end architecture model and serve as a portfolio project as well.

Right now the project is in early development so a lot of things may change before I start branching hard.

## Requirements

| Tool | Minimum | Notes |
|---|---|---|
| Go | 1.27 | `back-end/go.mod` |
| Bun | 1.3 | Reads the text `bun.lock` (`configVersion`) |
| TypeScript | 7.0 | Installed by `bun install` |
| Vite | 8.3 | Installed by `bun install` |
| Docker | Compose v2 | Uses `docker compose`, not `docker-compose` |

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
| `LOG_FORMAT` | `text` (`json` for structured logs) |
| `LOG_LEVEL` | `info` (`debug`, `warn`, `error`) |

Logs go to stdout, one access-log line per request. Each response carries an `X-Request-Id` header (an incoming one
is reused when it is safe to log), and the same id is in the log line.

- Tests: `go test ./...`. The integration tests in `cmd/api` run the real handler against a fresh database loaded
  with `init.sql` + `seed.sql`, created on the server in `TEST_DATABASE_URL` (e.g. the compose one) or, when that is
  unset, in a throwaway [testcontainers](https://golang.testcontainers.org/) Postgres. Without either they are skipped.
- Lint: `golangci-lint run ./...` (config in `.golangci.yml`, includes the `gofmt`/`goimports` check)
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

Mock mode is for demos and UI work only. It follows the API (filters, sorts, paging, 400 on bad input), but its
article search is only an approximation of the back end's `websearch_to_tsquery` search, and the back end, covered
by the Go integration tests, is the reference.

- Tests: `bun test` (suites in `test/`)
- Lint: `bun run lint` ([Biome](https://biomejs.dev), lint and import order only; formatting is by hand)
- Type check, including the tests: `bun run typecheck`

## API

All routes are `GET` under `/api/1`: `categories`, `circuits`, `drivers`, `races`, `teams`, `engines`, `chassis`,
`cars`, `articles`, `article-texts`, each as a list (`/races`) and by id (`/races/{id}`; categories use their slug).

Lists take `page` (25 rows per page; one extra row is returned when a next page exists),
`filter=column,operator,value` (`eq neq gt gte lt lte ilike in`; `in` takes `|`-separated values; repeat it, up to
10 times, to AND filters: `?filter=category_id,eq,f1&filter=season,eq,2026`) and `sort`, a comma-separated list of
columns each optionally followed by `asc`/`desc` (`?sort=season,desc,round`, or repeated: `?sort=season,desc&sort=round`);
`id` breaks the remaining ties. Swagger lists each resource's filterable and sortable columns. Errors are JSON: `{"status": 400, "error": "..."}`. Draft and future-dated articles,
and their texts, are never served.

`races.payload` and `cars.parts` are JSONB written by the ETL, checked on every read against their contracts
(spec 3.3/3.4: required keys, types, enums, ranges; `races.payload` per category, any object for categories without a
contract). A row that breaks its contract is still served with status 200 and its JSON exactly as stored, flagged with
`payload_valid: false` / `parts_valid: false` and a `payload_errors` / `parts_errors` list (`"fastest_lap: missing"`),
and the back end logs a warning. Valid rows carry `true` and no error list. Clients must check the flag before reading
nested fields; the race and car pages show such rows as raw JSON. Only database errors give a 500.

`/articles/search?q=&lang=&page=` searches article texts: whole words, case-insensitive, no stemming
(`"phrases"`, `or` and `-word` work). One result per article, using the text in `lang` when it matches, with a
Markdown-free preview whose matched words are marked. Without `q` it lists every article, newest first.
