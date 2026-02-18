package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
)

// TorrentTask represents an async torrent creation task
type TorrentTask struct {
	ID          string  `json:"id"`
	Status      string  `json:"status"` // "hashing", "building", "done", "error"
	Progress    float64 `json:"progress"` // 0.0 to 1.0
	TorrentPath string  `json:"torrentPath,omitempty"`
	Error       string  `json:"error,omitempty"`
	TotalBytes  int64   `json:"totalBytes"`
	HashedBytes int64   `json:"hashedBytes"`
}

// TaskManager manages async torrent creation tasks
type TaskManager struct {
	mu    sync.RWMutex
	tasks map[string]*TorrentTask
}

// NewTaskManager creates a new TaskManager
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*TorrentTask),
	}
}

func generateTaskID() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 12)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return fmt.Sprintf("%d_%s", time.Now().UnixMilli(), string(b))
}

func (tm *TaskManager) GetTask(id string) *TorrentTask {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	t, ok := tm.tasks[id]
	if !ok {
		return nil
	}
	// Return a copy
	cp := *t
	return &cp
}

func (tm *TaskManager) createTask() *TorrentTask {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	task := &TorrentTask{
		ID:     generateTaskID(),
		Status: "hashing",
	}
	tm.tasks[task.ID] = task
	return task
}

func (tm *TaskManager) updateTask(id string, fn func(t *TorrentTask)) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if t, ok := tm.tasks[id]; ok {
		fn(t)
	}
}

// CleanOldTasks removes completed/errored tasks older than 10 minutes
func (tm *TaskManager) CleanOldTasks() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	for id, t := range tm.tasks {
		if t.Status == "done" || t.Status == "error" {
			delete(tm.tasks, id)
		}
	}
}

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
type App struct {
	TaskMgr *TaskManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		TaskMgr: NewTaskManager(),
	}
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

// findFirstVideoFile searches for the first video file in a directory (sorted alphabetically).
// It checks the top level first, then one level of subdirectories.
func findFirstVideoFile(dirPath string) string {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return ""
	}
	// First pass: look for video files directly in the directory
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if isVideoFile(ext) {
			return filepath.Join(dirPath, entry.Name())
		}
	}
	// Second pass: look one level deep in subdirectories
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		subEntries, err := os.ReadDir(filepath.Join(dirPath, entry.Name()))
		if err != nil {
			continue
		}
		for _, subEntry := range subEntries {
			if subEntry.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(subEntry.Name()))
			if isVideoFile(ext) {
				return filepath.Join(dirPath, entry.Name(), subEntry.Name())
			}
		}
	}
	return ""
}

// GetMediaInfo executes mediainfo command on the file and returns output
// If filePath is a directory, it finds the first video file inside for analysis
func (a *App) GetMediaInfo(filePath string) (string, error) {
	// Check if path is a directory; if so, find the first video file
	info, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot stat path: %w", err)
	}
	if info.IsDir() {
		videoFile := findFirstVideoFile(filePath)
		if videoFile != "" {
			filePath = videoFile
		}
		// If no video file found, pass directory as-is to mediainfo (legacy fallback)
	}

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

// mediaTypeDirName maps a media type string to the corresponding subdirectory name
func mediaTypeDirName(mediaType string) string {
	switch mediaType {
	case "movie":
		return "Films"
	case "episode", "season":
		return "Séries"
	case "ebook":
		return "Ebooks"
	case "game":
		return "Jeux"
	default:
		return "Autres"
	}
}

// resolveOutputDir builds and creates the output directory:
// {outputDir}/{MediaTypeDir}/{torrentName}/
func resolveOutputDir(outputDir, mediaType, torrentName string) (string, error) {
	if outputDir == "" {
		outputDir = "/torrents"
	}
	dir := filepath.Join(outputDir, mediaTypeDirName(mediaType), torrentName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("cannot create output directory %s: %w", dir, err)
	}
	return dir, nil
}

