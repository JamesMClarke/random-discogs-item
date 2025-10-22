package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"random-discogs-item/models"
)

// TODO: Add error handling
// TODO: Add logic to only update cache if data has changed
// TODO: Add command line flags for folder selection, etc.
// TODO: Add command line flags for cache location, force refresh, etc.
// TODO: Get random item from cached data depending on criteria
var debug = true

func main() {
	auth := models.Auth{
		Token:    getToken(),
		Username: getUsername(),
	}
	if debug {
		fmt.Println("Token:", auth.Token)
		fmt.Println("Username:", auth.Username)
	}

	// Get length of collection
	collectionLength := getLengthOfCollection()
	cacheLength := getCacheLength()

	if cacheLength != collectionLength {
		fmt.Printf("Cached records is a different length (%d vs %d), updating cache...\n", cacheLength, collectionLength)
		updateCache()
	}
}

func cacheDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return home + "/.cache/random-discogs-item/"
}

func checkCacheDir() {
	dir := cacheDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getCacheLength() int {
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
	return len(records)
}

func updateCache() {
	records := []models.Record{}

	folders := getFolders()
	for _, folder := range folders {
		if folder.Name == "All" {
			continue
		}
		fmt.Printf("Folder ID: %d, Name: %s, Count: %d, Resource URL: %s\n",
			folder.ID, folder.Name, folder.Count, folder.ResourceURL)
		records = append(records, getFolderItems(folder.ID, folder.Name)...)
	}

	// for _, record := range records {
	// 	fmt.Printf("Record ID: %d, Title: %s, Year: %d, Format: %s, Resource URL: %s, Folder Name: %s\n",
	// 		record.ID, record.BasicInformation.Title, record.BasicInformation.Year, record.BasicInformation.Formats, record.BasicInformation.ResourceURL, record.FolderName)
	// }

	// Check and create cache directory
	checkCacheDir()

	// Save records to cache file
	cacheFile := cacheDir() + "records_cache.json"
	file, err := os.Create(cacheFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(records)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Records cached to", cacheFile)
}
