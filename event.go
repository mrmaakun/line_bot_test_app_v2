package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Source struct {
	Type    string `json:"type"`
	UserId  string `json:"userid"`
	GroupId string `json:"groupId"`
	RoomId  string `json:"roomId"`
}

type Postback struct {
	Data string `json:"data"`
}

type Event struct {
	ReplyToken string          `json:"replyToken"`
	Type       string          `json:"type"`
	Timestamp  int64           `json:"timestamp"`
	Source     Source          `json:"source"`
	Message    json.RawMessage `json:"message"`
	Postback   Postback        `json:"postback"`
}

// Function that handles postback events
func ProcessPostbackEvent(e Event) {

	log.Println("Processing Postback Event")
	log.Println("Postback Data: " + e.Postback.Data)

	switch e.Postback.Data {

	case "run":

		rand.Seed((time.Now().UTC().UnixNano()))

		coinFlip := rand.Intn(100000)

		log.Println("Coin flip number: ", coinFlip)

		if coinFlip%2 == 0 {

			replyMessage1 := ReplyMessage{
				Text: "I got your run postback... and your were able to escape!!",
				Type: "text",
			}

			// TODO: Put this url in config file
			image_url := "https://line-bot-test-app-v2.herokuapp.com/images/static/run.jpg"
			preview_image_url := "https://line-bot-test-app-v2.herokuapp.com/images/static/p_run.jpg"

			replyMessage2 := ReplyMessage{
				Type:               "image",
				OriginalContentUrl: image_url,
				PreviewImageUrl:    preview_image_url,
			}

			SendReplyMessage(e.ReplyToken, []ReplyMessage{replyMessage1, replyMessage2})

		} else {

			replyMessage1 := ReplyMessage{
				Text: "I got your run postback... and the zombie got you! Now you must EXPLODE!",
				Type: "text",
			}

			// TODO: Put this url in config file
			image_url := "https://line-bot-test-app-v2.herokuapp.com/images/static/explode.jpg"
			preview_image_url := "https://line-bot-test-app-v2.herokuapp.com/images/static/p_explode.jpg"

			replyMessage2 := ReplyMessage{
				Type:               "image",
				OriginalContentUrl: image_url,
				PreviewImageUrl:    preview_image_url,
			}

			SendReplyMessage(e.ReplyToken, []ReplyMessage{replyMessage1, replyMessage2})

		}

	case "noexplode":

		replyMessage1 := ReplyMessage{
			Text: "I got a postback saying that you do not want to explode... and I think you are a coward!",
			Type: "text",
		}

		replyMessage2 := ReplyMessage{
			Type:      "sticker",
			StickerId: "527",
			PackageId: "2",
		}

		SendReplyMessage(e.ReplyToken, []ReplyMessage{replyMessage1, replyMessage2})

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

	// Image Map
	if strings.Contains(strings.ToLower(m.Text), "imagemap") {

		SendImageMap(e.ReplyToken)
		return

	}

	// Leave API
	if strings.Contains(strings.ToLower(m.Text), "goodbye") {

		switch e.Source.Type {

		case "room":

			LeaveGroupOrRoom(e.Source.Type, e.Source.RoomId)

		case "group":

			LeaveGroupOrRoom(e.Source.Type, e.Source.GroupId)

		}

		return
	}

	// Carousel API
	if strings.Contains(strings.ToLower(m.Text), "carousel") {

		// Implement Carousel API

	}

	// Confirm Dialog API
	if strings.Contains(strings.ToLower(m.Text), "I want to explode") {

		templateAction1 := TemplateAction{
			Type:  "uri",
			Label: "YES!",
			Uri:   "https://line-bot-test-app-v2.herokuapp.com/images/static/explode.jpg",
		}

		templateAction2 := TemplateAction{
			Type:  "postback",
			Label: "NO!",
			Data:  "noexplode",
		}

		templateActions := []TemplateAction{templateAction1, templateAction2}

		template := Template{
			Type:    "confirm",
			Text:    "Are you SURE you want to explode?",
			Actions: templateActions,
		}

		confirmMessage := ReplyMessage{
			AltText:  "This is a confirm template",
			Type:     "template",
			Template: template,
		}

		SendReplyMessage(e.ReplyToken, []ReplyMessage{confirmMessage})
		return

	}

	// Buttons Dialog API
	if strings.Contains(strings.ToLower(m.Text), "find zombie") {

		templateAction1 := TemplateAction{
			Type:  "postback",
			Label: "Run!",
			Data:  "run",
			Text:  "I'm outta here!!",
		}

		templateAction2 := TemplateAction{
			Type:  "message",
			Label: "Scream!",
			Text:  "AHHHHHH!",
		}

		templateAction3 := TemplateAction{
			Type:  "uri",
			Label: "EXPLODE!",
			Uri:   "https://line-bot-test-app-v2.herokuapp.com/images/static/explode.jpg",
		}

		templateActions := []TemplateAction{templateAction1, templateAction2, templateAction3}

		template := Template{
			Type:              "buttons",
			ThumbnailImageUrl: "https://line-bot-test-app-v2.herokuapp.com/images/static/zombiemessage.jpg",
			Title:             "You have encountered a ZOMBIE!!",
			Text:              "What do you do?!?",
			Actions:           templateActions,
		}

		buttonMessage := ReplyMessage{
			AltText:  "This is a buttons template",
			Type:     "template",
			Template: template,
		}

		SendReplyMessage(e.ReplyToken, []ReplyMessage{buttonMessage})
		return

	}

	//	ReplyToMessage(e.ReplyToken, m)

}
