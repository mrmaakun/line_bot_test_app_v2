package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	Id   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

type Reply struct {
	SendReplyToken string         `json:"replyToken"`
	Messages       []ReplyMessage `json:"messages"`
}

type ReplyMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func SendReplyMessage(replyToken string, message string) {

	// Make Reply API Request
	url := "https://api.line-beta.me/v2/bot/message/reply"

	//Make reply message

	replyMessage := ReplyMessage{
		Type: "text",
		Text: message,
	}

	reply := Reply{
		SendReplyToken: replyToken,
		Messages:       []ReplyMessage{replyMessage},
	}
	jsonPayload, err := json.Marshal(reply)
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

	log.Println("Type: " + m.Type)
	log.Println("ID: " + m.Id)

	switch m.Type {
	case "text":
		log.Println("This is a text message")
		log.Println("The message sent was: " + m.Text)

		SendReplyMessage(e.ReplyToken, m.Text)

	default:
		log.Println("Invalid Message Type")
	}

}

func DefaultPathHandler(w http.ResponseWriter, r *http.Request) {

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

	http.HandleFunc("/", DefaultPathHandler)

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
