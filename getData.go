package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"random-discogs-item/models"

	"github.com/joho/godotenv"
)

func addAuth(req *http.Request) {
	req.Header.Add("Authorization", "Discogs token="+getToken())
}

func getToken() string {
	// Get token from env
	err := godotenv.Load() // loads from .env in current directory
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("DISCOGS_TOKEN")
}

func getUsername() string {
	// Get username from Discogs API using token
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.discogs.com/oauth/identity", nil)
	if err != nil {
		log.Fatal(err)
	}
	addAuth(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse JSON and return the username
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	if username, ok := data["username"].(string); ok {
		return username
	}

	log.Fatal("username not found in response")
	return ""
}

func getFolders() []models.CollectionFolder {
	// Get folders from Discogs API using token and username
	client := &http.Client{}
	username := getUsername()
	url := "https://api.discogs.com/users/" + username + "/collection/folders"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	addAuth(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse JSON and return the folders
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	folders := []models.CollectionFolder{}
	if f, ok := data["folders"].([]interface{}); ok {
		for _, folder := range f {
			folderBytes, err := json.Marshal(folder)
			if err != nil {
				log.Fatal(err)
			}
			var cf models.CollectionFolder
			if err := json.Unmarshal(folderBytes, &cf); err != nil {
				log.Fatal(err)
			}
			folders = append(folders, cf)
		}
	}
	return folders
}

func getFolderItems(folderID int, folderName string) []models.Record {
	// Get collection items from Discogs API using token and username
	// TODO: Handle pagination
	client := &http.Client{}
	username := getUsername()
	records := []models.Record{}
	for page := 1; ; page++ {
		url := "https://api.discogs.com/users/" + username + "/collection/folders/" +
			fmt.Sprintf("%d", folderID) + "/releases?page=" + fmt.Sprintf("%d", page)
		if debug {
			fmt.Println("Fetching URL:", url)
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		addAuth(req)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Parse JSON and return the collection items
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			log.Fatal(err)
		}

		if r, ok := data["releases"].([]interface{}); ok {
			for _, item := range r {
				itemBytes, err := json.Marshal(item)
				if err != nil {
					log.Fatal(err)
				}
				var record models.Record
				if err := json.Unmarshal(itemBytes, &record); err != nil {
					log.Fatal(err)
				}
				record.FolderName = folderName
				records = append(records, record)
			}
		}

		// Check for pagination
		if pagination, ok := data["pagination"].(map[string]interface{}); ok {
			if pages, ok := pagination["pages"].(float64); ok {
				if debug {
					fmt.Println("Folder:", folderName, "Page:", page, "of", int(pages))
				}
				if int(pages) <= page {
					break
				}
			}
		}
	}
	return records

}

func getLengthOfCollection() int {
	// Get length of collection from Discogs API using token and username
	client := &http.Client{}
	username := getUsername()
	url := "https://api.discogs.com/users/" + username + "/collection/folders/0"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	addAuth(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse JSON and return the collection length
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}
	if count, ok := data["count"].(float64); ok {
		return int(count)
	}

	log.Fatal("count not found in response")
	return 0
}

func getRecordsFromCache() []models.Record {
	cacheFile := cacheDir() + "records_cache.json"
	file, err := os.Open(cacheFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var records []models.Record
	err = json.NewDecoder(file).Decode(&records)
	if err != nil {
		log.Fatal(err)
	}
	return records
}
