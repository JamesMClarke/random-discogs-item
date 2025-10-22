package main

import (
	"fmt"
	"random-discogs-item/models"
)

func main() {
	auth := models.Auth{
		Token:    getToken(),
		Username: getUsername(),
	}
	fmt.Println("Token:", auth.Token)
	fmt.Println("Username:", auth.Username)

	records := []models.Record{}

	folders := getFolders()
	for _, folder := range folders {
		fmt.Printf("Folder ID: %d, Name: %s, Count: %d, Resource URL: %s\n",
			folder.ID, folder.Name, folder.Count, folder.ResourceURL)
		records = append(records, getFolderItems(folder.ID, folder.Name)...)
	}

	for _, record := range records {
		fmt.Printf("Record ID: %d, Title: %s, Year: %d, Format: %s, Resource URL: %s, Folder Name: %s\n",
			record.ID, record.BasicInformation.Title, record.BasicInformation.Year, record.BasicInformation.Formats, record.BasicInformation.ResourceURL, record.FolderName)
	}
}
