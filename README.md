# Gator - Blog Aggregator CLI

Gator is a command-line RSS feed aggregator written in Go. It allows you to subscribe to RSS feeds, automatically fetch new posts, and browse content from all your subscriptions in one place.

## Prerequisites

Before installing Gator, ensure you have the following installed:

- **Go** (version 1.23 or higher) - [Download Go](https://golang.org/dl/)
- **PostgreSQL** - [Download PostgreSQL](https://www.postgresql.org/download/)

## Installation

### Install via Go

```bash
go install github.com/Legendary-Coder-GT/blog_aggregator@latest
```

This will install the `blog_aggregator` binary to your `$GOPATH/bin` directory. Make sure this directory is in your `PATH`.

> **Note:** Go programs are statically compiled binaries. After installation, you can run the program without needing the Go toolchain installed.

### Build from Source

```bash
git clone https://github.com/Legendary-Coder-GT/blog_aggregator.git
cd blog_aggregator
go build -o gator
```

## Database Setup

1. **Create a PostgreSQL database:**

```bash
createdb blog_aggregator
```

2. **Run the database migrations:**

Using [Goose](https://github.com/pressly/goose):

```bash
goose -dir sql/schema postgres "postgres://username:password@localhost:5432/blog_aggregator?sslmode=disable" up
```

Or manually apply the SQL files in `sql/schema/` in order (001 through 005).

## Configuration

Gator requires a configuration file at `~/.gatorconfig.json` with the following format:

```json
{
  "db_url": "postgres://username:password@localhost:5432/blog_aggregator?sslmode=disable"
}
```

Replace `username` and `password` with your PostgreSQL credentials.

## Usage

### Getting Started

1. **Register a new user:**

```bash
gator register <username>
```

2. **Add an RSS feed:**

```bash
gator addfeed "Blog Name" "https://example.com/feed.xml"
```

3. **Start the aggregator to fetch posts:**

```bash
gator agg 1m
```

This will fetch new posts from all feeds every 1 minute. Press `Ctrl+C` to stop.

4. **Browse your posts:**

```bash
gator browse
```

### Available Commands

| Command | Description |
|---------|-------------|
| `register <name>` | Create a new user and set as current user |
| `login <name>` | Switch to an existing user |
| `users` | List all registered users |
| `addfeed <name> <url>` | Add a new RSS feed and automatically follow it |
| `feeds` | List all feeds in the database |
| `follow <url>` | Follow an existing feed |
| `unfollow <url>` | Unfollow a feed |
| `following` | List feeds you're following |
| `agg <duration>` | Start the aggregator (e.g., `30s`, `1m`, `5m`) |
| `browse [limit]` | View posts from followed feeds (default: 2 posts) |
| `reset` | Delete all users and data |

### Examples

**Register and add feeds:**

```bash
gator register alice
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
gator addfeed "Go Blog" "https://blog.golang.org/feed.atom"
```

**Follow a feed added by another user:**

```bash
gator follow "https://example.com/feed.xml"
```

**View more posts:**

```bash
gator browse 10
```

**Switch users:**

```bash
gator login bob
```

## Project Structure

```
.
├── main.go              # Entry point and command registration
├── commands.go          # Command registry implementation
├── handlers.go          # Command handlers
├── middleware.go        # Authentication middleware
├── rss_fetch.go         # RSS feed fetching and parsing
├── scrapeFeeds.go       # Feed aggregation logic
├── internal/
│   ├── config/          # Configuration file handling
│   └── database/        # Generated database code (sqlc)
└── sql/
    ├── schema/          # Database migrations
    └── queries/         # SQL query definitions
```

## Development

For development, you can use:

```bash
go run . <command>
```

To regenerate the database code after modifying SQL queries:

```bash
sqlc generate
```

## License

This project is open source.
