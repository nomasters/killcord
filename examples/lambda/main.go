package main

import (
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nomasters/killcord"
)

func main() {
	var opts killcord.ProjectOptions
	session := killcord.New()

	otps.Payload.Secret = os.Getenv("KILLCORD_PAYLOAD_SECRET")
	opts.Payload.RPCURL = os.Getenv("KILLCORD_PAYLOAD_RPCURL")
	opts.Contract.ID = os.Getenv("KILLCORD_CONTRACT_ID")
	opts.Contract.RPCURL = os.Getenv("KILLCORD_CONTRACT_RPCURL")
	opts.Contract.Publisher.Address = os.Getenv("KILLCORD_CONTRACT_PUBLISHER_ADDRESS")
	opts.Contract.Publisher.Password = os.Getenv("KILLCORD_CONTRACT_PUBLISHER_PASSWORD")
	opts.Contract.Publisher.KeyStore = os.Getenv("KILLCORD_CONTRACT_PUBLISHER_KEYSTORE")

	// check for dev mode string and set bool
	if os.Getenv("KILLCORD_DEV_MODE") == "true" {
		opts.DevMode = true
	}

	// check for warning threshold and if it exists, convert to int64
	if x := os.Getenv("KILLCORD_PUBLISHER_WARNING_THRESHOLD"); x != "" {
		y, err := strconv.Atoi(x)
		if err == nil {
			opts.Publisher.WarningThreshold = int64(y)
		}
	}

	// check for publish threshold and if it exists, convert to int64
	if x := os.Getenv("KILLCORD_PUBLISHER_PUBLISH_THRESHOLD"); x != "" {
		y, err := strconv.Atoi(x)
		if err == nil {
			opts.Publisher.PublishThreshold = int64(y)
		}
	}

	session.Options = opts
	session.Init()

	lambda.Start(session.PublisherLambdaHandler)
}
