package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
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
	lambda.Start(handler)
	//eliTesting()
}

//func eliTesting() {
//	ctx := context.Background()
//	err := postToWebHook(ctx)
//	if err != nil {
//		fmt.Println(err)
//	}
//}

func handler(ctx context.Context, e events.APIGatewayProxyRequest) (string, error) {
	// parse event and get the payload details
	body := e.Body
	spew.Dump(body)
	//err := ProcessBuildNotification(ctx, body
	// check redis for details and write new details if required

	// if unique or errored, post to redis, else drop it or throw or w/e
	//err := postToWebHook(ctx)
	//if err != nil {
	//	return "", err
	//}
	return "Eli says success!", nil
}

type SlackFailurePayload struct {
	text string
}

func ProcessBuildNotification(ctx context.Context, e events.APIGatewayProxyRequest)  {

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
