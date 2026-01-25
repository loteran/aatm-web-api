package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
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
func (a *App) UploadToLaCale(torrentPath string, nfoPath string, title string, description string, tmdbId string, mediaType string, releaseInfo ReleaseInfo, passkey string, email string, password string) error {
	if passkey == "" {
		return fmt.Errorf("passkey is missing in settings (required for metadata)")
	}
	if email == "" || password == "" {
		return fmt.Errorf("email and password are required for upload authentication")
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

	// 4. Authenticate (Get Session)
	client, err := a.LaCaleLogin(email, password)
	if err != nil {
		return fmt.Errorf("La Cale Login failed: %w", err)
	}

	// 5. Upload (Internal API)
	internalURL := "https://la-cale.space/api/internal"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Fields
	// writer.WriteField("passkey", passkey) // Removed as we use session
	writer.WriteField("title", title)
	writer.WriteField("description", description)
	writer.WriteField("categoryId", categoryId)
	writer.WriteField("isAnonymous", "false")
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

	req, err := http.NewRequest("POST", internalURL+"/torrents/upload", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

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

type LoginResponse struct {
	Success bool `json:"success"`
}

func (a *App) LaCaleLogin(email, password string) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: jar}

	// 1. Auth Login
	authBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	resp, err := client.Post("https://la-cale.space/api/internal/auth/login", "application/json", bytes.NewBuffer(authBody))
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(b))
	}

	var res LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode login response: %w", err)
	}

	if !res.Success {
		return nil, fmt.Errorf("login was unsuccessful (success: false)")
	}

	return client, nil
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
			valuesToCheck = []string{info.Audio}
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
