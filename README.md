# Exile Tracker

A Go backend service for tracking Path of Exile characters, snapshots.  
It provides a REST API, background data fetcher, and integrates with the Path of Exile public API.

---

## Why

Well... I am really bad building and understanding a lot of stuff in Path of Exile and 
my biggest issue to really enjoy the game is understanding build upgrades, when and which upgrades do i need to do.
So a wanted to track good players to see which updates they do.

---

## Plans
- **POB Export** The application will generate a snapshot automatically, this string will importable into PoB
- **SSH Application** There will be an option to connect to a server via ssh to see in there directly the snapshots 


## Features

- **REST API** for managing accounts, characters, and passive skill snapshots
- **Background fetcher** for periodically updating character data from the PoE API
- **SQLite database** with migration support (Goose)
- **Structured logging** with zerolog
- **Graceful shutdown** for both API and background services

---

## Current focus

- **Understanding PoB** Im trying to figure out how to generate something that PoB can understand

---

## Project Structure

```
cmd/
  main.go         # Application entrypoint
  api/            # API server setup
config/           # Configuration loading
db/               # Database and migrations
models/           # Data models (internal and API)
poeclient/        # Path of Exile API client
repository/       # Database access layer
services/         # Business logic and background fetcher
utils/            # Logging and helpers
```

---

## Getting Started

### Prerequisites

- Go 1.20+
- [Goose](https://github.com/pressly/goose) for migrations

### Setup

1. **Clone the repository**
   ```sh
   git clone https://github.com/yourusername/exile-tracker.git
   cd exile-tracker
   ```

2. **Configure environment variables**

   Create a `.env` file or set the following variables:
   ```
   DB_PATH=./data.db
   PORT=:3000
   ```

3. **Run database migrations**
   ```sh
   goose -dir ./migrations sqlite3 ./data.db up
   ```

4. **Build and run the server**
   ```sh
   go build -o exile-tracker ./cmd
   ./exile-tracker
   ```

---

## API Endpoints

All endpoints are prefixed with `/api/v1`.

### POB Snapshots

- `GET    /pobsnapshots/character/{characterId}`        — List snapshots for a character
- `GET    /pobsnapshots/character/{characterId}/latest` — Get latest snapshot for a character
- `GET    /pobsnapshots/{id}`                           — Get snapshot by ID

---

## Development

- Logging is handled by [zerolog](https://github.com/rs/zerolog).
- API routing uses [chi](https://github.com/go-chi/chi).
- Database access uses the standard `database/sql` package.
- Background fetcher runs on a configurable interval.

---

## License

MIT

---

## Credits

- [Path of Exile](https://www.pathofexile.com/)
- [Goose](https://github.com/pressly/goose)
- [Path of Building](https://github.com/PathOfBuildingCommunity/PathOfBuilding)
