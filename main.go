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

// EmbeddedCredentials is the JSON service account credentials embedded directly in the code
const EmbeddedCredentials = `{
  "type": "service_account",
  "project_id": "taptap-436017",
  "private_key_id": "941670d40e7f011ad34645f7e46a791fe932d67e",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDK/zyy2w5f7OYk\nHSXEaIpk61oFyVaiBjdXjp1oW00chb3pKLeSbkViXSPTawAPFroVIe7R5e68W1dv\nNvNasHY2wQ2vpOlTcqeKBiHP4a3iUW+AmVRj8T0eCdXDXFmfSRzx1xqfdbIvyTYs\nNlrfEhXAVdCV5Fx0EWjKlpdERMt+YnN4j6rZzgwUJDuGmu3YwCmzTNjfLCdTUSH9\nSFvrzOeCDl0DUUD4MuQpi8tZ6e1L+PWUxvPZqRitSze882eXs4ALvWPs3fSLg5oZ\n7jv3/dd6UL9hz7Acf5tXwRf7zr1lkgYjYdkufVJO0IHjI9T69CBypjJVWNrhGQzX\nwTwJPCbzAgMBAAECggEAT1xwyQXN/V+a69R/XuV18ZV3WNrJZUej3DWzwUgsgfGh\neOWDuxokQxvhtTZYTA3ZCwj8mo/XgUj+ikrD4hkp5iccaCZDV+3zpQjXsDNtLRUk\nT/Th4r6945/5s8pHeXf2em/bhyrW0krKRIetiBdEbLC//tHL+U6TFty/3587pTk5\nv/Ll4At6lbW0ff4j1CRDyM2uJSnu9wBfw5FfYeNulRb1BSFvgKMbKsEey1MseNQb\nOpOTW/3M8lX+39PtT5JlqR/oI0O34nfS/B5gkVC/42v8isQFL2HHnHIx86ZATkXq\npvrVSqwGpSl/q8biFEvAZUU/65cfqWQN3pO40KxtdQKBgQD/bw2KHd9oNRTavLZE\n4hQU9HnvZttNK8U4QHWED6NTSpvkqb5EZChOwanUNgZRIHsm6a4U5LhcUFa7FgQ/\nQ+XuJu5E8E5qf9XKGj+OnwkBCvqq+OP62jZ3TdwWk/2qGjyY819I/gnSc00BghMj\nZ2yWUxGnoyDdQpRvNf8JArbKRwKBgQDLcm3Acg+JREwR8WDz0WCo/X8gvDFq0tlq\nULAkkWbtgFsDgxwXFI6stSZgoIy645nc+DHBc267CGlfJMVAV5pM9iJDeq0c58l5\niC3X7gcs2cIUUoiug/xio4mWmg62y/Vc6YjK9HGCCKTzNwesNDnnIAIoZntckb8h\nBnP+/sZn9QKBgGNcOoMQbbfmdg9EOw5+dttT4h6q/wF81kG0aUIOpzSeIBgJo1aN\nM1S3Zq2CumBSZzVSzwGXmtNl+ObbgJlvewBxqlussoQg5/Ou1CxRVrpOIAXjvSL2\nQRuVcNhjhtflTs8cVGNbVkzDxx+gDnvGHmo3M/XmscD/xiegdG133cy/AoGBAMmc\nwmq+HetYBVKatAuraHDPlhYoqYhFHzQedhAnD6s5UfhvC31L4AADHN8Q+6WRO78h\nLp2Y+RjcQyAIXnle1wiBun7IqZlFkgGgFF4yAmZN/ekJyW24WnqduhHG5eH8yVCk\nFe2axImqa1yjIjVjJCaJL9o9hO69eH0P2g/PB2upAoGABHjnbzG3VCSPxBCY48vW\nrDFp5ei3l94FqJn0B1VW7dXewAYF+R9FBZzJu2r1iE7G0DPskLU8mlw2tALQx3tj\nVAZkNPf5NCHgaimxNNWwNPZTjg9cygXmx28YEn2Jx40zZoKcZHqUkFi/uG91V0Pa\nt169ZiHYgkZEfA3+z2X+ygY=\n-----END PRIVATE KEY-----\n",
  "client_email": "dat-test@taptap-436017.iam.gserviceaccount.com",
  "client_id": "116107819533727108566",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/dat-test%40taptap-436017.iam.gserviceaccount.com",
  "universe_domain": "googleapis.com"
}`

// writeToSheet appends a row of data to the Google Sheet
func writeToSheet(data []string) error {
	ctx := context.Background()

	// Try using embedded credentials directly
	var srv *sheets.Service
	var err error

	// Use JSON credentials in memory
	srv, err = sheets.NewService(ctx, option.WithCredentialsJSON([]byte(EmbeddedCredentials)))
	if err != nil {
		log.Printf("Error using embedded credentials: %v", err)

		// Fall back to file-based credentials as a backup
		credentialsFile := getCredentialsPath()
		log.Printf("Falling back to credentials file: %s", credentialsFile)

		srv, err = sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile))
		if err != nil {
			return fmt.Errorf("unable to create sheets service: %v", err)
		}
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
