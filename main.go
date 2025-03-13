package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Constants for Google Sheets API
const (
	// Replace with your actual spreadsheet ID
	spreadsheetID = "1qNaI1ezqHf-V7DDG8A445oas1tLuARAKXsnqrMsxX3w"
	sheetName     = "Log1"
)

// getCredentialsPath returns the path to the credentials file
// considering both local development and containerized environments
func getCredentialsPath() string {
	// Check if we're running in the Docker container
	if _, err := os.Stat("/root/credentials/credentials.json"); err == nil {
		return "/root/credentials/credentials.json"
	}

	// Default to local development path
	return "credentials/credentials.json"
}

// writeToSheet appends a row of data to the Google Sheet
func writeToSheet(data []string) error {
	ctx := context.Background()

	credentialsFile := getCredentialsPath()
	log.Printf("Using credentials file: %s", credentialsFile)

	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return fmt.Errorf("unable to create sheets service: %v", err)
	}

	// Define the range for appending data
	rangeData := fmt.Sprintf("%s!A:E", sheetName)

	// Create a ValueRange containing the data to be written
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{{
			data[0], data[1], data[2], data[3], data[4],
		}},
	}

	// Append the data to the spreadsheet
	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, rangeData, valueRange).
		ValueInputOption("RAW").
		Do()

	if err != nil {
		return fmt.Errorf("unable to append data to sheet: %v", err)
	}

	return nil
}

// submitHandler handles POST requests to /submit endpoint
func submitHandler(c *gin.Context) {
	var requestData struct {
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Name        string `json:"name"`
		TelegramID  string `json:"telegram_id"`
		Score       string `json:"score"`
	}

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate required fields
	if requestData.Email == "" || requestData.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and name are required fields"})
		return
	}

	// Prepare data for Google Sheets
	data := []string{
		requestData.Email,
		requestData.PhoneNumber,
		requestData.Name,
		requestData.TelegramID,
		requestData.Score,
	}

	// Write data to Google Sheets
	if err := writeToSheet(data); err != nil {
		log.Printf("Error writing to sheet: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write to sheet"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Data submitted successfully"})
}

func main() {
	// Create a new Gin router with default middleware
	r := gin.Default()

	// Register the POST endpoint for submitting data
	r.POST("/submit", submitHandler)

	// Start the server
	log.Println("Server running on port 8000")
	if err := r.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
