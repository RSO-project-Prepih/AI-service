package handlers

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/RSO-project-Prepih/AI-service/database"
)

// FetchImageData gets image data from the database
func FetchImageData(userID, imageID string) ([]byte, error) {
	db := database.NewDBConnection()
	defer db.Close()

	var imageData []byte
	query := "SELECT data FROM images WHERE user_id = $1 AND image_id = $2"
	row := db.QueryRow(query, userID, imageID)
	err := row.Scan(&imageData)
	if err != nil {
		log.Println("Error fetching image data:", err)
		return nil, err
	}
	return imageData, nil
}

// saveProcessedImage saves the processed image data in the imageprocessing table
func saveProcessedImage(userID, originalImageID string, processedImageData []byte) {
	db := database.NewDBConnection()
	defer db.Close()

	newImageID := originalImageID
	// Define the processing service type
	processingServiceType := "ColorEnhancement"

	// Insert the processed image data into the imageprocessing table
	query := "INSERT INTO imageprocessing (image_id, processing_service_type, data) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, newImageID, processingServiceType, processedImageData)
	if err != nil {
		log.Println("Error saving processed image:", err)
		return
	}
}

func PostColorEnhancementPhoto(userID, imageID string) {
	url := "https://www.ailabapi.com/api/image/enhance/image-color-enhancement"
	method := "POST"

	imageData, err := FetchImageData(userID, imageID)
	if err != nil {
		log.Println("Error fetching image data:", err)
		return
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	part, err := writer.CreateFormFile("image", "filename.jpg") // Ensure filename has correct extension
	if err != nil {
		log.Println("Error creating form file:", err)
		return
	}

	_, err = part.Write(imageData)
	if err != nil {
		log.Println("Error writing image data to form:", err)
		return
	}

	// Add output format field
	_ = writer.WriteField("output_format", "jpg") // Adjust the format as needed

	_ = writer.WriteField("mode", "LogC")

	err = writer.Close()
	if err != nil {
		log.Println("Error closing writer:", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	apiKey := os.Getenv("AILAB_API_KEY") // Ensure this is set in your environment
	req.Header.Add("ailabapi-api-key", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close()

	processedImageData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	log.Println("Processed image data:", string(processedImageData))

	saveProcessedImage(userID, imageID, processedImageData)
}
