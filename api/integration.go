package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/anacrolix/torrent/metainfo"
)

// ReleaseInfo matches the typescript interface
type ReleaseInfo struct {
	Title             string   `json:"title"`
	Year              string   `json:"year"`
	Season            string   `json:"season"`
	Episode           string   `json:"episode"`
	Resolution        string   `json:"resolution"`
	Source            string   `json:"source"`
	Codec             string   `json:"codec"`
	Audio             string   `json:"audio"`
	AudioCodecs       []string `json:"audioCodecs"`
	AudioChannels     string   `json:"audioChannels"`
	Language          string   `json:"language"`
	AudioLanguages    []string `json:"audioLanguages"`
	SubtitleLanguages []string `json:"subtitleLanguages"`
	Hdr               []string `json:"hdr"`
	Tags              []string `json:"tags"`
	ReleaseGroup      string   `json:"releaseGroup"`
	Container         string   `json:"container"`
	Genres            []string `json:"genres"`
}

// Local Tags Structure (from tags.json)
type LocalMetaRoot struct {
	Categories []LocalCategory `json:"quaiprincipalcategories"`
}

type LocalCategory struct {
	Name            string                `json:"name"`
	Slug            string                `json:"slug"`
	ID              string                `json:"id"`
	SubCategories   []LocalCategory       `json:"emplacementsouscategorie"`
	Characteristics []LocalCharacteristic `json:"caracteristiques"`
}

type LocalCharacteristic struct {
	Name string     `json:"name"`
	Slug string     `json:"slug"`
	Tags []LocalTag `json:"tags"`
}

