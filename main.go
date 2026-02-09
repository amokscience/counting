package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Entry struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

var (
	entries []Entry
	mu      sync.RWMutex
	counter int
)

func init() {
	entries = make([]Entry, 0)
	counter = 0

	// Start the background goroutine that adds entries every 20 seconds
	go addEntriesPeriodically()
}

func addEntriesPeriodically() {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	// Add initial entry
	addNewEntry()

	for range ticker.C {
		addNewEntry()
	}
}

func addNewEntry() {
	mu.Lock()
	defer mu.Unlock()

	counter++
	entry := Entry{
		ID:        counter,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("Entry #%d - %s", counter, time.Now().Format("2006-01-02 15:04:05")),
	}

	// Add to the beginning of the slice
	entries = append([]Entry{entry}, entries...)

	// Keep only the last 90 entries
	if len(entries) > 100 {
		entries = entries[:100]
	}
}

func handleData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mu.RLock()
	defer mu.RUnlock()

	// Return the first 25 entries (most recent)
	displayEntries := entries
	if len(entries) > 25 {
		displayEntries = entries[:25]
	}

	json.NewEncoder(w).Encode(displayEntries)
}

func handleCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mu.RLock()
	currentCounter := counter
	mu.RUnlock()

	json.NewEncoder(w).Encode(map[string]int{"counter": currentCounter})
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, htmlContent)
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/api/data", handleData)
	http.HandleFunc("/api/counter", handleCounter)

	addr := "0.0.0.0:8080"
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

const htmlContent = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Counting - Running Table</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://unpkg.com/vue@3/dist/vue.global.js"></script>
    <style>
        body {
            background-color: #f8f9fa;
            padding: 20px;
        }
        .container {
            max-width: 900px;
            margin: 0 auto;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px 0;
            border-radius: 8px;
            margin-bottom: 30px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 2.5rem;
        }
        .refresh-info {
            font-size: 0.9rem;
            margin-top: 10px;
            opacity: 0.9;
        }
        .table-container {
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .table {
            margin-bottom: 0;
        }
        .table thead {
            background-color: #667eea;
            color: white;
        }
        .table tbody tr {
            border-bottom: 1px solid #dee2e6;
        }
        .table tbody tr:hover {
            background-color: #f0f0f0;
        }
        .table td {
            vertical-align: middle;
            padding: 12px;
        }
        .id-badge {
            background-color: #667eea;
            color: white;
            padding: 4px 8px;
            border-radius: 4px;
            font-weight: bold;
            font-size: 0.9rem;
        }
        .timestamp {
            color: #6c757d;
            font-size: 0.9rem;
        }
        .loading {
            text-align: center;
            padding: 30px;
            color: #6c757d;
        }
        .spinner {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 3px solid rgba(102, 126, 234, 0.1);
            border-radius: 50%;
            border-top-color: #667eea;
            animation: spin 0.8s linear infinite;
            margin-right: 10px;
        }
        @keyframes spin {
            to { transform: rotate(360deg); }
        }
        .empty-state {
            text-align: center;
            padding: 50px 20px;
            color: #6c757d;
        }
        .status-badge {
            display: inline-block;
            padding: 4px 12px;
            background-color: #28a745;
            color: white;
            border-radius: 20px;
            font-size: 0.85rem;
            animation: pulse 2s infinite;
        }
        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.7; }
        }
    </style>
</head>
<body>
    <div id="app">
        <div class="header">
            <h1> 📊 Running Counter </h1>
            <div class="refresh-info">
                <span class="status-badge">● Live</span>
                <p style="margin: 10px 0 0 0;">Auto-refreshes every 20 seconds • Showing last 25 of {{ totalEntries }} entries</p>
            </div>
        </div>

        <div class="table-container">
            <div v-if="loading" class="loading">
                <div class="spinner"></div>
                Loading data...
            </div>

            <div v-else-if="entries.length === 0" class="empty-state">
                <p>No entries yet. The first entry will appear shortly...</p>
            </div>

            <div v-else>
                <table class="table">
                    <thead>
                        <tr>
                            <th style="width: 10%;">ID</th>
                            <th style="width: 35%;">Timestamp</th>
                            <th style="width: 55%;">Message</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="entry in entries" :key="entry.id">
                            <td>
                                <span class="id-badge">{{ entry.id }}</span>
                            </td>
                            <td>
                                <span class="timestamp">{{ formatTime(entry.timestamp) }}</span>
                            </td>
                            <td>{{ entry.message }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>

    <script>
        const { createApp } = Vue;

        createApp({
            data() {
                return {
                    entries: [],
                    loading: true,
                    totalEntries: 0
                };
            },
            methods: {
                async fetchData() {
                    try {
                        const response = await fetch('/api/data');
                        if (!response.ok) throw new Error('Network response was not ok');
                        
                        const data = await response.json();
                        this.entries = data || [];
                        
                        // Calculate total entries (this is an estimate based on the ID)
                        if (this.entries.length > 0) {
                            this.totalEntries = this.entries[0].id;
                        }
                        
                        this.loading = false;
                    } catch (error) {
                        console.error('Error fetching data:', error);
                        this.loading = false;
                    }
                },
                formatTime(timestamp) {
                    const date = new Date(timestamp);
                    return date.toLocaleString('en-US', {
                        year: 'numeric',
                        month: '2-digit',
                        day: '2-digit',
                        hour: '2-digit',
                        minute: '2-digit',
                        second: '2-digit'
                    });
                }
            },
            mounted() {
                this.fetchData();
                // Check counter and refresh every 20 seconds if needed
                setInterval(async () => {
                    try {
                        const response = await fetch('/api/counter');
                        if (response.ok) {
                            const data = await response.json();
                            // Check if counter is a multiple of 20 (and greater than 0)
                            if (data.counter > 0 && data.counter % 20 === 0) {
                                // Refresh the page
                                location.reload();
                                return;
                            }
                        }
                    } catch (error) {
                        console.error('Error fetching counter:', error);
                    }
                    // Fetch data if page wasn't refreshed
                    this.fetchData();
                }, 20000);
            }
        }).mount('#app');
    </script>
</body>
</html>
`
