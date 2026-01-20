package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
)

// Supported media extensions
var videoExtensions = map[string]bool{
	".mkv": true, ".mp4": true, ".avi": true, ".wmv": true, ".m4v": true,
}

var ebookExtensions = map[string]bool{
	".epub": true, ".pdf": true, ".mobi": true, ".azw3": true,
	".cbz": true, ".cbr": true, ".djvu": true,
}

// isMediaFile checks if the extension is a supported media file
func isMediaFile(ext string) bool {
	return videoExtensions[ext] || ebookExtensions[ext]
}

// isVideoFile checks if the extension is a video file
func isVideoFile(ext string) bool {
	return videoExtensions[ext]
}

// isEbookFile checks if the extension is an ebook file
func isEbookFile(ext string) bool {
	return ebookExtensions[ext]
}

// FileInfo struct to hold file details
type FileInfo struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	IsDir       bool   `json:"isDir"`
	IsProcessed bool   `json:"isProcessed"`
	HasMedia    bool   `json:"hasMedia,omitempty"`
	MediaType   string `json:"mediaType,omitempty"` // "video" or "ebook"
}

// App struct
type App struct{}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// ListDirectory returns the contents of the given directory
func (a *App) ListDirectory(path string) ([]FileInfo, error) {
	if path == "" {
		return []FileInfo{}, nil
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := []FileInfo{}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())
		isProc := isProcessed(fullPath)

		if entry.IsDir() {
			// Show all directories to allow navigation
			// Check if directory contains media files (directly or in subdirs)
			hasMedia := dirContainsMedia(fullPath)
			files = append(files, FileInfo{
				Name:        entry.Name(),
				Size:        info.Size(),
				IsDir:       true,
				IsProcessed: isProc,
				HasMedia:    hasMedia,
			})
		} else {
			// Check file extension
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if isMediaFile(ext) {
				mediaType := "video"
				if isEbookFile(ext) {
					mediaType = "ebook"
				}
				files = append(files, FileInfo{
					Name:        entry.Name(),
					Size:        info.Size(),
					IsDir:       false,
					IsProcessed: isProc,
					MediaType:   mediaType,
				})
			}
		}
	}
	return files, nil
}

// dirContainsMedia checks if a directory contains media files (recursively, max 2 levels)
func dirContainsMedia(path string) bool {
	return dirContainsMediaDepth(path, 0, 2)
}

func dirContainsMediaDepth(path string, currentDepth, maxDepth int) bool {
	if currentDepth > maxDepth {
		return false
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if dirContainsMediaDepth(filepath.Join(path, entry.Name()), currentDepth+1, maxDepth) {
				return true
			}
		} else {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if isMediaFile(ext) {
				return true
			}
		}
	}
	return false
}

// GetMediaInfo executes mediainfo command on the file and returns output
func (a *App) GetMediaInfo(filePath string) (string, error) {
	// Check if mediainfo is in PATH
	path, err := exec.LookPath("mediainfo")
	if err != nil {
		return "", fmt.Errorf("mediainfo not found in PATH: %w", err)
	}

	cmd := exec.Command(path, filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// CreateTorrent creates a .torrent file for the given source path
func (a *App) CreateTorrent(sourcePath string, trackers []string, comment string, isPrivate bool) (string, error) {
	info := metainfo.Info{
		PieceLength: 256 * 1024,
	}

	if isPrivate {
		info.Private = new(bool)
		*info.Private = true
	}

	err := info.BuildFromFilePath(sourcePath)
	if err != nil {
		return "", err
	}

	mi := metainfo.MetaInfo{
		AnnounceList: func() [][]string {
			var list [][]string
			for _, url := range trackers {
				if strings.TrimSpace(url) != "" {
					list = append(list, []string{url})
				}
			}
			return list
		}(),
		Comment:   comment,
		CreatedBy: "AATM-API",
	}
	mi.SetDefaults()

	infoBytes, err := bencode.Marshal(info)
	if err != nil {
		return "", err
	}
	mi.InfoBytes = infoBytes

	// Determine output path
	// If source is in /host (read-only), save torrent to /torrents instead
	baseName := filepath.Base(sourcePath)
	var outputPath string
	if strings.HasPrefix(sourcePath, "/host") {
		outputPath = filepath.Join("/torrents", baseName+".torrent")
	} else {
		outputPath = sourcePath + ".torrent"
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	err = mi.Write(outFile)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

// SaveNfo saves the NFO content to a file derived from the source path
func (a *App) SaveNfo(sourcePath string, content string) (string, error) {
	// Get base name without extension
	baseName := filepath.Base(sourcePath)
	ext := filepath.Ext(sourcePath)
	lowerExt := strings.ToLower(ext)

	// Strip extension if it's a known media file
	if isMediaFile(lowerExt) {
		baseName = strings.TrimSuffix(baseName, ext)
	}

	// Determine output directory
	// If source is in /host (read-only), save to /torrents instead
	var outputDir string
	if strings.HasPrefix(sourcePath, "/host") {
		outputDir = "/torrents"
	} else {
		outputDir = filepath.Dir(sourcePath)
	}

	outputPath := filepath.Join(outputDir, baseName+".nfo")

	err := os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		return "", err
	}
	return outputPath, nil
}

// DeleteFile deletes the specified file
func (a *App) DeleteFile(path string) error {
	if path == "" {
		return nil
	}
	return os.Remove(path)
}

// GetDirectorySize calculates the total size of a directory recursively
func (a *App) GetDirectorySize(path string) (string, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return formatSize(size), nil
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
