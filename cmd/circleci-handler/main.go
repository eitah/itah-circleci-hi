package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/davecgh/go-spew/spew"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config contains all the configurable env variables absorbed from the environment. split_words turns camel in to snake case.
type Config struct {
	SlackWebhookURL string `split_words:"true" envconfig:"SLACK_WEBHOOK_URL"`

	// logging
	LogLevel string `split_words:"true" default:"info"`
	LogJSON  bool   `split_words:"true" default:"true"`
}

var (
	log    *logrus.Logger
	config Config
)

func init() {
	if err := envconfig.Process("", &config); err != nil {
		fmt.Fprintf(os.Stderr, "invalid environment config: %s\n", err)
		os.Exit(1)
	}

	log = logrus.New()
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid log level %s: %s\n", config.LogLevel, err)
		os.Exit(1)
	}
	log.SetLevel(level)

	if config.LogJSON {
		log.SetFormatter(&logrus.JSONFormatter{})
	}
}

func main() {
	// todo I want header auth to happen in the gateway not in the lambda, so I bypass using that config value for now
	//lambda.Start(handler)
	eliTesting()
}

func eliTesting() {
	ctx := context.Background()
	build, err := processBuildNotification(ctx, "{ \"attachments\": [             {               \"fallback\": \":red_circle: A $CIRCLE_JOB job has failed!\",               \"text\": \":red_circle: A $CIRCLE_JOB job has failed! $SLACK_MENTIONS\",               \"fields\": [                 {                   \"title\": \"Project\",                   \"value\": \"$CIRCLE_PROJECT_REPONAME\",                   \"short\": true                 },                 {                   \"title\": \"Job Number\",                   \"value\": \"$CIRCLE_BUILD_NUM\",                   \"short\": true                 }               ],               \"actions\": [                 {                   \"type\": \"button\",                   \"text\": \"Visit Job\",                   \"url\": \"$CIRCLE_BUILD_URL\"                 }               ],               \"color\": \"#ed5c5c\"             }           ]         }")

	sendMe, err := Debounce()
	if err != nill {
		//  assume if redis fails  we should still proceed to post
		sendMe = true
	}

	if sendMe == true {
		//postErr := postToWebHook(ctx)
		// if postErr != nil {
		// 	fmt.Println(err)
		// }
	}

	if err != nil {
		// only after the post occurs do we want to return failure
		fmt.Println(err)
	}
	spew.Dump(sendMe)

}

func handler(ctx context.Context, e events.APIGatewayProxyRequest) (string, error) {
	// parse event and get the payload details
	//body := e.Body

	//build, err := ProcessBuildNotification(ctx, e.Body);
	// check redis for details and write new details if required

	// if unique or errored, post to redis, else drop it or throw or w/e
	//err := postToWebHook(ctx)
	//if err != nil {
	//	return "", err
	//}
	return "Eli says success!", nil
}

// Debounce queries redis for the notification entry and returns sendMe
func Debounce(build BuildStatus) (bool, error) {

}

// BuildStatusNotification is an incoming circleci build status
type BuildStatusNotification struct {
	Attachments []struct {
		Fallback string `json:"fallback"`
		Fields   []struct {
			Title string `json:"title"`
			Value string `json:"value"`
		} `json:"fields"`
		Actions []struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		} `json:"actions"`
	} `json:"attachments"`
}

// BuildStatus is the outgoing flatter
type BuildStatus struct {
	ProjectReponame string
	CircleBuildNum  string
	CircleBuildURL  string
}

func processBuildNotification(_ context.Context, e string) (BuildStatus, error) {
	var raw BuildStatusNotification
	var build BuildStatus
	err := json.Unmarshal([]byte(e), &raw)
	if err != nil {
		return BuildStatus{}, err
	}

	if raw.Attachments[0].Fields[0].Title == "Project" {
		build.ProjectReponame = raw.Attachments[0].Fields[0].Value
	}

	if raw.Attachments[0].Fields[1].Title == "Job Number" {
		build.CircleBuildNum = raw.Attachments[0].Fields[1].Value
	}

	if raw.Attachments[0].Actions[0].URL != "" {
		build.CircleBuildURL = raw.Attachments[0].Actions[0].URL
	}

	return build, nil
}

func postToWebHook(ctx context.Context) error {
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	requestBody, err := json.Marshal(map[string]string{
		"text": "I am deployed!",
	})
	if err != nil {
		return err
	}

	spew.Dump("webhook url is", config.SlackWebhookURL)

	req, err := http.NewRequestWithContext(ctx, "POST", config.SlackWebhookURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Webhook failed! Status Code: %v \n", resp.StatusCode)
	}

	// this is not needed rn but i thought I might want to reference it
	//_, err := ioutil.ReadAll(resp.Body)
	//spew.Dump(string(body))

	return nil
}
