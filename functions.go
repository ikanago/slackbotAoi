package slackbotAoi

import (
	"encoding/json"
	"fmt"
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

	isAuthor, _ := regexp.MatchString("70_pocky", tweet.UserName)
	isReTweet, _ := regexp.MatchString("RT", tweet.Text)
	isDesiredTweet, _ := regexp.MatchString("創作2コマ漫画", tweet.Text)
	if isAuthor && !isReTweet && isDesiredTweet {
		client := slack.New(verificationToken)
		_, _, err = client.PostMessage("CBHMC8KF0", slack.MsgOptionText(tweet.LinkToTweet, false))
		if err != nil {
			log.Printf("error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, "success")
	}
}

func HelloCommand(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(verificationToken) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/hello":
		params := &slack.Msg{ResponseType: "in_channel", Text: "Hello, <@" + s.UserID + ">"}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
