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

// Meta API response structs (from /api/external/meta)
type MetaResponse struct {
	Categories    []MetaCategory `json:"categories"`
	TagGroups     []MetaTagGroup `json:"tagGroups"`
	UngroupedTags []MetaTag      `json:"ungroupedTags"`
}

type MetaCategory struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Slug     string         `json:"slug"`
	Children []MetaCategory `json:"children"`
}

type MetaTagGroup struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
	Tags []MetaTag `json:"tags"`
}

type MetaTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// fetchLaCaleMetadata fetches categories and tags from the La-Cale API
func fetchLaCaleMetadata(apiKey string) (*MetaResponse, error) {
	req, err := http.NewRequest("GET", "https://la-cale.space/api/external/meta", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create meta request: %w", err)
	}
	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch meta: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("meta API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, _ := io.ReadAll(resp.Body)
	var meta MetaResponse
	if err := json.Unmarshal(body, &meta); err != nil {
		return nil, fmt.Errorf("failed to decode meta response: %w", err)
	}

	// Inject local tags for groups where the API returns null
	for i := range meta.TagGroups {
		if len(meta.TagGroups[i].Tags) == 0 {
			if localTags, ok := localTagsDB[meta.TagGroups[i].Slug]; ok {
				meta.TagGroups[i].Tags = localTags
			}
		}
	}

	return &meta, nil
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
func (a *App) UploadToLaCale(torrentPath string, nfoPath string, title string, description string, tmdbId string, mediaType string, releaseInfo ReleaseInfo, passkey string, apiKey string, overrideTags []string) error {
	if apiKey == "" {
		return fmt.Errorf("La-Cale API key is missing in settings")
	}

	// 1. Fetch Metadata from La-Cale API
	meta, err := fetchLaCaleMetadata(apiKey)
	if err != nil {
		return fmt.Errorf("failed to fetch La-Cale metadata: %w", err)
	}

	// 2. Identify Category
	categoryId := findCategory(meta.Categories, mediaType)
	if categoryId == "" {
		return fmt.Errorf("could not find a matching category for type: %s", mediaType)
	}

	// 3. Identify Tags (use overrideTags if provided)
	var matchedTags []string
	if len(overrideTags) > 0 {
		matchedTags = overrideTags
	} else {
		matchedTags = findMatchingTags(meta.TagGroups, releaseInfo)
	}

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
		fmt.Printf("UploadToLaCale: HTTP request error: %v\n", err)
		return fmt.Errorf("upload request failed: %w", err)
	}
	defer uploadResp.Body.Close()

	respBody, _ := io.ReadAll(uploadResp.Body)
	fmt.Printf("UploadToLaCale: Status=%d\n", uploadResp.StatusCode)
	fmt.Printf("UploadToLaCale: Response body: %s\n", string(respBody))
	fmt.Printf("UploadToLaCale: Response headers: %v\n", uploadResp.Header)

	if uploadResp.StatusCode < 200 || uploadResp.StatusCode >= 300 {
		return fmt.Errorf("API Error (Status %d) | Response: %s", uploadResp.StatusCode, string(respBody))
	}

	return nil
}


// Preview structs
type PreviewTag struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

type PreviewTagGroup struct {
	ID   string       `json:"id"`
	Name string       `json:"name"`
	Slug string       `json:"slug"`
	Tags []PreviewTag `json:"tags"`
}

type LaCalePreviewResponse struct {
	CategoryId   string            `json:"categoryId"`
	CategoryName string            `json:"categoryName"`
	MatchedTags  []PreviewTag      `json:"matchedTags"`
	AllTagGroups []PreviewTagGroup `json:"allTagGroups"`
}

// PreviewLaCale returns preview data (category, matched tags, all tag groups) without uploading
func (a *App) PreviewLaCale(mediaType string, releaseInfo ReleaseInfo, apiKey string) (*LaCalePreviewResponse, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("La-Cale API key is missing in settings")
	}

	meta, err := fetchLaCaleMetadata(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch La-Cale metadata: %w", err)
	}

	categoryId := findCategory(meta.Categories, mediaType)
	categoryName := findCategoryName(meta.Categories, categoryId)

	matchedTagIds := findMatchingTags(meta.TagGroups, releaseInfo)

	// Build a set of matched IDs for quick lookup
	matchedSet := make(map[string]bool)
	for _, id := range matchedTagIds {
		matchedSet[id] = true
	}

	// Build matched tags with names and groups
	var matchedTags []PreviewTag
	for _, group := range meta.TagGroups {
		for _, tag := range group.Tags {
			if matchedSet[tag.ID] {
				matchedTags = append(matchedTags, PreviewTag{
					ID:    tag.ID,
					Name:  tag.Name,
					Group: group.Name,
				})
			}
		}
	}

	// Build all tag groups
	var allTagGroups []PreviewTagGroup
	for _, group := range meta.TagGroups {
		pg := PreviewTagGroup{
			ID:   group.ID,
			Name: group.Name,
			Slug: group.Slug,
		}
		for _, tag := range group.Tags {
			pg.Tags = append(pg.Tags, PreviewTag{
				ID:    tag.ID,
				Name:  tag.Name,
				Group: group.Name,
			})
		}
		allTagGroups = append(allTagGroups, pg)
	}

	return &LaCalePreviewResponse{
		CategoryId:   categoryId,
		CategoryName: categoryName,
		MatchedTags:  matchedTags,
		AllTagGroups: allTagGroups,
	}, nil
}

