# Counting App

A simple Go application that serves a live-updating Vue/Bootstrap dashboard showing a running table of entries.

## Features

- **Real-time updates**: Automatically adds new entries every 20 seconds
- **Live dashboard**: Auto-refreshes every 20 seconds using Vue.js
- **Responsive design**: Built with Bootstrap 5 for mobile compatibility
- **Data retention**: Keeps the last 90 entries, displays the latest 50
- **Docker-ready**: Binds to 0.0.0.0:8080 for containerization

## Running Locally

### Prerequisites
- Go 1.21 or later

### Build and Run
```bash
go run main.go
```

The application will start on `http://localhost:8080`

## Running with Docker

### Build the Docker image
```bash
docker build -t counting:latest .
```

### Run the Docker container
```bash
docker run -p 8080:8080 counting:latest
```

Access the application at `http://localhost:8080`

## API Endpoints

- `GET /` - Serves the HTML dashboard
- `GET /api/data` - Returns the last 50 entries as JSON

## How It Works

1. **Backend (Go)**:
   - Maintains a slice of the last 90 entries
   - Generates a new entry every 20 seconds
   - Exposes entries via JSON API

2. **Frontend (Vue + Bootstrap)**:
   - Fetches data every 20 seconds via `/api/data`
   - Displays entries in a responsive table
   - Shows live status and entry count

## Entry Structure

```json
{
  "id": 1,
  "timestamp": "2026-01-08T10:30:45Z",
  "message": "Entry #1 - 2026-01-08 10:30:45"
}
```
