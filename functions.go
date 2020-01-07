package slackbotAoi

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var verificationToken string

func init() {
	verificationToken = os.Getenv("VERIFICATION_TOKEN")
}

type Tweet struct {
	Text        string `json:"text"`
	UserName    string `json:"userName"`
	LinkToTweet string `json:"linkToTweet"`
}

func SendTweet(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var tweet Tweet
	if err := json.Unmarshal(body, &tweet); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if tweet.isMatchTweet("70_pocky", "創作2コマ漫画") {
		err := postMessage("C9DTMB6GZ", tweet.LinkToTweet)
		if err != nil {
			log.Printf("error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if tweet.isMatchTweet("yuukikikuchi", "100日後に死ぬワニ") {
		err := postMessage("C9DTMB6GZ", tweet.LinkToTweet)
		if err != nil {
			log.Printf("error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// Check `tweet` satisfies a certain requirements.
func (tweet *Tweet) isMatchTweet(author string, content string) bool {
	isAuthor, _ := regexp.MatchString(author, tweet.UserName)
	isReTweet, _ := regexp.MatchString("RT", tweet.Text)
	isDesiredTweet, _ := regexp.MatchString(content, tweet.Text)
	return isAuthor && !isReTweet && isDesiredTweet
}

func postMessage(channel string, message string) error {
	client := slack.New(verificationToken)
	_, _, err := client.PostMessage(channel, slack.MsgOptionText(message, false))
	return err
}