// findCategoryName recursively finds the name of a category by its ID
func findCategoryName(categories []MetaCategory, targetId string) string {
	for _, cat := range categories {
		if cat.ID == targetId {
			return cat.Name
		}
		if len(cat.Children) > 0 {
			if name := findCategoryName(cat.Children, targetId); name != "" {
				return name
			}
		}
	}
	return ""
}

// Helpers

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func findCategory(categories []MetaCategory, mediaType string) string {
	keywords := []string{}
	switch mediaType {
	case "movie":
		keywords = []string{"films", "film"}
	case "ebook":
		keywords = []string{"e-books", "ebooks", "ebook", "e-book"}
	case "game":
		keywords = []string{"jeux", "jeu", "games", "game", "pc"}
	default: // episode, season
		keywords = []string{"series", "séries", "serie"}
	}

	var search func(cats []MetaCategory) string
	search = func(cats []MetaCategory) string {
		for _, cat := range cats {
			lowerSlug := strings.ToLower(cat.Slug)
			for _, kw := range keywords {
				if strings.Contains(lowerSlug, kw) && cat.ID != "" {
					return cat.ID
				}
			}
			if len(cat.Children) > 0 {
				if id := search(cat.Children); id != "" {
					return id
				}
			}
		}
		return ""
	}
	return search(categories)
}

