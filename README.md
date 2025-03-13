# Google Sheets API Integration with Golang

This project provides a simple API built with Golang that allows you to write data to Google Sheets. It uses the Google Sheets API v4 and the Gin web framework.

## Prerequisites

- Go 1.19+ installed
- A Google Cloud Platform account
- A Google Sheets spreadsheet

## Setup Instructions

### 1. Configure Google Cloud Project

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project
3. Enable the Google Sheets API:
   - Navigate to "API & Services" > "Library"
   - Search for "Google Sheets API" and enable it
4. Create a Service Account:
   - Go to "IAM & Admin" > "Service Accounts"
   - Click "Create Service Account"
   - Enter a name and description
   - Grant the "Editor" role for Google Sheets
   - Click "Done"
5. Create and download credentials:
   - Click on the newly created service account
   - Go to the "Keys" tab
   - Click "Add Key" > "Create new key"
   - Choose JSON format and create
   - Save the downloaded JSON file as `credentials.json` in the `credentials` directory of this project

### 2. Share Your Google Sheet

1. Create a new Google Sheet or use an existing one
2. Copy the Spreadsheet ID from the URL:
   - The URL will look like: `https://docs.google.com/spreadsheets/d/YOUR_SPREADSHEET_ID/edit`
3. Share your Google Sheet with the service account email (found in the credentials.json file)
   - The email will look like: `your-service-account@your-project.iam.gserviceaccount.com`
   - Give it "Editor" access

### 3. Configure the Application

1. Open `main.go` and update the following constants:
   ```go
   const (
       spreadsheetID = "YOUR_SPREADSHEET_ID" // Replace with your actual spreadsheet ID
       sheetName     = "Sheet1"              // Replace with your sheet name if different
   )
   ```

### 4. Run the Application

```bash
# Move the credentials file to the credentials directory
mv path/to/downloaded/credentials.json credentials/

# Update permissions
chmod 600 credentials/credentials.json

# Run the application
go run main.go
```

The server will start on port 8080.

## API Usage

### Submit Data to Google Sheets

**Endpoint:** `POST /submit`

**Request Body:**

```json
{
  "email": "user@example.com",
  "phone_number": "123456789",
  "name": "John Doe",
  "telegram_id": "@johndoe",
  "score": "95"
}
```

**Success Response:**

```json
{
  "message": "Data submitted successfully"
}
```

**Error Response:**

```json
{
  "error": "Failed to write to sheet"
}
```

## Testing the API

You can test the API using curl:

```bash
curl -X POST "http://localhost:8080/submit" \
-H "Content-Type: application/json" \
-d '{
    "email": "test@example.com",
    "phone_number": "123456789",
    "name": "John Doe",
    "telegram_id": "@johndoe",
    "score": "100"
}'
```

## Security Considerations

- Keep your `credentials.json` file secure and never commit it to your repository
- Consider implementing API authentication if this service will be publicly accessible
- Use HTTPS in production environments

## Troubleshooting

- If you encounter a "Permission denied" error, make sure you've shared your Google Sheet with the service account email
- Check that the spreadsheet ID and sheet name are correct
- Verify that the Google Sheets API is enabled in your Google Cloud project

## License

MIT
