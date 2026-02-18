package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// AppSettings defines the structure of the settings to be saved
type AppSettings struct {
	RootPath         string   `json:"rootPath"`
	TorrentTrackers  string   `json:"torrentTrackers"`
	IsPrivateTorrent bool     `json:"isPrivateTorrent"`
	Passkey      string `json:"passkey"`
	LaCaleApiKey string `json:"laCaleApiKey"`
	// Torrent client selection: "qbittorrent", "transmission", "deluge", "none"
	TorrentClient string `json:"torrentClient"`
	// qBittorrent settings
	QbitUrl      string `json:"qbitUrl"`
	QbitUsername string `json:"qbitUsername"`
	QbitPassword string `json:"qbitPassword"`
	// Transmission settings
	TransmissionUrl      string `json:"transmissionUrl"`
	TransmissionUsername string `json:"transmissionUsername"`
	TransmissionPassword string `json:"transmissionPassword"`
	// Deluge settings
	DelugeUrl      string `json:"delugeUrl"`
	DelugePassword string `json:"delugePassword"`
	// Display settings
	ShowProcessed    bool `json:"showProcessed"`
	ShowNotProcessed bool `json:"showNotProcessed"`
	IsFullAuto       bool `json:"isFullAuto"`
	EnableHardlink   bool `json:"enableHardlink"`
	HardlinkDirs     []string `json:"hardlinkDirs"`
	// Output directory for .torrent and .nfo files
	OutputDir string `json:"outputDir"`
}

// InitDB initializes the SQLite database
func InitDB() {
	// Use /data directory in container for persistence
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "/data"
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("Warning: could not create data dir %s: %v", dataDir, err)
		// Fallback to local directory
		dataDir = "."
	}

	dbPath := filepath.Join(dataDir, "aatm.db")

	var errOpen error
	db, errOpen = sql.Open("sqlite", dbPath)
	if errOpen != nil {
		log.Fatal(errOpen)
	}

	createTables()
	log.Printf("Database initialized at %s", dbPath)
}

func createTables() {
	query := `
    CREATE TABLE IF NOT EXISTS settings (
        id INTEGER PRIMARY KEY CHECK (id = 1),
        data TEXT
    );
    CREATE TABLE IF NOT EXISTS processed_files (
        path TEXT PRIMARY KEY,
        processed_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

// SaveSettings saves the application settings to the database
func (a *App) SaveSettings(settings AppSettings) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT OR REPLACE INTO settings (id, data) VALUES (1, ?)", string(data))
	return err
}

// GetSettings retrieves the application settings from the database
func (a *App) GetSettings() AppSettings {
	var data string
	err := db.QueryRow("SELECT data FROM settings WHERE id = 1").Scan(&data)
	if err != nil {
		// Return default settings
		return getDefaultSettings()
	}
	var settings AppSettings
	json.Unmarshal([]byte(data), &settings)
	// Fill in defaults for empty values
	defaults := getDefaultSettings()
	if settings.RootPath == "" {
		settings.RootPath = defaults.RootPath
	}
	if settings.TorrentClient == "" {
		settings.TorrentClient = defaults.TorrentClient
	}
	if settings.QbitUrl == "" {
		settings.QbitUrl = defaults.QbitUrl
	}
	if settings.QbitUsername == "" {
		settings.QbitUsername = defaults.QbitUsername
	}
	if settings.QbitPassword == "" {
		settings.QbitPassword = defaults.QbitPassword
	}
	if settings.TransmissionUrl == "" {
		settings.TransmissionUrl = defaults.TransmissionUrl
	}
	if settings.DelugeUrl == "" {
		settings.DelugeUrl = defaults.DelugeUrl
	}
	if settings.DelugePassword == "" {
		settings.DelugePassword = defaults.DelugePassword
	}
	if settings.OutputDir == "" {
		settings.OutputDir = defaults.OutputDir
	}
	return settings
}

// getDefaultSettings returns the default application settings
func getDefaultSettings() AppSettings {
	return AppSettings{
		RootPath:             "/host",
		TorrentTrackers:      "",
		IsPrivateTorrent:     true,
		TorrentClient:        "qbittorrent",
		QbitUrl:              "http://localhost:8081",
		QbitUsername:         "admin",
		QbitPassword:         "adminadmin",
		TransmissionUrl:      "http://localhost:9091",
		TransmissionUsername: "",
		TransmissionPassword: "",
		DelugeUrl:            "http://localhost:8112",
		DelugePassword:       "deluge",
		ShowProcessed:        false,
		ShowNotProcessed:     true,
		OutputDir:            "/torrents",
	}
}

// MarkProcessed marks a file as processed in the database
func (a *App) MarkProcessed(path string) error {
	_, err := db.Exec("INSERT OR IGNORE INTO processed_files (path) VALUES (?)", path)
	return err
}

// ClearProcessedFiles removes all records from the processed_files table
func (a *App) ClearProcessedFiles() error {
	_, err := db.Exec("DELETE FROM processed_files")
	return err
}

// isProcessed checks if a file has explicitly been marked as processed
func isProcessed(path string) bool {
	if db == nil {
		return false
	}
	var exists int
	err := db.QueryRow("SELECT 1 FROM processed_files WHERE path = ?", path).Scan(&exists)
	return err == nil
}

// GetAllProcessedFiles returns all processed files from the database
func (a *App) GetAllProcessedFiles() ([]map[string]string, error) {
	rows, err := db.Query("SELECT path, processed_at FROM processed_files ORDER BY processed_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []map[string]string
	for rows.Next() {
		var path, processedAt string
		if err := rows.Scan(&path, &processedAt); err != nil {
			continue
		}
		files = append(files, map[string]string{
			"path":        path,
			"processedAt": processedAt,
		})
	}
	return files, nil
}