func findMatchingTags(tagGroups []MetaTagGroup, info ReleaseInfo) []string {
	matched := []string{}
	unique := make(map[string]bool)

	addTag := func(t MetaTag) {
		if !unique[t.ID] && t.ID != "" {
			unique[t.ID] = true
			matched = append(matched, t.ID)
		}
	}

	normalize := func(s string) string {
		s = strings.ToLower(s)
		s = strings.ReplaceAll(s, "-", "")
		s = strings.ReplaceAll(s, " ", "")
		s = strings.ReplaceAll(s, "é", "e")
		s = strings.ReplaceAll(s, "è", "e")
		return s
	}

	for _, group := range tagGroups {
		slug := strings.ToLower(group.Slug)
		tagsFoundForGroup := false

		var valuesToCheck []string
		strictMatch := false // used for codec-audio: exact match only to avoid AC3 matching E-AC3

		switch {
		case strings.Contains(slug, "genre"):
			valuesToCheck = info.Genres
		case strings.Contains(slug, "qualit") || strings.Contains(slug, "resolution"):
			valuesToCheck = []string{info.Resolution}
		case strings.Contains(slug, "codec-vid") || strings.Contains(slug, "codec-video"):
			valuesToCheck = []string{info.Codec}
		case strings.Contains(slug, "codec-audio"):
			valuesToCheck = info.AudioCodecs
			if len(valuesToCheck) == 0 && info.Audio != "" {
				valuesToCheck = []string{info.Audio}
			}
			strictMatch = true
		case strings.Contains(slug, "langue"):
			valuesToCheck = info.AudioLanguages
			if len(valuesToCheck) == 0 && info.Language != "" {
				valuesToCheck = []string{info.Language}
			}
		case strings.Contains(slug, "sous-titre"):
			valuesToCheck = info.SubtitleLanguages
			// note: isSubtitle set below after switch
		case strings.Contains(slug, "extension") || strings.Contains(slug, "format"):
			valuesToCheck = []string{info.Container}
		case strings.Contains(slug, "source"):
			valuesToCheck = []string{info.Source}
		case strings.Contains(slug, "caract") || strings.Contains(slug, "hdr"):
			valuesToCheck = info.Hdr
			valuesToCheck = append(valuesToCheck, info.Tags...)
		default:
			continue
		}

		// Filter empty and build valid values with genre/language translations
		validValues := []string{}
		isGenre := strings.Contains(slug, "genre")
		isLang := strings.Contains(slug, "langue")
		isSubtitle := strings.Contains(slug, "sous-titre")

		for _, v := range valuesToCheck {
			if v != "" {
				validValues = append(validValues, strings.ToLower(v))
				lowerV := strings.ToLower(v)
				if isGenre {
					switch lowerV {
					case "adventure":
						validValues = append(validValues, "aventure")
					case "fantasy":
						validValues = append(validValues, "fantastique")
					case "science fiction", "sci-fi":
						validValues = append(validValues, "science-fiction")
					case "mystery":
						validValues = append(validValues, "mystere")
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
				if isLang || isSubtitle {
					// Add both French and English equivalents for language/subtitle matching
					// For subtitles, also add short codes (fr, eng) to match tags like "FR", "ENG"
					switch lowerV {
					case "français", "french":
						validValues = append(validValues, "french", "français", "fr")
					case "anglais", "english":
						validValues = append(validValues, "english", "anglais", "eng")
					case "espagnol", "spanish":
						validValues = append(validValues, "spanish", "espagnol")
					case "allemand", "german", "deutsch":
						validValues = append(validValues, "german", "deutsch", "allemand")
					case "italien", "italian", "italiano":
						validValues = append(validValues, "italian", "italiano", "italien")
					case "portugais", "portuguese":
						validValues = append(validValues, "portuguese", "portugais")
					case "japonais", "japanese":
						validValues = append(validValues, "japanese", "japonais")
					case "coréen", "korean":
						validValues = append(validValues, "korean", "coréen")
					case "chinois", "chinese":
						validValues = append(validValues, "chinese", "chinois")
					case "russe", "russian":
						validValues = append(validValues, "russian", "russe")
					case "arabe", "arabic":
						validValues = append(validValues, "arabic", "arabe")
					}
					// For subtitles: values often contain language + suffix ("Anglais (Forcés PGS)")
					// so also add translations based on Contains, not just exact match
					if isSubtitle {
						switch {
						case strings.Contains(lowerV, "anglais") || strings.Contains(lowerV, "english"):
							validValues = append(validValues, "english", "anglais", "eng")
						case strings.Contains(lowerV, "français") || strings.Contains(lowerV, "french"):
							validValues = append(validValues, "french", "français", "fr", "francais")
						case strings.Contains(lowerV, "espagnol") || strings.Contains(lowerV, "spanish"):
							validValues = append(validValues, "spanish", "espagnol")
						case strings.Contains(lowerV, "allemand") || strings.Contains(lowerV, "german"):
							validValues = append(validValues, "german", "deutsch", "allemand")
						case strings.Contains(lowerV, "italien") || strings.Contains(lowerV, "italian"):
							validValues = append(validValues, "italian", "italien")
						}
					}
				}
			}
		}

		if len(validValues) == 0 {
			continue
		}

		for _, tag := range group.Tags {
			isMatch := false
			normTag := normalize(tag.Name)
			// Also normalize the tag slug for matching
			normSlug := normalize(tag.Slug)

			for _, val := range validValues {
				normVal := normalize(val)

				if normTag == normVal || normSlug == normVal {
					isMatch = true
				} else if !strictMatch {
					if strings.Contains(normTag, normVal) || strings.Contains(normSlug, normVal) {
						isMatch = true
					} else if strings.Contains(normVal, normTag) {
						isMatch = true
					} else if strings.Contains(strings.ToLower(tag.Name), strings.ToLower(val)) {
						isMatch = true
					}
				}

				if isMatch {
					break
				}
			}

			if isMatch {
				addTag(tag)
				tagsFoundForGroup = true
			}
		}

		// Fallback "Autre"
		// For subtitles: also activate fallback additively if any original subtitle value
		// has no corresponding non-fallback tag (e.g. Spanish subtitles when only FR/ENG tags exist)
		needsFallback := !tagsFoundForGroup
		if !needsFallback && isSubtitle {
			for _, origVal := range valuesToCheck {
				if origVal == "" {
					continue
				}
				lo := strings.ToLower(origVal)
				// Build check values including translations for this original value
				checkVals := []string{normalize(lo)}
				switch {
				case strings.Contains(lo, "anglais") || strings.Contains(lo, "english"):
					checkVals = append(checkVals, "english", "eng", "anglais")
				case strings.Contains(lo, "français") || strings.Contains(lo, "french"):
					checkVals = append(checkVals, "french", "fr", "francais")
				case strings.Contains(lo, "espagnol") || strings.Contains(lo, "spanish"):
					checkVals = append(checkVals, "spanish", "espagnol")
				case strings.Contains(lo, "allemand") || strings.Contains(lo, "german"):
					checkVals = append(checkVals, "german", "deutsch", "allemand")
				case strings.Contains(lo, "italien") || strings.Contains(lo, "italian"):
					checkVals = append(checkVals, "italian", "italien")
				case strings.Contains(lo, "vff"):
					checkVals = append(checkVals, "vff")
				case strings.Contains(lo, "vfq"):
					checkVals = append(checkVals, "vfq")
				}
				covered := false
				for _, tag := range group.Tags {
					if strings.HasPrefix(strings.ToLower(tag.Name), "autre") {
						continue
					}
					normTagName := normalize(tag.Name)
					for _, cv := range checkVals {
						if normTagName == cv || strings.Contains(cv, normTagName) || strings.Contains(normTagName, cv) {
							covered = true
							break
						}
					}
					if covered {
						break
					}
				}
				if !covered {
					needsFallback = true
					break
				}
			}
		}
		if needsFallback {
			for _, tag := range group.Tags {
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
