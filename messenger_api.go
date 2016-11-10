package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Source struct {
	Type    string `json:"type"`
	UserId  string `json:"userid"`
	GroupId string `json:"groupId"`
	RoomId  string `json:"roomId"`
}

type Profile struct {
	DisplayName   string `json:"displayName"`
	UserId        string `json:"userId"`
	PictureUrl    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

type Event struct {
	ReplyToken string          `json:"replyToken"`
	Type       string          `json:"type"`
	Timestamp  int64           `json:"timestamp"`
	Source     Source          `json:"source"`
	Message    json.RawMessage `json:"message"`
	Postback   json.RawMessage `json:"postback"`
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

type ImagemapArea struct {
	X      int32 `json:"x"`
	Y      int32 `json:"y"`
	Width  int32 `json:"width"`
	Height int32 `json:"height"`
}

type ImagemapActions struct {
	Type    string       `json:"type"`
	Text    string       `json:"text"`
	LinkUri string       `json:"linkUrl"`
	Area    ImagemapArea `json:"area"`
}

type ImagemapBaseSize struct {
	Height int32 `json:"height"`
	Width  int32 `json:"width"`
}

type Reply struct {
	SendReplyToken string         `json:"replyToken"`
	Messages       []ReplyMessage `json:"messages"`
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
}

// This function checks to see if the number of files in the images directory is less than the max number.
// If it is, it deletes the oldest image

func CleanImageDirectory() {

	//Get a slice of files in the images directory
	files, _ := ioutil.ReadDir("images")

	numberOfStoredImages := len(files)

	// TODO: Change the max number of stored images to a config item
	if numberOfStoredImages > 30 {

		var earliestModifiedTime time.Time
		var earliestModifiedFileName string

		for _, f := range files {

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

func GetProfile(userId string) Profile {

	client := &http.Client{}

	url := "https://api.line-beta.me/v2/bot/profile/" + userId

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(resp.Body)

	var userProfile Profile

	err = decoder.Decode(&userProfile)

	if err != nil {
		panic(err)
	}

	return userProfile

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

func SendImageMap(replyToken string) {

	zone1 := ImagemapActions{
		Type:    "uri",
		LinkUri: "http://www.explodingkittens.com/",
		Area:    ImagemapArea{X: 47, Y: 54, Width: 293, Height: 528},
	}

	/*
		zone2 := ImagemapActions{
			Type: "message",
			Text: "ZOMBIES!!",
			Area: ImagemapArea{X: 549, Y: 49, Width: 293, Height: 528},
		}
	*/
	replyMessage := ReplyMessage{

		Type:     "imagemap",
		BaseUrl:  "https://line-bot-test-app-v2.herokuapp.com/images/imagemap",
		AltText:  "This is an imagemap",
		BaseSize: ImagemapBaseSize{Height: 636, Width: 1040},
		Actions:  []ImagemapActions{zone1},
	}

	SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

}

func SendReplyMessage(replyToken string, replyMessages []ReplyMessage) {

	url := "https://api.line-beta.me/v2/bot/message/reply"

	var jsonPayload []byte = nil
	var err error

	reply := Reply{
		SendReplyToken: replyToken,
		Messages:       replyMessages,
	}

	jsonPayload, err = json.Marshal(reply)

	//Make reply message

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	s := buf.String()

	log.Println("Json Req:" + s)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	log.Println("Response Status:", resp.Status)
	log.Println("Response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("Response Body:", string(body))

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

// Function to handle follow events
func ProcessFollowEvent(e Event) {

	log.Println("Processing Follow Event")

	replyMessage1 := ReplyMessage{
		Text: "Hi, " + GetProfile(e.Source.GroupId).DisplayName + "!!",
		Type: "text",
	}

	replyMessage2 := ReplyMessage{
		Text: "Thank you for being my friend!",
		Type: "text",
	}

	replyMessage3 := ReplyMessage{
		Type:      "sticker",
		StickerId: "144",
		PackageId: "2",
	}

	SendReplyMessage(e.ReplyToken, []ReplyMessage{replyMessage1, replyMessage2, replyMessage3})

}

// Function to handle follow events
func ProcessJoinEvent(e Event) {

	log.Println("Processing Join Event")

	replyMessage1 := ReplyMessage{
		Text: "Hello everybody!",
		Type: "text",
	}

	replyMessage2 := ReplyMessage{
		Text: "Thank you for inviting me to this group!",
		Type: "text",
	}

	replyMessage3 := ReplyMessage{
		Type:      "sticker",
		StickerId: "144",
		PackageId: "2",
	}

	SendReplyMessage(e.ReplyToken, []ReplyMessage{replyMessage1, replyMessage2, replyMessage3})

}

// Function to handle follow events
func ProcessUnfollowEvent(e Event) {

	log.Println("Bot has been unfollowed by user: " + e.Source.UserId)

}

// Function to handle follow events
func ProcessLeaveEvent(e Event) {

	log.Println("Bot has left group: " + e.Source.GroupId)

}

// Function to handle all message events
func ProcessMessageEvent(e Event) {

	var m Message

	log.Println("Entered ProcessMessageEvent")

	err := json.Unmarshal(e.Message, &m)

	log.Println("Finished Unmarshall")

	if err != nil {
		log.Fatalln("error unmarshalling message: ", err)
	}

	if m.Text == "Imagemap" {

		SendImageMap(e.ReplyToken)

	} else {

		ReplyToMessage(e.ReplyToken, m)
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
