package internal

import (
	"net/http"
	"net/url"
	"fmt"
	"os"
	"time"
	"strings"
)

const TwilioUrl = "https://api.twilio.com/2010-04-01/Accounts/"
const MessageEndPoint = "Messages.json"

func ExampleTwilio() {
	body := "Hello world"
	SendMessageTwilio(body, "+32496952214")
}

func SendMessageTwilio(content string, target string) {
	sid := os.Getenv("TWILIO_ACCOUNT_SID")
	token := os.Getenv("TWILIO_AUTH_TOKEN")
	TwilioSource := os.Getenv("TWILIO_SOURCE_NUMBER")

	endpoint := getMessageEndpoint()

	body := url.Values{}
	body.Set("To", "whatsapp:" + target)
	body.Set("From", TwilioSource)
	body.Set("Body", content)

	request, err := http.NewRequest("POST", endpoint, strings.NewReader(body.Encode()))
	if err != nil {
		panic(err)
	}
	request.SetBasicAuth(sid, token)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{
		Timeout: time.Duration(1 * time.Minute),
	}

	resp, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println(resp.Body)
}

func getMessageEndpoint() string {
	sid := os.Getenv("TWILIO_ACCOUNT_SID")
	url := TwilioUrl + sid + "/" + MessageEndPoint
	fmt.Println("URL is", url)
	return url
}