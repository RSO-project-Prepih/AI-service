package handlers

import (
	"encoding/json"
	"log"

	"github.com/RSO-project-Prepih/AI-service/database"
	_ "github.com/RSO-project-Prepih/AI-service/docs"
)

// ImageData is a struct to parse the data field.
type ImageData struct {
	Data struct {
		ImageURL string `json:"image_url"`
	} `json:"data"`
}

// ImageProcessing represents the structure of your image processing data.
type ImageProcessing struct {
	ImageID  string `json:"image_id"`
	ImageURL string `json:"image_url"`
}

// GetImageProcessingPhotos godoc
// @Summary Retrieve processed images
// @Description Fetches data for images that have undergone processing
// @Tags images
// @Accept  json
// @Produce  json
// @Success 200 {array} ImageProcessing
// @Router /image-processing [get]
func GetImageProcessingPhotos() ([]ImageProcessing, error) {
	log.Println("Fetching the images from processing data...")
	db := database.NewDBConnection()
	defer db.Close()

	var processedImages []ImageProcessing

	query := "SELECT image_id, data FROM imageprocessing"
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Error fetching image processing data:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var data []byte
		var imageData ImageData

		err := rows.Scan(&id, &data)
		if err != nil {
			log.Println("Error scanning image processing data:", err)
			return nil, err
		}

		err = json.Unmarshal(data, &imageData)
		if err != nil {
			log.Println("Error unmarshalling image processing data:", err)
			continue // Skip this row and move to the next one
		}

		processedImage := ImageProcessing{
			ImageID:  id,
			ImageURL: imageData.Data.ImageURL,
		}

		processedImages = append(processedImages, processedImage)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over image processing rows:", err)
		return nil, err
	}

	return processedImages, nil
}