type LocalTag struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// RemoveFromQBittorrent removes the torrent from qBittorrent without deleting files
func (a *App) RemoveFromQBittorrent(torrentPath string, qbitUrl string, username string, password string) error {
	if qbitUrl == "" {
		return nil
	}

	// 1. Get InfoHash from file
	mi, err := metainfo.LoadFromFile(torrentPath)
	if err != nil {
		return fmt.Errorf("failed to load torrent file: %w", err)
	}
	infoHash := mi.HashInfoBytes().HexString()

	// 2. Login
	qbitUrl = strings.TrimSuffix(qbitUrl, "/")
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	client := &http.Client{Jar: jar}

	if username != "" || password != "" {
		vals := url.Values{}
		vals.Set("username", username)
		vals.Set("password", password)
		resp, err := client.PostForm(qbitUrl+"/api/v2/auth/login", vals)
		if err != nil {
			return fmt.Errorf("failed to login to qBittorrent: %w", err)
		}
		defer resp.Body.Close()
		// Some versions ensure cookie is set
	}

	// 3. Delete
	vals := url.Values{}
	vals.Set("hashes", infoHash)
	vals.Set("deleteFiles", "false")

	req, err := http.NewRequest("POST", qbitUrl+"/api/v2/torrents/delete", strings.NewReader(vals.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete torrent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to delete torrent, status: %d", resp.StatusCode)
	}

	return nil
}

// UploadToQBittorrent uploads the .torrent file to the configured qBittorrent instance
func (a *App) UploadToQBittorrent(torrentPath string, qbitUrl string, username string, password string) error {
	if qbitUrl == "" {
		return fmt.Errorf("qBittorrent URL is not configured")
	}

	// Remove trailing slash
	qbitUrl = strings.TrimSuffix(qbitUrl, "/")

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	client := &http.Client{Jar: jar}

	// 1. Login
	if username != "" || password != "" {
		vals := url.Values{}
		vals.Set("username", username)
		vals.Set("password", password)
		resp, err := client.PostForm(qbitUrl+"/api/v2/auth/login", vals)
		if err != nil {
			return fmt.Errorf("failed to login to qBittorrent: %w", err)
		}
		defer resp.Body.Close()

		// Some qbit versions return 200 even on failure with "Fails." in body
		body, _ := io.ReadAll(resp.Body)
		if string(body) == "Fails." {
			return fmt.Errorf("qBittorrent login failed: invalid credentials")
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("qBittorrent login failed: status %d", resp.StatusCode)
		}
	}

	// 2. Add Torrent
	// Prepare multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	file, err := os.Open(torrentPath)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("torrents", "release.torrent")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	// Ensure torrent starts automatically
	writer.WriteField("paused", "false")

	writer.Close()

	req, err := http.NewRequest("POST", qbitUrl+"/api/v2/torrents/add", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to add torrent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to add torrent, status: %d", resp.StatusCode)
	}

	return nil
}

// UploadToLaCale uploads the release metadata and files to la-cale.space
func (a *App) UploadToLaCale(torrentPath string, nfoPath string, title string, description string, tmdbId string, mediaType string, releaseInfo ReleaseInfo, passkey string, apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("La-Cale API key is missing in settings")
	}

	// 1. Fetch Metadata (Load from embedded tagsData)
	var meta LocalMetaRoot
	if err := json.Unmarshal([]byte(tagsData), &meta); err != nil {
		return fmt.Errorf("failed to parse embedded tags data: %w", err)
	}

	// 2. Identify Category
	categoryId, relevantChars := findLocalCategory(meta.Categories, mediaType)
	if categoryId == "" {
		return fmt.Errorf("could not find a matching category for type: %s", mediaType)
	}

	// 3. Identify Tags
	matchedTags := findLocalMatchingTags(relevantChars, releaseInfo)

	// 4. Upload (External API with API Key)
	client := &http.Client{}
	uploadURL := "https://la-cale.space/api/external/upload"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Fields
	if passkey != "" {
		writer.WriteField("passkey", passkey)
	}
	writer.WriteField("title", title)
	writer.WriteField("description", description)
	writer.WriteField("categoryId", categoryId)
	if tmdbId != "" {
		writer.WriteField("tmdbId", tmdbId)
		// Map simple mediaType to likely tmdb type
		tmdbType := "MOVIE"
		if mediaType == "episode" || mediaType == "season" {
			tmdbType = "TV"
		}
		writer.WriteField("tmdbType", tmdbType)
	}

	for _, tag := range matchedTags {
		writer.WriteField("tags", tag)
	}

	tFile, err := os.Open(torrentPath)
	if err != nil {
		return err
	}
	defer tFile.Close()

	// Create Torrent Part custom
	h := make(map[string][]string)
	h["Content-Disposition"] = []string{fmt.Sprintf(`form-data; name="file"; filename="%s.torrent"`, title)}
	h["Content-Type"] = []string{"application/x-bittorrent"}
	tPart, err := writer.CreatePart(h)
	if err != nil {
		return err
	}
	io.Copy(tPart, tFile)

	// NFO
	nFile, err := os.Open(nfoPath)
	if err != nil {
		return err
	}
	defer nFile.Close()

	// Create NFO Part custom
	hNfo := make(map[string][]string)
	hNfo["Content-Disposition"] = []string{fmt.Sprintf(`form-data; name="nfoFile"; filename="%s.nfo"`, title)}
	hNfo["Content-Type"] = []string{"text/x-nfo"}
	nPart, err := writer.CreatePart(hNfo)
	if err != nil {
		return err
	}
	io.Copy(nPart, nFile)

	writer.Close()

	// Debug Print (Truncated for sanity, but enough to see structure)
	fmt.Printf("--- Payload Preview (First 2000 chars) ---\n%s\n--- End Preview ---\n", string(body.Bytes()[:min(2000, body.Len())]))

	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Api-Key", apiKey)

	uploadResp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("upload request failed: %w", err)
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != 200 {
		respBody, err := io.ReadAll(uploadResp.Body)
		if err != nil {
			return fmt.Errorf("API Error (Status %d) - Failed to read body: %v", uploadResp.StatusCode, err)
		}
		return fmt.Errorf("API Error (Status %d) | Response: %s", uploadResp.StatusCode, string(respBody))
	}

	return nil
}