// calculatePieceLength returns an appropriate piece length based on total content size
func calculatePieceLength(sourcePath string) int64 {
	var totalSize int64
	filepath.Walk(sourcePath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case totalSize > 50*GB:
		return 16 * MB
	case totalSize > 16*GB:
		return 8 * MB
	case totalSize > 8*GB:
		return 4 * MB
	case totalSize > 4*GB:
		return 2 * MB
	case totalSize > 1*GB:
		return 1 * MB
	default:
		return 512 * KB
	}
}

// CreateTorrent creates a .torrent file for the given source path
// torrentName is the name that will appear in the torrent (the release name)
func (a *App) CreateTorrent(sourcePath string, trackers []string, comment string, isPrivate bool, torrentName string) (string, error) {
	pieceLength := calculatePieceLength(sourcePath)
	info := metainfo.Info{
		PieceLength: pieceLength,
	}

	if isPrivate {
		info.Private = new(bool)
		*info.Private = true
	}

	info.Source = "lacale"

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

	// Determine output path using configured outputDir
	var baseName string
	if torrentName != "" {
		baseName = torrentName
	} else {
		baseName = filepath.Base(sourcePath)
	}
	outDir, err := resolveOutputDir("/torrents", "", baseName)
	if err != nil {
		return "", err
	}
	outputPath := filepath.Join(outDir, baseName+".torrent")

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

// StartCreateTorrent launches async torrent creation and returns a task ID
func (a *App) StartCreateTorrent(sourcePath string, trackers []string, comment string, isPrivate bool, torrentName string, outputDir string, mediaType string) string {
	task := a.TaskMgr.createTask()
	taskID := task.ID

	go func() {
		result, err := a.createTorrentWithProgress(taskID, sourcePath, trackers, comment, isPrivate, torrentName, outputDir, mediaType)
		if err != nil {
			a.TaskMgr.updateTask(taskID, func(t *TorrentTask) {
				t.Status = "error"
				t.Error = err.Error()
			})
			return
		}
		a.TaskMgr.updateTask(taskID, func(t *TorrentTask) {
			t.Status = "done"
			t.TorrentPath = result
			t.Progress = 1.0
		})
	}()

	return taskID
}

// createTorrentWithProgress builds a torrent file with progress tracking
func (a *App) createTorrentWithProgress(taskID string, sourcePath string, trackers []string, comment string, isPrivate bool, torrentName string, outputDir string, mediaType string) (string, error) {
	pieceLength := calculatePieceLength(sourcePath)

	// Phase 1: Walk filesystem to build file list and total size
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return "", err
	}
	isDir := sourceInfo.IsDir()

	type fileEntry struct {
		absPath string
		relPath []string
		length  int64
	}

	var fileList []fileEntry
	var totalSize int64

	if isDir {
		filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			relPath, _ := filepath.Rel(sourcePath, path)
			parts := strings.Split(relPath, string(os.PathSeparator))
			fileList = append(fileList, fileEntry{
				absPath: path,
				relPath: parts,
				length:  info.Size(),
			})
			totalSize += info.Size()
			return nil
		})
	} else {
		fileList = append(fileList, fileEntry{
			absPath: sourcePath,
			length:  sourceInfo.Size(),
		})
		totalSize = sourceInfo.Size()
	}

	a.TaskMgr.updateTask(taskID, func(t *TorrentTask) {
		t.TotalBytes = totalSize
	})

	// Phase 2: Hash pieces with progress tracking
	var pieces []byte
	var hashedBytes int64

	pieceBuf := make([]byte, pieceLength)
	pieceBufLen := int64(0)
	readBuf := make([]byte, 64*1024) // 64KB read buffer
	lastUpdate := time.Now()

	flushPiece := func() {
		h := sha1.Sum(pieceBuf[:pieceBufLen])
		pieces = append(pieces, h[:]...)
		pieceBufLen = 0
	}

	for _, fe := range fileList {
		f, err := os.Open(fe.absPath)
		if err != nil {
			return "", fmt.Errorf("open %s: %w", fe.absPath, err)
		}

		for {
			remaining := pieceLength - pieceBufLen
			toRead := int64(len(readBuf))
			if toRead > remaining {
				toRead = remaining
			}

			n, readErr := f.Read(readBuf[:toRead])
			if n > 0 {
				copy(pieceBuf[pieceBufLen:], readBuf[:n])
				pieceBufLen += int64(n)
				hashedBytes += int64(n)

				if pieceBufLen == pieceLength {
					flushPiece()
				}

				// Update progress every 500ms to avoid lock contention
				if time.Since(lastUpdate) > 500*time.Millisecond {
					hb := hashedBytes
					a.TaskMgr.updateTask(taskID, func(t *TorrentTask) {
						t.HashedBytes = hb
						t.Progress = float64(hb) / float64(totalSize)
					})
					lastUpdate = time.Now()
				}
			}
			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				f.Close()
				return "", fmt.Errorf("read %s: %w", fe.absPath, readErr)
			}
		}
		f.Close()
	}

	// Flush remaining partial piece
	if pieceBufLen > 0 {
		flushPiece()
	}

	// Final progress update
	a.TaskMgr.updateTask(taskID, func(t *TorrentTask) {
		t.Status = "building"
		t.HashedBytes = totalSize
		t.Progress = 1.0
	})

	// Phase 3: Build metainfo struct
	info := metainfo.Info{
		PieceLength: pieceLength,
		Pieces:      pieces,
		Name:        filepath.Base(sourcePath),
	}

	if isPrivate {
		info.Private = new(bool)
		*info.Private = true
	}
	info.Source = "lacale"

	if isDir {
		for _, fe := range fileList {
			info.Files = append(info.Files, metainfo.FileInfo{
				Path:   fe.relPath,
				Length: fe.length,
			})
		}
	} else {
		info.Length = totalSize
	}

	if torrentName != "" {
		info.Name = torrentName
	}

	// Phase 4: Encode and write .torrent file
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

	var baseName string
	if torrentName != "" {
		baseName = torrentName
	} else {
		baseName = filepath.Base(sourcePath)
	}

	outDir, err := resolveOutputDir(outputDir, mediaType, baseName)
	if err != nil {
		return "", err
	}
	outputPath := filepath.Join(outDir, baseName+".torrent")

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
// outputDir and mediaType are used to build the output path:
// {outputDir}/{MediaTypeDir}/{torrentName}/{torrentName}.nfo
func (a *App) SaveNfo(sourcePath string, content string, torrentName string, outputDir string, mediaType string) (string, error) {
	// Determine base name: use torrentName if provided, otherwise derive from source
	var baseName string
	if torrentName != "" {
		baseName = torrentName
	} else {
		baseName = filepath.Base(sourcePath)
		ext := filepath.Ext(sourcePath)
		lowerExt := strings.ToLower(ext)
		if isMediaFile(lowerExt) {
			baseName = strings.TrimSuffix(baseName, ext)
		}
	}

	dir, err := resolveOutputDir(outputDir, mediaType, baseName)
	if err != nil {
		return "", err
	}

	outputPath := filepath.Join(dir, baseName+".nfo")
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
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

	// For files: the source's parent dir must not be used as hardlink destination
	// (would result in linking a file to itself, then deleting the original)
	sourceParent := filepath.Dir(filepath.Clean(sourcePath))

	for _, dir := range hardlinkDirs {
		if dir == "" {
			continue
		}
		cleanDir := filepath.Clean(dir)
		// Skip if the hardlink dir is the same as (or a parent of) the source's directory
		if cleanDir == sourceParent {
			continue
		}
		// Ensure the directory exists
		if _, err := os.Stat(cleanDir); os.IsNotExist(err) {
			continue
		}
		dirDevID, err := getDeviceID(cleanDir)
		if err != nil {
			continue
		}
		if dirDevID == sourceDevID {
			return cleanDir, nil
		}
	}

	return "", fmt.Errorf("no hardlink directory found on the same device as %s (source parent dir is excluded)", sourcePath)
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

	// Safety check: never create a hardlink that points to itself
	if filepath.Clean(destPath) == filepath.Clean(sourcePath) {
		return "", fmt.Errorf("source and destination are the same path, cannot hardlink: %s", sourcePath)
	}

	// If destination already exists, remove it first
	if _, err := os.Stat(destPath); err == nil {
		if err := os.RemoveAll(destPath); err != nil {
			return "", fmt.Errorf("failed to remove existing destination: %w", err)
		}
	}

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
