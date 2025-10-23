package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"random-discogs-item/models"

	"github.com/joho/godotenv"
)

// Reuse a single HTTP client and cache token/username to avoid repeated work
var (
	httpClient  = &http.Client{Timeout: 15 * time.Second}
	cachedToken string
	cachedUser  string
)

func addAuth(req *http.Request) {
	token := getToken()
	if token != "" {
		req.Header.Set("Authorization", "Discogs token="+token)
	}
}

func getToken() string {
	if cachedToken != "" {
		return cachedToken
	}

	// Load env if necessary
	_ = godotenv.Load()
	cachedToken = os.Getenv("DISCOGS_TOKEN")
	return cachedToken
}

func getUsername() string {
	if cachedUser != "" {
		return cachedUser
	}

	req, err := http.NewRequest("GET", "https://api.discogs.com/oauth/identity", nil)
	if err != nil {
		log.Fatal(err)
	}
	addAuth(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Fatalf("unexpected status from /oauth/identity: %s", resp.Status)
	}

	// Use Decoder directly to stream and avoid double-marshalling
	var data struct {
		Username string `json:"username"`
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		log.Fatal("failed to decode oauth/identity response:", err)
	}
	if data.Username == "" {
		log.Fatal("username not found in response")
	}
	cachedUser = data.Username
	return cachedUser
}

func getFolders() []models.CollectionFolder {
	username := getUsername()
	url := fmt.Sprintf("https://api.discogs.com/users/%s/collection/folders", username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	addAuth(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Fatalf("unexpected status from folders endpoint: %s", resp.Status)
	}

	var data struct {
		Folders []models.CollectionFolder `json:"folders"`
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		log.Fatal("failed to decode folders response:", err)
	}
	return data.Folders
}

func getFolderItems(folderID int, folderName string) []models.Record {
	username := getUsername()
	records := []models.Record{}

	for page := 1; ; page++ {
		url := fmt.Sprintf("https://api.discogs.com/users/%s/collection/folders/%d/releases?page=%d", username, folderID, page)
		if debug {
			fmt.Println("Fetching URL:", url)
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		addAuth(req)

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		if resp.Body == nil {
			resp.Body.Close()
			break
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			log.Fatalf("unexpected status fetching folder items: %s", resp.Status)
		}

		var data struct {
			Releases   []models.Record `json:"releases"`
			Pagination struct {
				Pages int `json:"pages"`
			} `json:"pagination"`
		}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&data); err != nil {
			resp.Body.Close()
			log.Fatal("failed to decode folder items response:", err)
		}
		resp.Body.Close()

		for _, r := range data.Releases {
			r.FolderName = folderName
			records = append(records, r)
		}

		if data.Pagination.Pages <= page || data.Pagination.Pages == 0 {
			break
		}
	}
	return records
}

func getLengthOfCollection() int {
	username := getUsername()
	url := fmt.Sprintf("https://api.discogs.com/users/%s/collection/folders/0", username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	addAuth(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Fatalf("unexpected status from collection length endpoint: %s", resp.Status)
	}

	var data struct {
		Count int `json:"count"`
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		log.Fatal("failed to decode collection length response:", err)
	}
	if data.Count == 0 {
		log.Fatal("count not found in response")
	}
	return data.Count
}

func getRecordsFromCache() []models.Record {
	cacheFile := cacheDir() + "records_cache.json"
	file, err := os.Open(cacheFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var records []models.Record
	dec := json.NewDecoder(file)
	if err := dec.Decode(&records); err != nil {
		log.Fatal(err)
	}
	return records
}