// Helpers

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func findLocalCategory(categories []LocalCategory, mediaType string) (string, []LocalCharacteristic) {
	// Recursive search for keywords
	// mediaType: "movie", "episode", "season", "ebook" -> "films", "series", "e-books"

	keywords := []string{}
	// Map mediaType to potential slugs in tags.json
	// Movie -> Video -> Films
	// Episode/Season -> Video -> Séries TV
	// Ebook -> E-books

	switch mediaType {
	case "movie":
		keywords = []string{"films", "film"}
	case "ebook":
		keywords = []string{"e-books", "ebooks", "ebook", "e-book"}
	case "game":
		keywords = []string{"jeux", "jeu", "games", "game", "pc"}
	default:
		keywords = []string{"series", "séries", "serie"}
	}

	for _, cat := range categories {
		// Verify if we are in "Vidéo" branch for movies/series
		// If root category is "Vidéo", search children

		// Recursion in children
		if len(cat.SubCategories) > 0 {
			if id, chars := findLocalCategory(cat.SubCategories, mediaType); id != "" {
				// If found in children, merge characteristics?
				// Usually characteristics are at leaf or inherited.
				// The JSON shows characteristics at leaf level mostly (e.g. Films has Genres, Codec...)
				// But root "Vidéo" might provide shared ones? JSON says "Vidéo" has emplacementsouscategorie but no direct caracteristiques listed in valid json provided in prompt for Video root, wait,
				// "Vidéo" has "emplacementsouscategorie".
				return id, chars
			}
		}

		// Check self path
		lowerSlug := strings.ToLower(cat.Slug)
		for _, kw := range keywords {
			if strings.Contains(lowerSlug, kw) {
				if cat.ID != "" {
					return cat.ID, cat.Characteristics
				}
				// If ID is missing but slug matches, maybe we are at correct category but ID logic failed?
				// But we expect ID.
				// For now return Slug and let uploader fail (or we fix logic)
				// Actually, we want ID. If new JSON has ID, we should use it.
				// Fallback to Slug just in case
				// return cat.Slug, cat.Characteristics
				// The previous code returned Slug.
				// Now we return ID if possible.
			}
		}
	}
	return "", nil
}

func findLocalMatchingTags(characteristics []LocalCharacteristic, info ReleaseInfo) []string {
	matched := []string{}
	unique := make(map[string]bool)

	addTag := func(t LocalTag) {
		if !unique[t.ID] && t.ID != "" {
			unique[t.ID] = true
			matched = append(matched, t.ID)
		}
	}

	for _, char := range characteristics {
		slug := strings.ToLower(char.Slug)
		tagsFoundForChar := false

		var valuesToCheck []string

		// Map characteristic to info field based on slug
		switch {
		case strings.Contains(slug, "genre"):
			valuesToCheck = info.Genres
		case strings.Contains(slug, "qualit") || strings.Contains(slug, "resolution"):
			valuesToCheck = []string{info.Resolution}
		case strings.Contains(slug, "codec-vid") || strings.Contains(slug, "codec-video"):
			valuesToCheck = []string{info.Codec}
			// Fallback: If codec is "x264", maybe tag "AVC/H264/x264" matches
		case strings.Contains(slug, "codec-audio"):
			valuesToCheck = info.AudioCodecs
			if len(valuesToCheck) == 0 && info.Audio != "" {
				valuesToCheck = []string{info.Audio}
			}
		case strings.Contains(slug, "langues") || strings.Contains(slug, "langue"):
			valuesToCheck = info.AudioLanguages
			if len(valuesToCheck) == 0 && info.Language != "" {
				valuesToCheck = []string{info.Language}
			}
		case strings.Contains(slug, "sous-titres"):
			valuesToCheck = info.SubtitleLanguages
		case strings.Contains(slug, "extension") || strings.Contains(slug, "format"):
			valuesToCheck = []string{info.Container}
		case strings.Contains(slug, "source"):
			valuesToCheck = []string{info.Source}
		case strings.Contains(slug, "caract") || strings.Contains(slug, "hdr"):
			valuesToCheck = info.Hdr
			valuesToCheck = append(valuesToCheck, info.Tags...)
		}

		// Filter empty
		validValues := []string{}

		// If dealing with Genres, try to map English -> French common TMDB values
		isGenre := strings.Contains(slug, "genre")

		for _, v := range valuesToCheck {
			if v != "" {
				validValues = append(validValues, strings.ToLower(v))
				if isGenre {
					// Add translated fallback for common English terms
					lowerV := strings.ToLower(v)
					switch lowerV {
					case "adventure":
						validValues = append(validValues, "aventure")
					case "fantasy":
						validValues = append(validValues, "fantastique")
					case "science fiction", "sci-fi":
						validValues = append(validValues, "science-fiction")
					case "mystery":
						validValues = append(validValues, "mystere") // Normalized usually handles accent, but spelling differs
					case "war":
						validValues = append(validValues, "guerre")
					case "family":
						validValues = append(validValues, "famille")
					case "history":
						validValues = append(validValues, "historique")
					case "comedy":
						validValues = append(validValues, "comedie")
					case "action & adventure":
						validValues = append(validValues, "action", "aventure")
					case "sci-fi & fantasy", "science-fiction & fantastique":
						validValues = append(validValues, "science-fiction", "fantastique")
					}
				}
			}
		}

		// If no values to check for this characteristic, we might want to skip "Autre" logic?
		// Actually, if we have no info about Audio Codec, should we set "Autre"?
		// Probably not, "Autre" usually implies "I have a value but it's not in the list".
		// But user said: "if the audio codec is not found in the tags ... there should be a "Autre" tag".
		// If Info.Audio is empty, then "not found" fits? Or does it mean "Detected but not matched"?
		// I assume if we detected something. If we didn't detect Audio Codec, we shouldn't tag it.
		if len(validValues) == 0 {
			continue
		}

		for _, tag := range char.Tags {
			isMatch := false

			// Normalize for comparison: remove hyphens, spaces, accents
			normalize := func(s string) string {
				s = strings.ToLower(s)
				s = strings.ReplaceAll(s, "-", "")
				s = strings.ReplaceAll(s, " ", "")
				// Simple accent removal if needed (can be expanded)
				s = strings.ReplaceAll(s, "é", "e")
				s = strings.ReplaceAll(s, "è", "e")
				return s
			}

			normTag := normalize(tag.Name)

			for _, val := range validValues {
				normVal := normalize(val)

				// Fuzzy match logic
				if normTag == normVal {
					isMatch = true
				} else if strings.Contains(normTag, normVal) {
					isMatch = true
				} else if strings.Contains(normVal, normTag) {
					isMatch = true
				} else {
					// Direct check without normalization for some cases?
					// Retain original check just in case
					if strings.Contains(strings.ToLower(tag.Name), strings.ToLower(val)) {
						isMatch = true
					}
				}

				if isMatch {
					break
				}
			}

			if isMatch {
				addTag(tag)
				tagsFoundForChar = true
			}
		}

		// Fallback "Autre"
		if !tagsFoundForChar {
			for _, tag := range char.Tags {
				if strings.EqualFold(tag.Name, "Autre") || strings.EqualFold(tag.Name, "Autres") || strings.HasPrefix(strings.ToLower(tag.Name), "autre") {
					addTag(tag)
					break
				}
			}
		}
	}

	return matched
}

