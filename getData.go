package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

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
