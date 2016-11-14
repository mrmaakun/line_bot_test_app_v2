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

type Template struct {
	Type              string           `json:"type,omitempty"`
	ThumbnailImageUrl string           `json:"thumbnailImageUrl,omitempty"`
	Title             string           `json:"menu,omitempty"`
	Text              string           `json:"text,omitempty"`
	Actions           []TemplateAction `json:"actions,omitempty"`
	Columns           []Column         `json:"columns,omitempty"`
}

type TemplateAction struct {
	Type  string `json:"type,omitempty"`
	Label string `json:"label,omitempty"`
	Data  string `json:"data,omitempty"`
	Text  string `json:"text,omitempty"`
	Uri   string `json:"uri,omitempty"`
}

type Column struct {
	ThumbnailImageUrl string           `json:"thumbnailImageUrl,omitempty"`
	Title             string           `json:"title,omitempty"`
	Text              string           `json:"text,omitempty"`
	Actions           []TemplateAction `json:"actions,omitempty"`
}

type Reply struct {
	SendReplyToken string         `json:"replyToken,omitempty"`
	Messages       []ReplyMessage `json:"messages,omitempty"`
}

type Profile struct {
	DisplayName   string `json:"displayName,omitempty"`
	UserId        string `json:"userId,omitempty"`
	PictureUrl    string `json:"pictureUrl,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
}

type ImagemapArea struct {
	X      int32 `json:"x,omitempty"`
	Y      int32 `json:"y,omitempty"`
	Width  int32 `json:"width,omitempty"`
	Height int32 `json:"height,omitempty"`
}

type ImagemapActions struct {
	Type    string       `json:"type,omitempty"`
	Text    string       `json:"text,omitempty"`
	LinkUri string       `json:"linkUri,omitempty"`
	Area    ImagemapArea `json:"area,omitempty"`
}

type ImagemapBaseSize struct {
	Height int32 `json:"height,omitempty"`
	Width  int32 `json:"width,omitempty"`
}

func SendImageMap(replyToken string) error {

	zone1 := ImagemapActions{
		Type:    "uri",
		LinkUri: "http://www.explodingkittens.com/",
		Area:    ImagemapArea{X: 47, Y: 54, Width: 293, Height: 528},
	}

	zone2 := ImagemapActions{
		Type: "message",
		Text: "ZOMBIES!!",
		Area: ImagemapArea{X: 549, Y: 49, Width: 293, Height: 528},
	}

	replyMessage := ReplyMessage{

		Type:     "imagemap",
		BaseUrl:  "https://line-bot-test-app-v2.herokuapp.com/images/imagemap",
		AltText:  "This is an imagemap",
		BaseSize: ImagemapBaseSize{Height: 636, Width: 1040},
		Actions:  []ImagemapActions{zone1, zone2},
	}

	err := SendReplyMessage(replyToken, []ReplyMessage{replyMessage})

	if err != nil {
		return err
	}

	return nil
}

func SendReplyMessage(replyToken string, replyMessages []ReplyMessage) error {

	url := "https://api.line-beta.me/v2/bot/message/reply"

	var jsonPayload []byte = nil
	var err error

	reply := Reply{
		SendReplyToken: replyToken,
		Messages:       replyMessages,
	}

	jsonPayload, err = json.Marshal(reply)

	log.Printf("SendReplyMessage(): Request JSON: " + string(jsonPayload))

	//Make reply message

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

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

	if resp.StatusCode != http.StatusOK {

		return &APIError{
			Code:     resp.StatusCode,
			Response: string(body),
		}
	}

	return nil

}

func LeaveGroupOrRoom(leaveType string, Id string) error {

	var url string

	// Set the API url based on the type of group/room that is being left
	switch leaveType {

	case "room":

		url = "https://api.line-beta.me/v2/bot/room/" + Id + "/leave"

	case "group":

		url = "https://api.line-beta.me/v2/bot/group/" + Id + "/leave"

	default:

		panic(fmt.Sprintf("%s", "Calling LeaveGroupOrRoom on invalid leaveType!"))

	}

	var jsonPayload []byte = nil
	var err error

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))

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

	if resp.StatusCode != http.StatusOK {

		return &APIError{

			Code:     resp.StatusCode,
			Response: string(body),
		}
	}

	return nil

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
