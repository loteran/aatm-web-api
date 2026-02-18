package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// Initialize database
	InitDB()

	// Create app instance
	app := NewApp()

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Serve static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(staticFS))

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, staticFS, "index.html")
	})

	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", fileServer).ServeHTTP(w, r)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Directory operations
	r.Get("/api/files", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		if path == "" {
			http.Error(w, "path parameter required", http.StatusBadRequest)
			return
		}
		files, err := app.ListDirectory(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(files)
	})

	r.Get("/api/directory-size", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		if path == "" {
			http.Error(w, "path parameter required", http.StatusBadRequest)
			return
		}
		size, err := app.GetDirectorySize(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"size": size})
	})

	// MediaInfo
	r.Get("/api/mediainfo", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		if path == "" {
			http.Error(w, "path parameter required", http.StatusBadRequest)
			return
		}
		info, err := app.GetMediaInfo(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"mediainfo": info})
	})

	// Torrent creation (async)
	r.Post("/api/torrent/create", func(w http.ResponseWriter, r *http.Request) {
		var req CreateTorrentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		settings := app.GetSettings()
		taskID := app.StartCreateTorrent(req.SourcePath, req.Trackers, req.Comment, req.IsPrivate, req.TorrentName, settings.OutputDir, req.MediaType)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"taskId": taskID})
	})

	// Torrent creation status
	r.Get("/api/torrent/status/{taskId}", func(w http.ResponseWriter, r *http.Request) {
		taskID := chi.URLParam(r, "taskId")
		task := app.TaskMgr.GetTask(taskID)
		if task == nil {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	})

	// NFO operations
	r.Post("/api/nfo/save", func(w http.ResponseWriter, r *http.Request) {
		var req SaveNfoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		settings := app.GetSettings()
		nfoPath, err := app.SaveNfo(req.SourcePath, req.Content, req.TorrentName, settings.OutputDir, req.MediaType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"nfoPath": nfoPath})
	})

	// Steam API proxy (to avoid CORS issues)
	r.Get("/api/steam/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "query parameter required", http.StatusBadRequest)
			return
		}
		resp, err := http.Get("https://store.steampowered.com/api/storesearch/?term=" + query + "&l=french&cc=FR")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		io.Copy(w, resp.Body)
	})

	r.Get("/api/steam/details", func(w http.ResponseWriter, r *http.Request) {
		appid := r.URL.Query().Get("appid")
		if appid == "" {
			http.Error(w, "appid parameter required", http.StatusBadRequest)
			return
		}
		resp, err := http.Get("https://store.steampowered.com/api/appdetails?appids=" + appid + "&l=french")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		io.Copy(w, resp.Body)
	})

	// qBittorrent integration
	r.Post("/api/qbittorrent/upload", func(w http.ResponseWriter, r *http.Request) {
		var req QBittorrentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := app.UploadToQBittorrent(req.TorrentPath, req.QbitUrl, req.Username, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "uploaded"})
	})

	r.Post("/api/qbittorrent/remove", func(w http.ResponseWriter, r *http.Request) {
		var req QBittorrentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := app.RemoveFromQBittorrent(req.TorrentPath, req.QbitUrl, req.Username, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "removed"})
	})

	// Generic torrent client integration (uses settings to determine which client)
	r.Post("/api/torrent-client/upload", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			TorrentPath string `json:"torrentPath"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		settings := app.GetSettings()
		err := app.UploadToTorrentClient(req.TorrentPath, settings)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "uploaded", "client": settings.TorrentClient})
	})

	r.Post("/api/torrent-client/remove", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			TorrentPath string `json:"torrentPath"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		settings := app.GetSettings()
		err := app.RemoveFromTorrentClient(req.TorrentPath, settings)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "removed", "client": settings.TorrentClient})
	})

	// Hardlink creation
	r.Post("/api/hardlink/create", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			SourcePath   string   `json:"sourcePath"`
			HardlinkDirs []string `json:"hardlinkDirs"`
			DestName     string   `json:"destName"` // Optional: custom name for the hardlink
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Find the matching hardlink directory on the same device
		destDir, err := app.FindMatchingHardlinkDir(req.SourcePath, req.HardlinkDirs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create the hardlink (with optional custom name)
		hardlinkPath, err := app.CreateHardlink(req.SourcePath, destDir, req.DestName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":       "created",
			"hardlinkPath": hardlinkPath,
		})
	})

	// La Cale integration
	r.Post("/api/lacale/preview", func(w http.ResponseWriter, r *http.Request) {
		var req LaCalePreviewRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		preview, err := app.PreviewLaCale(req.MediaType, req.ReleaseInfo, req.ApiKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(preview)
	})

	r.Post("/api/lacale/upload", func(w http.ResponseWriter, r *http.Request) {
		var req LaCaleUploadRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := app.UploadToLaCale(
			req.TorrentPath,
			req.NfoPath,
			req.Title,
			req.Description,
			req.TmdbId,
			req.MediaType,
			req.ReleaseInfo,
			req.Passkey,
			req.ApiKey,
			req.Tags,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "uploaded"})
	})

	// Settings
	r.Get("/api/settings", func(w http.ResponseWriter, r *http.Request) {
		settings := app.GetSettings()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	})

	r.Post("/api/settings", func(w http.ResponseWriter, r *http.Request) {
		var settings AppSettings
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := app.SaveSettings(settings); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "saved"})
	})

	// Processed files
	r.Post("/api/processed/mark", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Path string `json:"path"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := app.MarkProcessed(req.Path); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "marked"})
	})

	r.Delete("/api/processed", func(w http.ResponseWriter, r *http.Request) {
		if err := app.ClearProcessedFiles(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "cleared"})
	})

	r.Get("/api/processed", func(w http.ResponseWriter, r *http.Request) {
		files, err := app.GetAllProcessedFiles()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(files)
	})

	// File operations
	r.Delete("/api/file", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		if path == "" {
			http.Error(w, "path parameter required", http.StatusBadRequest)
			return
		}
		if err := app.DeleteFile(path); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	})

	// Get port from env or default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("AATM API server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Request types
type CreateTorrentRequest struct {
	SourcePath  string   `json:"sourcePath"`
	Trackers    []string `json:"trackers"`
	Comment     string   `json:"comment"`
	IsPrivate   bool     `json:"isPrivate"`
	TorrentName string   `json:"torrentName"`
	MediaType   string   `json:"mediaType"`
}

type SaveNfoRequest struct {
	SourcePath  string `json:"sourcePath"`
	Content     string `json:"content"`
	TorrentName string `json:"torrentName"`
	MediaType   string `json:"mediaType"`
}

type QBittorrentRequest struct {
	TorrentPath string `json:"torrentPath"`
	QbitUrl     string `json:"qbitUrl"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type LaCaleUploadRequest struct {
	TorrentPath string      `json:"torrentPath"`
	NfoPath     string      `json:"nfoPath"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	TmdbId      string      `json:"tmdbId"`
	MediaType   string      `json:"mediaType"`
	ReleaseInfo ReleaseInfo `json:"releaseInfo"`
	Passkey     string      `json:"passkey"`
	ApiKey      string      `json:"apiKey"`
	Tags        []string    `json:"tags"`
}

type LaCalePreviewRequest struct {
	MediaType   string      `json:"mediaType"`
	ReleaseInfo ReleaseInfo `json:"releaseInfo"`
	ApiKey      string      `json:"apiKey"`
}
