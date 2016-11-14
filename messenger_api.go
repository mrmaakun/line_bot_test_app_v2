package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Id        string  `json:"id,omitempty"`
	Type      string  `json:"type,omitempty"`
	Text      string  `json:"text,omitempty"`
	PackageId string  `json:"packageId,omitempty"`
	StickerId string  `json:"stickerId,omitempty"`
	Title     string  `json:"title,omitempty"`
	Address   string  `json:"address,omitempty"`
	Latitude  float32 `json:"latitude,omitempty"`
	Longitude float32 `json:"longitude,omitempty"`
}

type ReplyMessage struct {
	Type               string            `json:"type,omitempty"`
	Text               string            `json:"text,omitempty"`
	OriginalContentUrl string            `json:"originalContentUrl,omitempty"`
	PreviewImageUrl    string            `json:"previewImageUrl,omitempty"`
	PackageId          string            `json:"packageId,omitempty"`
	StickerId          string            `json:"stickerId,omitempty"`
	Duration           string            `json:"duration,omitempty"`
	Title              string            `json:"title,omitempty"`
	Address            string            `json:"address,omitempty"`
	Latitude           float32           `json:"latitude,omitempty"`
	Longitude          float32           `json:"longitude,omitempty"`
	BaseUrl            string            `json:"baseUrl,omitempty"`
	AltText            string            `json:"altText,omitempty"`
	BaseSize           ImagemapBaseSize  `json:"baseSize,omitempty"`
	Actions            []ImagemapActions `json:"actions,omitempty"`
	Template           Template          `json:"template,omitempty"`
}

func ReplyToMessage(replyToken string, m Message) error {

	// Make Reply API Request

	switch m.Type {

	case "text":

		replyMessage := ReplyMessage{
			Text: m.Text,
			Type: m.Type,
		}

		err := SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

		if err != nil {
			return err
		}

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

		err := SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

		if err != nil {
			return err
		}
	case "video":

		videoPath := GetContent(m.Type, m.Id)
		video_url := "https://line-bot-test-app-v2.herokuapp.com/images/" + videoPath
		preview_image_url := "https://line-bot-test-app-v2.herokuapp.com/images/video_thumbnail.jpg"

		replyMessage := ReplyMessage{

			Type:               m.Type,
			OriginalContentUrl: video_url,
			PreviewImageUrl:    preview_image_url,
		}

		err := SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

		if err != nil {
			return err
		}
	case "audio":

		audioPath := GetContent(m.Type, m.Id)
		audio_url := "https://line-bot-test-app-v2.herokuapp.com/images/" + audioPath

		replyMessage := ReplyMessage{

			Type:               m.Type,
			OriginalContentUrl: audio_url,
			Duration:           "240000",
		}

		err := SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

		if err != nil {
			return err
		}
	case "sticker":

		replyMessage := ReplyMessage{
			Type:      m.Type,
			PackageId: m.PackageId,
			StickerId: m.StickerId,
		}

		log.Println("PackageId: " + m.PackageId)
		log.Println("Stickerid: " + m.StickerId)

		err := SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

		if err != nil {
			return err
		}
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

		err := SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

		if err != nil {
			return err
		}

	}

	return nil

}

func CheckMAC(message, messageMAC, key []byte) bool {

	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func APIPathHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("This is the Default Path Handler")
	log.Println("Entered the default Path Handler")

	// Verify Request Signature

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read the response body: "+err.Error(), http.StatusInternalServerError)
	}

	if os.Getenv("SKIP_SIGNATURE_VERIFICATION") != "TRUE" {

		decoded_signature, err := base64.StdEncoding.DecodeString(r.Header.Get("X-Line-Signature"))

		if err != nil {
			http.Error(w, "Failed to read the response body: "+err.Error(), http.StatusInternalServerError)

		}

		channel_secret := os.Getenv("LINE_CHANNEL_SECRET")

		mac := hmac.New(sha256.New, []byte(channel_secret))
		mac.Write(body)
		mac.Sum(nil)

		if CheckMAC(body, decoded_signature, []byte(channel_secret)) == false {

			log.Println("ERROR: Message Verification Has Failed.")
			http.Error(w, "ERROR: Message Authentication Failed", http.StatusUnauthorized)
			return
		} else {

			log.Println("Message Verification Has Succeeded")

		}

	} else {

		log.Println("Bot is set to bypass Signature Verification")
	}

	request := &struct {
		Events []*Event `json:"events"`
	}{}

	err = json.Unmarshal(body, &request)

	if err != nil {
		panic(err)
	}

	for _, event := range request.Events {
		log.Println("replytoken: " + event.ReplyToken)
		log.Println("type: " + event.Type)
		log.Println("timestamp: ", event.Timestamp)

		switch event.Type {
		case "message":
			err = ProcessMessageEvent(*event)
		case "follow":
			err = ProcessFollowEvent(*event)
		case "unfollow":
			ProcessUnfollowEvent(*event)
		case "join":
			err = ProcessJoinEvent(*event)
		case "leave":
			ProcessLeaveEvent(*event)
		case "postback":
			err = ProcessPostbackEvent(*event)
		default:
			log.Println("Caught invalid event type!")
			err = &APIError{
				Code:     500,
				Response: "Caught invalid event type: " + event.Type,
			}
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
