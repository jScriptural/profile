# profile

This project is a Go-based API that generates rich user profiles by aggregating data from three external sources (Genderize, Agify, and Nationalize).
It utilizes a time-ordered UUIDv7 system for data persistence and is optimized for high-performance indexing in SQLite.


## Project structure

profile/
├── main.go             # Server initialization & routing
├── go.mod              # Dependency management
├── database.db         # Local SQLite file (not tracked in git)
├── internal/
│   ├── api/            # HTTP transport layer
│   │   ├── handlers.go # Request/Response processing
│   │   └── routes.go   # Endpoint definitions
│   ├── models/         # Data structures (The "source of truth")
│   │   └── models.go  # The Profile struct + UUIDv7 logic
│   └── service/        # The "Brain" (External API calls & logic)
│       └── service.go # Functions to call Genderize, Agify, etc.
└── store/              # Persistence layer
    └── store.go       # SQL queries



## API Endpoints Summary

| Method | Endpoint | Description | Status Codes |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/profiles` | Create or retrieve profile via JSON body. | `201`, `200`, `400`, `422` |
| `GET` | `/api/profiles/{name}` | Get a single profile by name. | `200`, `404` |
| `GET` | `/api/profiles` | List all profiles (supports query params). | `200`, `500` |
| `DELETE` | `/api/profiles/{name}` | Remove a profile from the database. | `204`, `404`, `500` |



##Local setup

Clone the repo: `git clone https://github.com/jscriptural/profile`

Install dependencies: `go mod tidy`

Run the server: `go run main.go`


```markdown
### Environment Variables
- `PORT`: The port the server listens on (Default: 8080).
- `DATABASE_URL`: Path to your SQLite database file(default: "./profile.db").

```
The server will listen on the port specified in the PORT environment variable (defaults to 8080)



##Live demo

**Base url: ** `https://profiler.pxxl.click`


