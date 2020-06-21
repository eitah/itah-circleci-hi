package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
	// lambda.Start(handler)
	eliTesting()
}

type postBody struct {
	text string
}

func eliTesting() (string, error) {
	spew.Dump("eli is the best!")
	// curl -X POST -H 'Content-type: application/json' --data '{"text":"Hello, World!"}'
	var buf bytes.Buffer
	p := postBody{text: "Hello, Pipeline!"}

	err := json.NewEncoder(&buf).Encode(p)
	if err != nil {
		return "", err
	}
	spew.Dump("eprepost")

	response, err := http.Post(config.SlackWebhookURL, "application/json", &buf)
	if err != nil {
		spew.Dump("err", err)

		return "", err
	}
	defer response.Body.Close()
	spew.Dump("post")

	return "success", nil
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (string, error) {
	spew.Dump("eli is the best!")
	// curl -X POST -H 'Content-type: application/json' --data '{"text":"Hello, World!"}'
	var buf bytes.Buffer
	p := postBody{text: "Hello, Pipeline!"}

	err := json.NewEncoder(&buf).Encode(p)
	if err != nil {
		return "", err
	}
	response, err := http.Post(config.SlackWebhookURL, "application/json", &buf)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	spew.Dump(err)
	return "success", nil

	// response, err := http.Post("")
}
