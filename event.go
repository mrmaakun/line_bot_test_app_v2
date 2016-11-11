package main

import (
	"encoding/json"
	"log"
	"strings"
)

type Event struct {
	ReplyToken string          `json:"replyToken"`
	Type       string          `json:"type"`
	Timestamp  int64           `json:"timestamp"`
	Source     Source          `json:"source"`
	Message    json.RawMessage `json:"message"`
	Postback   json.RawMessage `json:"postback"`
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

	// Image Map
	if strings.Contains(strings.ToLower(m.Text), "imagemap") {

		SendImageMap(e.ReplyToken)

	} else {

		ReplyToMessage(e.ReplyToken, m)
	}

	// Leave API
	if strings.Contains(strings.ToLower(m.Text), "goodbye") {

		switch e.Source.Type {

		case "room":

			LeaveGroupOrRoom(e.Source.Type, e.Source.RoomId)

		case "group":

			LeaveGroupOrRoom(e.Source.Type, e.Source.GroupId)

		}
	}

	// Carousel API
	if strings.Contains(strings.ToLower(m.Text), "carousel") {

		// Implement Carousel API

	}

	// Confirm Dialog API
	if strings.Contains(strings.ToLower(m.Text), "confirm dialogue") {

		// Implement Confirm Dialogue

	}

	// Buttons Dialog API
	if strings.Contains(strings.ToLower(m.Text), "buttons dialogue") {

		// Implement Buttons Dialogue

	}

}