// ==================== TRANSMISSION INTEGRATION ====================

// TransmissionRPCRequest represents a Transmission RPC request
type TransmissionRPCRequest struct {
	Method    string      `json:"method"`
	Arguments interface{} `json:"arguments,omitempty"`
}

// TransmissionRPCResponse represents a Transmission RPC response
type TransmissionRPCResponse struct {
	Result    string                 `json:"result"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// UploadToTransmission uploads a torrent to Transmission via RPC
func (a *App) UploadToTransmission(torrentPath string, transmissionUrl string, username string, password string) error {
	if transmissionUrl == "" {
		return fmt.Errorf("Transmission URL is not configured")
	}

	transmissionUrl = strings.TrimSuffix(transmissionUrl, "/")
	rpcUrl := transmissionUrl + "/transmission/rpc"

	// Read torrent file and encode to base64
	torrentData, err := os.ReadFile(torrentPath)
	if err != nil {
		return fmt.Errorf("failed to read torrent file: %w", err)
	}
	b64Torrent := base64.StdEncoding.EncodeToString(torrentData)

	client := &http.Client{}
	sessionId := ""

	// Helper to make RPC request
	makeRequest := func(req TransmissionRPCRequest) (*TransmissionRPCResponse, error) {
		body, _ := json.Marshal(req)
		httpReq, err := http.NewRequest("POST", rpcUrl, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		httpReq.Header.Set("Content-Type", "application/json")
		if sessionId != "" {
			httpReq.Header.Set("X-Transmission-Session-Id", sessionId)
		}
		if username != "" {
			httpReq.SetBasicAuth(username, password)
		}

		resp, err := client.Do(httpReq)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Handle CSRF token
		if resp.StatusCode == 409 {
			sessionId = resp.Header.Get("X-Transmission-Session-Id")
			return nil, fmt.Errorf("retry with session id")
		}

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}

		var rpcResp TransmissionRPCResponse
		if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
			return nil, err
		}
		return &rpcResp, nil
	}

	addReq := TransmissionRPCRequest{
		Method: "torrent-add",
		Arguments: map[string]interface{}{
			"metainfo": b64Torrent,
			"paused":   false,
		},
	}

	// First attempt (will likely fail to get session ID)
	resp, err := makeRequest(addReq)
	if err != nil && strings.Contains(err.Error(), "retry with session id") {
		// Retry with session ID
		resp, err = makeRequest(addReq)
	}
	if err != nil {
		return fmt.Errorf("failed to add torrent to Transmission: %w", err)
	}

	if resp.Result != "success" {
		return fmt.Errorf("Transmission returned: %s", resp.Result)
	}

	return nil
}

// RemoveFromTransmission removes a torrent from Transmission
func (a *App) RemoveFromTransmission(torrentPath string, transmissionUrl string, username string, password string) error {
	if transmissionUrl == "" {
		return nil
	}

	// Get info hash
	mi, err := metainfo.LoadFromFile(torrentPath)
	if err != nil {
		return fmt.Errorf("failed to load torrent file: %w", err)
	}
	infoHash := mi.HashInfoBytes().HexString()

	transmissionUrl = strings.TrimSuffix(transmissionUrl, "/")
	rpcUrl := transmissionUrl + "/transmission/rpc"

	client := &http.Client{}
	sessionId := ""

	makeRequest := func(req TransmissionRPCRequest) (*TransmissionRPCResponse, error) {
		body, _ := json.Marshal(req)
		httpReq, err := http.NewRequest("POST", rpcUrl, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		httpReq.Header.Set("Content-Type", "application/json")
		if sessionId != "" {
			httpReq.Header.Set("X-Transmission-Session-Id", sessionId)
		}
		if username != "" {
			httpReq.SetBasicAuth(username, password)
		}

		resp, err := client.Do(httpReq)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 409 {
			sessionId = resp.Header.Get("X-Transmission-Session-Id")
			return nil, fmt.Errorf("retry with session id")
		}

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}

		var rpcResp TransmissionRPCResponse
		if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
			return nil, err
		}
		return &rpcResp, nil
	}

	removeReq := TransmissionRPCRequest{
		Method: "torrent-remove",
		Arguments: map[string]interface{}{
			"ids":               []string{infoHash},
			"delete-local-data": false,
		},
	}

	resp, err := makeRequest(removeReq)
	if err != nil && strings.Contains(err.Error(), "retry with session id") {
		resp, err = makeRequest(removeReq)
	}
	if err != nil {
		return fmt.Errorf("failed to remove torrent from Transmission: %w", err)
	}

	if resp.Result != "success" {
		return fmt.Errorf("Transmission returned: %s", resp.Result)
	}

	return nil
}

// ==================== DELUGE INTEGRATION ====================

// DelugeRPCRequest represents a Deluge JSON-RPC request
type DelugeRPCRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     int           `json:"id"`
}

// DelugeRPCResponse represents a Deluge JSON-RPC response
type DelugeRPCResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
	ID     int         `json:"id"`
}

// UploadToDeluge uploads a torrent to Deluge via JSON-RPC
func (a *App) UploadToDeluge(torrentPath string, delugeUrl string, password string) error {
	if delugeUrl == "" {
		return fmt.Errorf("Deluge URL is not configured")
	}

	delugeUrl = strings.TrimSuffix(delugeUrl, "/")
	rpcUrl := delugeUrl + "/json"

	// Read and encode torrent
	torrentData, err := os.ReadFile(torrentPath)
	if err != nil {
		return fmt.Errorf("failed to read torrent file: %w", err)
	}
	b64Torrent := base64.StdEncoding.EncodeToString(torrentData)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	requestId := 1

	makeRequest := func(method string, params []interface{}) (*DelugeRPCResponse, error) {
		req := DelugeRPCRequest{
			Method: method,
			Params: params,
			ID:     requestId,
		}
		requestId++

		body, _ := json.Marshal(req)
		httpReq, err := http.NewRequest("POST", rpcUrl, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(httpReq)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var rpcResp DelugeRPCResponse
		if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
			return nil, err
		}
		if rpcResp.Error != nil {
			return nil, fmt.Errorf("Deluge error: %v", rpcResp.Error)
		}
		return &rpcResp, nil
	}

	// 1. Login
	_, err = makeRequest("auth.login", []interface{}{password})
	if err != nil {
		return fmt.Errorf("Deluge login failed: %w", err)
	}

	// 2. Check connection (optional but recommended)
	connResp, err := makeRequest("web.connected", []interface{}{})
	if err != nil {
		return fmt.Errorf("failed to check Deluge connection: %w", err)
	}
	if connResp.Result == false {
		// Try to connect to first available host
		hostsResp, _ := makeRequest("web.get_hosts", []interface{}{})
		if hosts, ok := hostsResp.Result.([]interface{}); ok && len(hosts) > 0 {
			if host, ok := hosts[0].([]interface{}); ok && len(host) > 0 {
				makeRequest("web.connect", []interface{}{host[0]})
			}
		}
	}

	// 3. Add torrent
	filename := filepath.Base(torrentPath)
	_, err = makeRequest("web.add_torrents", []interface{}{
		[]map[string]interface{}{
			{
				"path":    filename,
				"options": map[string]interface{}{"add_paused": false},
			},
		},
	})
	// Alternative method: core.add_torrent_file
	_, err = makeRequest("core.add_torrent_file", []interface{}{
		filename,
		b64Torrent,
		map[string]interface{}{"add_paused": false},
	})
	if err != nil {
		return fmt.Errorf("failed to add torrent to Deluge: %w", err)
	}

	return nil
}

// RemoveFromDeluge removes a torrent from Deluge
func (a *App) RemoveFromDeluge(torrentPath string, delugeUrl string, password string) error {
	if delugeUrl == "" {
		return nil
	}

	// Get info hash
	mi, err := metainfo.LoadFromFile(torrentPath)
	if err != nil {
		return fmt.Errorf("failed to load torrent file: %w", err)
	}
	infoHash := strings.ToLower(mi.HashInfoBytes().HexString())

	delugeUrl = strings.TrimSuffix(delugeUrl, "/")
	rpcUrl := delugeUrl + "/json"

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	requestId := 1

	makeRequest := func(method string, params []interface{}) (*DelugeRPCResponse, error) {
		req := DelugeRPCRequest{
			Method: method,
			Params: params,
			ID:     requestId,
		}
		requestId++

		body, _ := json.Marshal(req)
		httpReq, err := http.NewRequest("POST", rpcUrl, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(httpReq)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var rpcResp DelugeRPCResponse
		if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
			return nil, err
		}
		return &rpcResp, nil
	}

	// Login
	_, err = makeRequest("auth.login", []interface{}{password})
	if err != nil {
		return fmt.Errorf("Deluge login failed: %w", err)
	}

	// Remove torrent (false = don't delete data)
	_, err = makeRequest("core.remove_torrent", []interface{}{infoHash, false})
	if err != nil {
		return fmt.Errorf("failed to remove torrent from Deluge: %w", err)
	}

	return nil
}

// ==================== GENERIC TORRENT CLIENT DISPATCHER ====================

// UploadToTorrentClient uploads a torrent to the configured client
func (a *App) UploadToTorrentClient(torrentPath string, settings AppSettings) error {
	switch settings.TorrentClient {
	case "qbittorrent":
		return a.UploadToQBittorrent(torrentPath, settings.QbitUrl, settings.QbitUsername, settings.QbitPassword)
	case "transmission":
		return a.UploadToTransmission(torrentPath, settings.TransmissionUrl, settings.TransmissionUsername, settings.TransmissionPassword)
	case "deluge":
		return a.UploadToDeluge(torrentPath, settings.DelugeUrl, settings.DelugePassword)
	case "none", "":
		return nil // No client configured, skip upload
	default:
		return fmt.Errorf("unknown torrent client: %s", settings.TorrentClient)
	}
}

// RemoveFromTorrentClient removes a torrent from the configured client
func (a *App) RemoveFromTorrentClient(torrentPath string, settings AppSettings) error {
	switch settings.TorrentClient {
	case "qbittorrent":
		return a.RemoveFromQBittorrent(torrentPath, settings.QbitUrl, settings.QbitUsername, settings.QbitPassword)
	case "transmission":
		return a.RemoveFromTransmission(torrentPath, settings.TransmissionUrl, settings.TransmissionUsername, settings.TransmissionPassword)
	case "deluge":
		return a.RemoveFromDeluge(torrentPath, settings.DelugeUrl, settings.DelugePassword)
	case "none", "":
		return nil
	default:
		return fmt.Errorf("unknown torrent client: %s", settings.TorrentClient)
	}
}
