package main

import (
	"encoding/json"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

type Source struct {
	Type    string `json:"type"`
	UserId  string `json:"userid"`
	GroupId string `json:"groupId"`
	RoomId  string `json:"roomId"`
}

type Message struct {
	Id        string  `json:"id"`
	Type      string  `json:"type"`
	Text      string  `json:"text"`
	PackageId string  `json:"packageId"`
	StickerId string  `json:"stickerId"`
	Title     string  `json:"title"`
	Address   string  `json:"address"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type ReplyMessage struct {
	Type               string            `json:"type"`
	Text               string            `json:"text"`
	OriginalContentUrl string            `json:"originalContentUrl"`
	PreviewImageUrl    string            `json:"previewImageUrl"`
	PackageId          string            `json:"packageId"`
	StickerId          string            `json:"stickerId"`
	Duration           string            `json:"duration"`
	Title              string            `json:"title"`
	Address            string            `json:"address"`
	Latitude           float32           `json:"latitude"`
	Longitude          float32           `json:"longitude"`
	BaseUrl            string            `json:"baseUrl"`
	AltText            string            `json:"altText"`
	BaseSize           ImagemapBaseSize  `json:"baseSize"`
	Actions            []ImagemapActions `json:"actions"`
	Template           Template          `json:"template"`
}

// Create a preview image from the original image
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

func ReplyToMessage(replyToken string, m Message) {

	// Make Reply API Request

	switch m.Type {

	case "text":

		replyMessage := ReplyMessage{
			Text: m.Text,
			Type: m.Type,
		}

		SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

	case "image":

		// TODO: Put this url in config file
		imagePath := GetContent(m.Type, m.Id)
		image_url := "https://line-bot-test-app-v2.herokuapp.com/images/" + imagePath
		preview_image_url := "https://line-bot-test-app-v2.herokuapp.com/images/" + CreatePreviewImage(imagePath)

		replyMessage := ReplyMessage{
			Type:               m.Type,
			OriginalContentUrl: image_url,
			PreviewImageUrl:    preview_image_url,
		}

		SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

	case "video":

		videoPath := GetContent(m.Type, m.Id)
		video_url := "https://line-bot-test-app-v2.herokuapp.com/images/" + videoPath
		preview_image_url := "https://line-bot-test-app-v2.herokuapp.com/images/video_thumbnail.jpg"

		replyMessage := ReplyMessage{

			Type:               m.Type,
			OriginalContentUrl: video_url,
			PreviewImageUrl:    preview_image_url,
		}

		SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

	case "audio":

		audioPath := GetContent(m.Type, m.Id)
		audio_url := "https://line-bot-test-app-v2.herokuapp.com/images/" + audioPath

		replyMessage := ReplyMessage{

			Type:               m.Type,
			OriginalContentUrl: audio_url,
			Duration:           "240000",
		}

		SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

	case "sticker":

		replyMessage := ReplyMessage{
			Type:      m.Type,
			PackageId: m.PackageId,
			StickerId: m.StickerId,
		}

		log.Println("PackageId: " + m.PackageId)
		log.Println("Stickerid: " + m.StickerId)

		SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

	case "location":

		replyMessage := ReplyMessage{
			Type:      m.Type,
			Title:     m.Title,
			Address:   m.Address,
			Latitude:  m.Latitude,
			Longitude: m.Longitude,
		}

		log.Println("Message Type: " + m.Type)
		log.Println("Title: " + m.Title)
		log.Println("Address: " + m.Address)
		log.Println("Latitude: ", m.Latitude)
		log.Println("Longitude: ", m.Longitude)

		SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

	}

}

func APIPathHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("This is the Default Path Handler")
	log.Println("Entered the default Path Handler")

	//Convert io.ReadCloser to String

	//	buf := new(bytes.Buffer)
	//	buf.ReadFrom(r.Body)
	//	requestString := buf.String()

	//	log.Println("Request Body: \n" + requestString)

	decoder := json.NewDecoder(r.Body)

	request := &struct {
		Events []*Event `json:"events"`
	}{}

	err := decoder.Decode(&request)

	if err != nil {
		panic(err)
	}

	for _, event := range request.Events {
		log.Println("replytoken: " + event.ReplyToken)
		log.Println("type: " + event.Type)
		log.Println("timestamp: ", event.Timestamp)

		switch event.Type {
		case "message":
			ProcessMessageEvent(*event)
		case "follow":
			ProcessFollowEvent(*event)
		case "unfollow":
			ProcessUnfollowEvent(*event)
		case "join":
			ProcessJoinEvent(*event)
		case "leave":
			ProcessLeaveEvent(*event)
		default:
		}
	}

}

func registerRouteHandlers() {

	log.Println("Registering Route Handlers")

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	http.HandleFunc("/api/", APIPathHandler)

	var endpoint_port string
	// If port is set an the environment variables, use that
	if endpoint_port = os.Getenv("PORT"); endpoint_port == "" {

		// Default endpoint is 12345
		endpoint_port = "12345"
		log.Println("setting port")

	}

	log.Println("Listening on port " + endpoint_port)
	log.Fatal(http.ListenAndServe(":"+endpoint_port, nil)) //nil means the default Router Server is used

}

func main() {

	log.Println("V2 Test Bot Started")

	registerRouteHandlers()

	log.Println("Registered Route Handlers")

}
