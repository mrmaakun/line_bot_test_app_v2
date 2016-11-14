package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Source struct {
	Type    string `json:"type,omitempty"`
	UserId  string `json:"userid,omitempty"`
	GroupId string `json:"groupId,omitempty"`
	RoomId  string `json:"roomId,omitempty"`
}

type Postback struct {
	Data string `json:"data,omitempty"`
}

type Event struct {
	ReplyToken string          `json:"replyToken,omitempty"`
	Type       string          `json:"type,omitempty"`
	Timestamp  int64           `json:"timestamp,omitempty"`
	Source     Source          `json:"source,omitempty"`
	Message    json.RawMessage `json:"message,omitempty"`
	Postback   Postback        `json:"postback,omitempty"`
}

// Function that handles postback events
func ProcessPostbackEvent(e Event) error {

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

			err := SendReplyMessage(e.ReplyToken, []ReplyMessage{replyMessage1, replyMessage2})

			if err != nil {
				return err
			}

			return nil

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

			err := SendReplyMessage(e.ReplyToken, []ReplyMessage{replyMessage1, replyMessage2})

			if err != nil {
				return err
			}

			return nil

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

		err := SendReplyMessage(e.ReplyToken, []ReplyMessage{replyMessage1, replyMessage2})

		if err != nil {
			return err
		}

		return nil

	}

	return nil
}

// Function to handle follow events
func ProcessFollowEvent(e Event) error {

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

	err := SendReplyMessage(e.ReplyToken, []ReplyMessage{replyMessage1, replyMessage2, replyMessage3})

	if err != nil {
		return err
	}

	return nil

}

// Function to handle follow events
func ProcessJoinEvent(e Event) error {

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

	err := SendReplyMessage(e.ReplyToken, []ReplyMessage{replyMessage1, replyMessage2, replyMessage3})

	if err != nil {
		return err
	}

	return nil

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
func ProcessMessageEvent(e Event) error {

	var m Message

	log.Println("Entered ProcessMessageEvent")

	err := json.Unmarshal(e.Message, &m)

	log.Println("Finished Unmarshall")

	if err != nil {
		log.Fatalln("error unmarshalling message: ", err)
	}

	// Image Map
	if strings.Contains(strings.ToLower(m.Text), "imagemap") {

		err := SendImageMap(e.ReplyToken)

		if err != nil {
			return err
		}

		return nil

	}

	// Leave API
	if strings.Contains(strings.ToLower(m.Text), "goodbye") {

		var err error

		switch e.Source.Type {

		case "room":

			err = LeaveGroupOrRoom(e.Source.Type, e.Source.RoomId)

		case "group":

			err = LeaveGroupOrRoom(e.Source.Type, e.Source.GroupId)

		default:

			err = &APIError{
				Code:     500,
				Response: "Invalid Source Type",
			}

		}

		if err != nil {
			return err
		} else {
			return nil
		}
	}

	// Carousel API
	if strings.Contains(strings.ToLower(m.Text), "carousel") {

		// Implement Carousel API

	}

	// Confirm Dialog API
	if strings.Contains(strings.ToLower(m.Text), "explode") {

		log.Println("Processing explode event")

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

		err := SendReplyMessage(e.ReplyToken, []ReplyMessage{confirmMessage})

		if err != nil {
			return err
		}

		return nil

	}

	// Buttons Dialog API
	if strings.Contains(strings.ToLower(m.Text), "find zombie") {

		log.Println("Processing zombie event")

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

		err := SendReplyMessage(e.ReplyToken, []ReplyMessage{buttonMessage})

		if err != nil {
			return err
		}

		return nil

	}

	// Carousel Dialog API
	if strings.Contains(strings.ToLower(m.Text), "multizombie") {

		log.Println("Processing Multizombie Event")

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

		column1 := Column{
			ThumbnailImageUrl: "https://line-bot-test-app-v2.herokuapp.com/images/static/zombiemessage.jpg",
			Title:             "Zombie 1",
			Text:              "You have encoutered Zombie 1!",
			Actions:           templateActions,
		}

		column2 := Column{
			ThumbnailImageUrl: "https://line-bot-test-app-v2.herokuapp.com/images/static/zombiemessage.jpg",
			Title:             "Zombie 2",
			Text:              "You have encoutered Zombie 2!",
			Actions:           templateActions,
		}

		column3 := Column{
			ThumbnailImageUrl: "https://line-bot-test-app-v2.herokuapp.com/images/static/zombiemessage.jpg",
			Title:             "Zombie 3",
			Text:              "You have encoutered Zombie 3!",
			Actions:           templateActions,
		}

		//Declare Columns Array
		columns := []Column{column1, column2, column3}

		template := Template{
			Type:              "carousel",
			ThumbnailImageUrl: "https://line-bot-test-app-v2.herokuapp.com/images/static/zombiemessage.jpg",
			Title:             "You have encountered a ZOMBIE!!",
			Text:              "What do you do?!?",
			Actions:           templateActions,
			Columns:           columns,
		}

		carouselMessage := ReplyMessage{
			AltText:  "This is a Carousel template",
			Type:     "template",
			Template: template,
		}

		err := SendReplyMessage(e.ReplyToken, []ReplyMessage{carouselMessage})

		if err != nil {
			return err
		}

		return nil

	}

	//	ReplyToMessage(e.ReplyToken, m)

	return nil

}
