# Simple Go API

A minimal Go REST API with the following features:
- Counter endpoint that returns sequential numbers
- Fake LLM endpoint with Server-Sent Events streaming
- Swagger documentation
- Structured logging with slog
- Environment-based configuration

## Setup

Copy `.env.example` to `.env` and adjust values as needed:

```bash
cp .env.example .env
```

## Installation

```bash
go mod download
```

## Running

```bash
go run .
```

The API will start on the configured port (default: 8080).

## API Endpoints

### Counter API
- **POST** `/api/counter/next` - Get next sequential number

Example:
```bash
curl -X POST http://localhost:8080/api/counter/next
```

Response:
```json
{"value": 1}
```

### Fake LLM API
- **POST** `/api/llm/chat` - Stream tokens via Server-Sent Events

Example:
```bash
curl -X POST http://localhost:8080/api/llm/chat \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Hello, AI!"}'
```

Response (SSE stream):
```
data: This
data:  is
data:  a
...
```

## Swagger Documentation

Visit `http://localhost:8080/swagger/` to view the API documentation.

## Configuration

All configuration is done via environment variables (see `.env.example`):

- `PORT` - Server port (default: 8080)
- `LOG_LEVEL` - Logging level: debug, info (default: info)
- `DEBUG` - Enable debug mode (default: false)

## Logging

The application uses `log/slog` for structured logging. All log messages use attributes instead of string interpolation, ensuring structured, queryable logs.

Example log output:
```json
{"time":"2024-01-15T10:30:45.123Z","level":"INFO","msg":"counter incremented","value":1}
{"time":"2024-01-15T10:30:46.456Z","level":"INFO","msg":"llm request received","prompt_length":"13"}
```
