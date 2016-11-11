package main

import (
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Create a preview image from the original image
// TODO: Make this method work for the static images too
func CreatePreviewImage(originalFileName string) string {

	// Open File
	file, err := os.Open("images/" + originalFileName)
	if err != nil {
		log.Fatal(err)
	}

	//Read Image
	image, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Image Read")

	previewImageFileName := "p_" + originalFileName

	previewImageFile, err := os.Create("images/" + previewImageFileName)

	//Resize image
	resizedImage := resize.Resize(240, 240, image, resize.Lanczos3)

	jpeg.Encode(previewImageFile, resizedImage, nil)

	return previewImageFileName

}

// This function checks to see if the number of files in the images directory is less than the max number.
// If it is, it deletes the oldest image

func CleanImageDirectory() {

	//Get a slice of files in the images directory
	files, _ := ioutil.ReadDir("images")

	//Debug statement

	numberOfStoredImages := len(files)

	// TODO: Change the max number of stored images to a config item
	if numberOfStoredImages > 30 {

		var earliestModifiedTime time.Time
		var earliestModifiedFileName string

		for _, f := range files {

			// Ignore file if it is a directory
			if f.IsDir() == true {
				continue
			}

			// If this is the first element, set it as the earliest one
			if earliestModifiedFileName == "" {

				earliestModifiedTime = f.ModTime()
				earliestModifiedFileName = f.Name()
				continue
			}

			if earliestModifiedTime.Before(f.ModTime()) {

				earliestModifiedTime = f.ModTime()
				earliestModifiedFileName = f.Name()
			}
		}

		err := os.Remove(earliestModifiedFileName)
		if err != nil {
			log.Fatal(err)

		}

	}

}

// Function for downloading and temporarily storing images, sound, and videos
// Returns the file name of the stored image
func GetContent(mediaType string, mediaId string) string {

	client := &http.Client{}
	rand.Seed((time.Now().UTC().UnixNano()))
	url := "https://api.line-beta.me/v2/bot/message/" + mediaId + "/content"

	switch mediaType {

	case "image":
		// Clean the image directory before getting content
		CleanImageDirectory()

		imageFileName := "image_" + strconv.Itoa(rand.Intn(10000)) + ".jpg"
		// Create output file
		newFile, err := os.Create("images/" + imageFileName)

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
		resp, err := client.Do(req)

		numBytesWritten, err := io.Copy(newFile, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Media ID: " + mediaId)
		log.Printf("Downloaded %d byte file.\n", numBytesWritten)
		log.Println("File name: " + imageFileName)

		//return the file name
		return imageFileName

	case "video":

		CleanImageDirectory()

		videoFileName := "video_" + strconv.Itoa(rand.Intn(10000)) + ".mp4"
		newFile, err := os.Create("images/" + videoFileName)

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
		resp, err := client.Do(req)

		numBytesWritten, err := io.Copy(newFile, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Media ID: " + mediaId)
		log.Printf("Downloaded %d byte file.\n", numBytesWritten)
		log.Println("File name: " + videoFileName)

		return videoFileName

	case "audio":

		CleanImageDirectory()

		audioFileName := "audio_" + strconv.Itoa(rand.Intn(10000)) + ".m4a"
		newFile, err := os.Create("images/" + audioFileName)

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
		resp, err := client.Do(req)

		numBytesWritten, err := io.Copy(newFile, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Media ID: " + mediaId)
		log.Printf("Downloaded %d byte file.\n", numBytesWritten)
		log.Println("File name: " + audioFileName)

		return audioFileName

	default:

		log.Println("Unknown media type")

		return ""

	}

}
