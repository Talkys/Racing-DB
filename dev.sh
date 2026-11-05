#!/usr/bin/env bash
# Local dev boot: postgres (docker) + go api + vite front end.
# Ctrl+C (or either service dying) tears everything down, including the db volume.
# I'm basically too lazy to keep two terminals running so this 1 hour task will save the 10 seconds of doing by hand.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

pgids=()
cleaning=0

cleanup() {
  ((cleaning)) && return
  cleaning=1
  trap - INT TERM
  echo
  echo "[dev] shutting down..."

  # Each service runs in its own process group, so signalling the group also
  # reaches the real children (the compiled go binary, the vite node process).
  for pgid in "${pgids[@]}"; do
    kill -TERM -- "-$pgid" 2>/dev/null || true
  done
  for _ in {1..50}; do
    alive=0
    for pgid in "${pgids[@]}"; do
      kill -0 -- "-$pgid" 2>/dev/null && alive=1
    done
    ((alive)) || break
    sleep 0.1
  done
  for pgid in "${pgids[@]}"; do
    kill -KILL -- "-$pgid" 2>/dev/null || true
  done

  docker compose down -v --remove-orphans
  echo "[dev] done"
}
trap cleanup EXIT
trap 'exit 130' INT TERM

# start <label> <dir> <command...>: runs the command in a new process group,
# prefixing every output line with the label.
start() {
  local label=$1 dir=$2
  shift 2
  setsid bash -c 'cd "$1" && shift 2 && "$@" 2>&1 | sed -u "s/^/[$0] /"' \
    "$label" "$dir" _ "$@" &
  pgids+=("$!")
}

echo "[dev] starting postgres..."
docker compose up -d --wait

start api "$ROOT/back-end" go run ./cmd/api
start web "$ROOT/front-end" bun run dev

echo "[dev] running, Ctrl+C to stop"
# Return as soon as either service exits, so a crash doesn't leave the other running.
wait -n || true
echo "[dev] a service exited"
