package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

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

var gameExtensions = map[string]bool{
	".iso": true, ".nsp": true, ".xci": true, ".pkg": true,
	".zip": true, ".rar": true, ".7z": true,
}

// isMediaFile checks if the extension is a supported media file
func isMediaFile(ext string) bool {
	return videoExtensions[ext] || ebookExtensions[ext] || gameExtensions[ext]
}

// isVideoFile checks if the extension is a video file
func isVideoFile(ext string) bool {
	return videoExtensions[ext]
}

// isEbookFile checks if the extension is an ebook file
func isEbookFile(ext string) bool {
	return ebookExtensions[ext]
}

// isGameFile checks if the extension is a game file
func isGameFile(ext string) bool {
	return gameExtensions[ext]
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
				} else if isGameFile(ext) {
					mediaType = "game"
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

	// Replace full path with just filename in "Complete name" line
	result := string(output)
	fileName := filepath.Base(filePath)
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "Complete name") {
			lines[i] = "Complete name                            : " + fileName
			break
		}
	}
	return strings.Join(lines, "\n"), nil
}

// CreateTorrent creates a .torrent file for the given source path
// torrentName is the name that will appear in the torrent (the release name)
func (a *App) CreateTorrent(sourcePath string, trackers []string, comment string, isPrivate bool, torrentName string) (string, error) {
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

	// Utiliser le nom personnalisé si fourni, sinon garder le nom du fichier source
	if torrentName != "" {
		info.Name = torrentName
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
	// Use the torrent name for the output file if provided
	var baseName string
	if torrentName != "" {
		baseName = torrentName
	} else {
		baseName = filepath.Base(sourcePath)
	}
	var outputPath string
	if strings.HasPrefix(sourcePath, "/host") {
		outputPath = filepath.Join("/torrents", baseName+".torrent")
	} else {
		outputPath = filepath.Join(filepath.Dir(sourcePath), baseName+".torrent")
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

// SaveNfo saves the NFO content to a file
// If torrentName is provided, it will be used as the filename
func (a *App) SaveNfo(sourcePath string, content string, torrentName string) (string, error) {
	// Determine base name: use torrentName if provided, otherwise derive from source
	var baseName string
	if torrentName != "" {
		baseName = torrentName
	} else {
		baseName = filepath.Base(sourcePath)
		ext := filepath.Ext(sourcePath)
		lowerExt := strings.ToLower(ext)

		// Strip extension if it's a known media file
		if isMediaFile(lowerExt) {
			baseName = strings.TrimSuffix(baseName, ext)
		}
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

// getDeviceID returns the device ID of the filesystem containing the path
func getDeviceID(path string) (uint64, error) {
	var stat syscall.Stat_t
	err := syscall.Stat(path, &stat)
	if err != nil {
		return 0, err
	}
	return stat.Dev, nil
}

// FindMatchingHardlinkDir finds a hardlink directory on the same device as sourcePath
func (a *App) FindMatchingHardlinkDir(sourcePath string, hardlinkDirs []string) (string, error) {
	sourceDevID, err := getDeviceID(sourcePath)
	if err != nil {
		return "", fmt.Errorf("cannot get device ID for source: %w", err)
	}

	for _, dir := range hardlinkDirs {
		if dir == "" {
			continue
		}
		// Ensure the directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		dirDevID, err := getDeviceID(dir)
		if err != nil {
			continue
		}
		if dirDevID == sourceDevID {
			return dir, nil
		}
	}

	return "", fmt.Errorf("no hardlink directory found on the same device as %s", sourcePath)
}

// CreateHardlink creates hardlinks for the source path in the destination directory
// destName is optional: if provided, the hardlink will use this name instead of the original
func (a *App) CreateHardlink(sourcePath string, destDir string, destName string) (string, error) {
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return "", fmt.Errorf("cannot stat source: %w", err)
	}

	// Use custom destName if provided, otherwise use original basename
	baseName := filepath.Base(sourcePath)
	if destName != "" {
		baseName = destName
	}
	destPath := filepath.Join(destDir, baseName)

	if sourceInfo.IsDir() {
		// For directories, create directory structure and hardlink all files
		err = a.hardlinkDirectory(sourcePath, destPath)
		if err != nil {
			return "", err
		}
	} else {
		// For single files, just create the hardlink
		err = os.Link(sourcePath, destPath)
		if err != nil {
			return "", fmt.Errorf("failed to create hardlink: %w", err)
		}
	}

	return destPath, nil
}

// hardlinkDirectory recursively creates hardlinks for all files in a directory
func (a *App) hardlinkDirectory(srcDir, destDir string) error {
	// Create destination directory
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", srcDir, err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			// Recursively handle subdirectories
			if err := a.hardlinkDirectory(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Create hardlink for files
			if err := os.Link(srcPath, destPath); err != nil {
				return fmt.Errorf("failed to create hardlink %s -> %s: %w", srcPath, destPath, err)
			}
		}
	}

	return nil
}
