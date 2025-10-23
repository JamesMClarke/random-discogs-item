package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"random-discogs-item/models"
)

// TODO: Add error handling
// TODO: Add command line flags for folder selection, etc.
// TODO: Add command line flags for cache location, force refresh, etc.
// TODO: Get random item from cached data depending on criteria
var debug bool

func main() {
	auth := models.Auth{
		Token:    getToken(),
		Username: getUsername(),
	}
	if debug {
		fmt.Println("Token:", auth.Token)
		fmt.Println("Username:", auth.Username)
	}

	singles := flag.Bool("singles", false, "Whether to include singles in the selection")
	notShared := flag.Bool("not-shared", false, "Whether to exclude shared items")
	forceUpdate := flag.Bool("force-update", false, "Whether to force update the cache")
	debugFlag := flag.Bool("debug", false, "Enable debug mode")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [who] [options]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Positional arguments:")
		fmt.Fprintln(os.Stderr, "  who\tWho's folder to get the item from (choices: alice, james, both)")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults() // prints all the flag definitions automatically
	}

	// Parse command-line arguments
	flag.Parse()

	// Handle positional arguments
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [who] [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "who must be one of: alice, james, both\n")
		os.Exit(1)
	}

	who := args[0]
	validWho := map[string]bool{"alice": true, "james": true, "both": true}
	if !validWho[who] {
		fmt.Fprintf(os.Stderr, "Invalid argument for who: %s (must be one of alice, james, both)\n", who)
		os.Exit(1)
	}

	// Set package-level debug and print results
	debug = *debugFlag
	if debug {
		fmt.Printf("who: %s\n", who)
		fmt.Printf("singles: %v\n", *singles)
		fmt.Printf("notShared: %v\n", *notShared)
		fmt.Printf("debug: %v\n", debug)
		fmt.Printf("forceUpdate: %v\n", *forceUpdate)
	}

	// Get length of collection
	collectionLength := getLengthOfCollection()
	cacheLength := getCacheLength()

	if cacheLength != collectionLength || *forceUpdate {
		fmt.Printf("Cached records is a different length (%d vs %d), updating cache...\n", cacheLength, collectionLength)
		updateCache()
	} else if debug {
		fmt.Println("Cache is up to date, no need to update.")
	}

	allRecords := getRecordsFromCache()
	if debug {
		fmt.Printf("Total records in cache: %d\n", len(allRecords))
	}
	filteredRecords := []models.Record{}
	switch who {
	case "alice":
		filteredRecords = append(filteredRecords, filterRecordsByFolder(allRecords, "Alice LPs")...)
		if *singles {
			filteredRecords = append(filteredRecords, filterRecordsByFolder(allRecords, "Alice Singles")...)
		}
	case "james":
		filteredRecords = append(filteredRecords, filterRecordsByFolder(allRecords, "James LPs")...)
		if *singles {
			filteredRecords = append(filteredRecords, filterRecordsByFolder(allRecords, "James Singles")...)
		}
	case "both":
		filteredRecords = append(filteredRecords, filterRecordsByFolder(allRecords, "James LPs")...)
		if *singles {
			filteredRecords = append(filteredRecords, filterRecordsByFolder(allRecords, "James Singles")...)
		}
		filteredRecords = append(filteredRecords, filterRecordsByFolder(allRecords, "Alice LPs")...)
		if *singles {
			filteredRecords = append(filteredRecords, filterRecordsByFolder(allRecords, "Alice Singles")...)
		}
	}
	if !*notShared {
		filteredRecords = append(filteredRecords, filterRecordsByFolder(allRecords, "Shared LPs")...)
		if *singles {
			filteredRecords = append(filteredRecords, filterRecordsByFolder(allRecords, "Shared Singles")...)
		}
	}

	if debug {
		fmt.Printf("Total records after filtering: %d\n", len(filteredRecords))
		for _, record := range filteredRecords {
			displayRecord(record)
		}
	}

	item := getRandomItem(filteredRecords)
	displayRecord(item)
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

func filterRecordsByFolder(records []models.Record, folderName string) []models.Record {
	filtered := []models.Record{}
	for _, record := range records {
		if record.FolderName == folderName {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func displayRecord(record models.Record) {
	fmt.Printf("Record ID: %d\n", record.ID)
	fmt.Printf("Title: %s\n", record.BasicInformation.Title)
	fmt.Printf("Year: %d\n", record.BasicInformation.Year)
	fmt.Printf("Format: ")
	for _, format := range record.BasicInformation.Formats {
		fmt.Printf("%s ", format.Name)
	}
	fmt.Printf("\nResource URL: %s\n", record.BasicInformation.ResourceURL)
	fmt.Printf("Folder Name: %s\n", record.FolderName)
}

func getRandomItem(records []models.Record) models.Record {
	if len(records) == 0 {
		log.Fatal("No records available to select from.")
	}
	randIndex := rand.Intn(len(records))
	return records[randIndex]
}
