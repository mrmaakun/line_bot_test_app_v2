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
	Type              string           `json:"type"`
	ThumbnailImageUrl string           `json:"thumbnailImageUrl"`
	Title             string           `json:"menu"`
	Text              string           `json:"text"`
	Actions           []TemplateAction `json:"actions"`
}

type TemplateAction struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	Data  string `json:"data"`
	Text  string `json:"text"`
	Uri   string `json:"uri"`
}

type Column struct {
	ThumbnailImageUrl string           `json:"thumbnailImageUrl"`
	Title             string           `json:"title"`
	Text              string           `json:"text"`
	Actions           []TemplateAction `json:"actions"`
}

type Reply struct {
	SendReplyToken string         `json:"replyToken"`
	Messages       []ReplyMessage `json:"messages"`
}

type Profile struct {
	DisplayName   string `json:"displayName"`
	UserId        string `json:"userId"`
	PictureUrl    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
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
	LinkUri string       `json:"linkUri"`
	Area    ImagemapArea `json:"area"`
}

type ImagemapBaseSize struct {
	Height int32 `json:"height"`
	Width  int32 `json:"width"`
}

func SendImageMap(replyToken string) {

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

func LeaveGroupOrRoom(leaveType string, Id string) {

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
