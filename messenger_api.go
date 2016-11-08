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
	"time"
)

type Event struct {
	ReplyToken string          `json:"replyToken"`
	Type       string          `json:"type"`
	Timestamp  int64           `json:"timestamp"`
	Source     json.RawMessage `json:"source"`
	Message    json.RawMessage `json:"message"`
	Postback   json.RawMessage `json:"postback"`
}

type Message struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Text      string `json:"text"`
	PackageId string `json:"packageId"`
	StickerId string `json:"stickerId"`
}

type Reply struct {
	SendReplyToken string         `json:"replyToken"`
	Messages       []ReplyMessage `json:"messages"`
}

type ReplyMessage struct {
	Type               string `json:"type"`
	Text               string `json:"text"`
	OriginalContentUrl string `json:"originalContentUrl"`
	PreviewImageUrl    string `json:"previewImageUrl"`
	PackageId          string `json:"packageId"`
	StickerId          string `json:"stickerId"`
}

// Function for downloading and temporarily storing images, sound, and videos
// Returns the file name of the stored image
func GetContent(mediaType string, mediaId string) string {

	client := &http.Client{}
	rand.Seed((time.Now().UTC().UnixNano()))

	imageFileName := "image_" + string(rand.Intn(10000)) + ".jpg"
	// Create output file
	newFile, err := os.Create(imageFileName)

	url := "https://api.line.me/v2/bot/message/" + mediaId + "/content"

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

	previewImageFileName := "p" + originalFileName

	previewImageFile, err := os.Create(previewImageFileName)

	//Resize image
	resizedImage := resize.Resize(240, 240, image, resize.Lanczos3)

	var image_reader io.ReadWriter

	jpeg.Encode(image_reader, resizedImage, nil)

	io.Copy(previewImageFile, image_reader)

	return previewImageFileName

}

func SendReplyMessage(replyToken string, m Message) {

	// Make Reply API Request
	url := "https://api.line-beta.me/v2/bot/message/reply"

	var jsonPayload []byte = nil
	var err error

	switch m.Type {

	case "text":

		replyMessage := ReplyMessage{
			Text: m.Text,
			Type: m.Type,
		}

		reply := Reply{
			SendReplyToken: replyToken,
			Messages:       []ReplyMessage{replyMessage},
		}
		jsonPayload, err = json.Marshal(reply)

	case "image":

		image_url := GetContent(m.Type, m.Id)
		preview_image_url := CreatePreviewImage(image_url)

		replyMessage := ReplyMessage{
			Type:               m.Type,
			OriginalContentUrl: image_url,
			PreviewImageUrl:    preview_image_url,
		}

		reply := Reply{
			SendReplyToken: replyToken,
			Messages:       []ReplyMessage{replyMessage},
		}
		jsonPayload, err = json.Marshal(reply)

	case "sticker":

		replyMessage := ReplyMessage{
			Type:      m.Type,
			PackageId: m.PackageId,
			StickerId: m.StickerId,
		}

		reply := Reply{
			SendReplyToken: replyToken,
			Messages:       []ReplyMessage{replyMessage},
		}
		jsonPayload, err = json.Marshal(reply)

	}

	//Make reply message

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	//reqbody, _ := ioutil.ReadAll(req.Body)
	//log.Println("Reply Message Response Body:", string(reqbody))

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

// Function to handle all message events
func ProcessMessageEvent(e Event) {

	var m Message

	log.Println("Entered ProcessMessageEvent")

	err := json.Unmarshal(e.Message, &m)

	log.Println("Finished Unmarshall")

	if err != nil {
		log.Fatalln("error unmarshalling message: ", err)
	}

	SendReplyMessage(e.ReplyToken, m)

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
