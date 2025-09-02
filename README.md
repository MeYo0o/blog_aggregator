## Gator — Blog Aggregator CLI (Go)

Gator is a command-line blog/RSS aggregator written in Go. It connects to a PostgreSQL database, follows RSS feeds, periodically fetches new posts, and lets you browse posts from feeds you follow.

Built while following the Boot.dev course “Build a Blog Aggregator in Go” [course link]. See the course overview here: [Build a Blog Aggregator in Go](https://www.boot.dev/courses/build-blog-aggregator-golang).

### Features

- **Configurable CLI**: Configure DB URL and current user in `~/.gatorconfig.json`.
- **User accounts**: Register and switch the current user.
- **Feeds**: Create feeds and follow/unfollow feeds.
- **Aggregator**: Long-running worker that continuously fetches feeds on a schedule.
- **Posts**: Persisted posts with per-user browsing and configurable limits.

## Requirements

- **Go**: 1.21+ recommended
- **PostgreSQL**: 14+ recommended
- **psql**: PostgreSQL command-line client

Note: Go programs are statically compiled binaries. After `go build` or `go install`, you can run `gator` without the Go toolchain on the machine. `go run .` is for development; `gator` (the installed binary) is for production.

## Install Go and PostgreSQL

### Go

1. Download from `https://go.dev/dl/` for your OS.
2. Install and ensure your PATH includes Go’s bin directory.
3. Verify:

```bash
go version
```

### PostgreSQL and psql

1. macOS (Homebrew):

```bash
brew install postgresql@15
brew services start postgresql@15
```

2. Linux (Debian/Ubuntu):

```bash
sudo apt update && sudo apt install -y postgresql postgresql-contrib
sudo systemctl enable --now postgresql
```

3. Windows: Use the official installer from `https://www.postgresql.org/download/`.

Verify `psql`:

```bash
psql --version
```

## Database Setup

1. Create a database, for example `gator`:

```bash
createdb gator
```

2. Set your database URL (example for local Postgres with no password):

```
postgres://localhost:5432/gator?sslmode=disable
```

3. Run migrations (uses Goose). From the project root:

```bash
make migrate
# or
goose up
```

## Installation

Install the CLI into your GOPATH/bin (or GOBIN):

```bash
go install ./...
```

This builds and installs the `gator` binary to `$(go env GOPATH)/bin` (or `GOBIN` if set). Ensure that directory is in your PATH.

Alternatively, for local development:

```bash
go run . <command> [args]
```

To produce a standalone binary in the project directory:

```bash
go build -o gator
./gator <command> [args]
```

## Configuration

Gator reads config from `~/.gatorconfig.json` with the following keys:

```json
{
  "db_url": "postgres://localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Create this file manually, or it will be expected to already exist before running.

## Usage Overview

Run commands either via the dev workflow:

```bash
go run . <command> [args]
```

or via the installed/built binary:

```bash
gator <command> [args]
```

### Commands

The following commands are registered in `main.go` and implemented under `internal/commands`:

- `login <username>`: Set the current user in the config. Fails if the user doesn’t exist in DB.
- `register <username>`: Create a new user and set it as the current user.
- `reset`: Delete all users (and cascade-owned data via FKs where applicable).
- `users`: List all users, marking the current user.
- `agg <duration>`: Start the aggregator worker. Example durations: `10s`, `5m`, `1h`. The worker:
  - Picks the next feed to fetch by `last_fetched_at` (NULLS FIRST, oldest first)
  - Fetches RSS with timeouts and response validation
  - Parses items and stores posts (ignores duplicates via DB uniqueness)
  - Marks the feed as fetched after a successful parse
- `feeds`: List all feeds and their owners.
- `addfeed <name> <url>`: Create a feed owned by the current user and automatically follow it.
- `follow <url>`: Follow an existing feed by URL for the current user.
- `following`: List feeds followed by the current user.
- `unfollow <url>`: Unfollow a feed by URL for the current user.
- `browse [limit]`: Show recent posts from feeds the current user follows. Defaults to a small limit; pass an integer to override.

### Examples

```bash
# set config file first (~/.gatorconfig.json) with your db_url

# register and login
gator register alice
gator login alice

# add and follow a feed
gator addfeed "TechCrunch" https://techcrunch.com/feed/

# start aggregator to pull feeds every 30 seconds
gator agg 30s

# browse latest 10 posts for current user
gator browse 10
```

## Development Notes

- Uses: Go, PostgreSQL, `sqlc`, and `goose`.
- SQL is defined in `sql/schema` (migrations) and `sql/queries` (sqlc queries).
- Generated Go code lives under `internal/database`.
- The aggregator uses a `time.Ticker` and safe HTTP fetching with timeouts, content-type checks, and limited body size.
- Posts are deduped by a DB uniqueness constraint (typically on URL or `(feed_id, url)` depending on the schema). Duplicate insert errors are safely ignored in code.

## Project Capabilities and Extensions

- Current:
  - Multiple users, follow/unfollow feeds
  - Periodic fetching with fair scheduling by `last_fetched_at`
  - Robust RSS parsing with multiple date formats (RFC1123Z/RFC1123/RFC822Z/RFC822/RFC3339)
  - Store and browse persisted posts with configurable limits
- Possible Expansions:
  - Full-text search over posts
  - Web UI or TUI for browsing
  - Background job concurrency across multiple workers
  - Rate limiting and per-feed backoff on errors
  - Enriched post parsing (images, content extraction)
  - Notifications (email/webhook) for new posts

## What I Learned (from Boot.dev Course)

From the Boot.dev course: [Build a Blog Aggregator in Go](https://www.boot.dev/courses/build-blog-aggregator-golang)

- **Config**: Built a configuration system to get/set values used by the CLI
- **Database**: Set up PostgreSQL, Goose (migrations), and SQLc (type-safe queries)
- **RSS**: Implemented HTTP fetching with timeouts and robust XML parsing
- **Following**: Enabled users to follow/unfollow feeds
- **Aggregate**: Built a long-running worker to continuously aggregate posts

## Troubleshooting

- Ensure `~/.gatorconfig.json` exists and `db_url` is correct
- Verify Postgres is running and reachable from your machine
- Run migrations before first use: `make migrate` or `goose up`
- Make sure `$(go env GOPATH)/bin` (or `GOBIN`) is in your PATH to use `gator`
- For development-only runs, use `go run . ...` from the project root
