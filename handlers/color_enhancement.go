package handlers

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/RSO-project-Prepih/AI-service/database"
	"github.com/RSO-project-Prepih/AI-service/prometheus"
	"github.com/nfnt/resize"
)

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

func saveProcessedImage(userID, originalImageID string, processedImageData []byte) {
	db := database.NewDBConnection()
	defer db.Close()

	newImageID := originalImageID

	processingServiceType := "ColorEnhancement"

	query := "INSERT INTO imageprocessing (image_id, processing_service_type, data) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, newImageID, processingServiceType, processedImageData)
	if err != nil {
		log.Println("Error saving processed image:", err)
		return
	}
}

const maxFileSize = 8 * 1024 * 1024 // 8 MB

func checkAndResizeImage(imageData []byte, maxResolution, maxFileSize uint) ([]byte, error) {
	log.Println("Checking image size...")
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		log.Println("Error decoding image:", err)
		return nil, err
	}

	bounds := img.Bounds()
	width := uint(bounds.Dx())
	height := uint(bounds.Dy())
	log.Println("Image size:", width, "x", height)

	// Resize if resolution is too high
	if width > maxResolution || height > maxResolution {
		log.Println("Resizing image...")
		img = resize.Thumbnail(maxResolution, maxResolution, img, resize.Lanczos3)
	}

	// Encode the image to check the file size
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 75}) // Start with a default quality
	if err != nil {
		log.Println("Error encoding image:", err)
		return nil, err
	}

	// If the file size is too large, reduce quality
	for buf.Len() > int(maxFileSize) {
		log.Println("Reducing image quality...")
		buf.Reset()
		quality := 75 * int(maxFileSize) / buf.Len() // Reduce quality proportional to excess size
		if quality < 10 {
			log.Println("Error reducing image quality: quality too low")
			quality = 10 // Set a minimum quality to prevent too much degradation
		}
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
		if err != nil {
			log.Println("Error encoding image:", err)
			return nil, err
		}
	}

	log.Println("Image size after processing:", buf.Len(), "bytes")
	return buf.Bytes(), nil
}

func PostColorEnhancementPhoto(userID, imageID string) {

	starteTime := time.Now()

	url := "https://www.ailabapi.com/api/image/enhance/image-color-enhancement"
	method := "POST"

	log.Println("Fetching image data...")
	imageData, err := FetchImageData(userID, imageID)
	if err != nil {
		log.Println("Error fetching image data:", err)
		return
	}

	log.Println("Image data fetched successfully")

	log.Println("Resizing image...")
	maxResolution := uint(3000)
	imageDataResize, err := checkAndResizeImage(imageData, maxResolution, maxFileSize)
	if err != nil {
		log.Println("Error processing image:", err)
		return
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	part, err := writer.CreateFormFile("image", "filename.jpg") // Ensure filename has correct extension
	if err != nil {
		log.Println("Error creating form file:", err)
		return
	}

	_, err = part.Write(imageDataResize)
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

	duration := time.Since(starteTime)
	prometheus.HTTPRequestDuration.WithLabelValues("/uploadphoto").Observe(duration.Seconds())
}
