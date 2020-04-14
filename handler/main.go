package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/pkg/errors"
)

// form represents the form present on the website as filled in by a user.
type form struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func errorResponse(statusCode int, message string, err error) (events.APIGatewayProxyResponse, error) {
	err = errors.Wrap(err, message)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       err.Error(),
	}, err
}

func sendEmail(fromAddress, toAddress, subject, message string) error {
	emailClient := ses.New(session.New(), aws.NewConfig().WithRegion("eu-west-1"))

	_, err := emailClient.SendEmail(&ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{&toAddress},
		},
		Source: &fromAddress,
		Message: &ses.Message{
			Subject: &ses.Content{
				Data: &subject,
			},
			Body: &ses.Body{
				Text: &ses.Content{
					Data: &message,
				},
			},
		},
	})

	return err
}

func getEmailAddress() string {
	// Hardcoded but base64 encoded to foil (some simple) web scraper spammers
	const emailAddressBase64 = "dG9taGpwQGdtYWlsLmNvbQ=="
	bytes, err := base64.StdEncoding.DecodeString(emailAddressBase64)
	if err != nil {
		panic("failed to decode email address")
	}

	return string(bytes)
}

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	emailAddress := getEmailAddress()
	const emailSubject = "Contact form from seafordcruk.org"

	log.Printf("Handling request: %+v\n", req)

	body := req.Body
	if req.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return errorResponse(500, "failed to decode request", err)
		}
		body = string(decoded)
	}
	formValues, err := url.ParseQuery(body)
	if err != nil {
		return errorResponse(400, "failed to parse form values", err)
	}

	contactName := formValues.Get("name")
	contactEmail := formValues.Get("email")
	contactMessage := formValues.Get("message")
	emailMessage := fmt.Sprintf("Name: %s\nEmail: %s\nMessage: %s", contactName, contactEmail, contactMessage)

	// Honey pot values
	honeyPotAddress := formValues.Get("address")
	honeyPotTimeString := formValues.Get("time") // integer number of ms since the javascript loaded

	// If the hidden field for address was filled in, or time is less than 15s, drop it as spam
	if honeyPotAddress != "" {
		log.Printf("Dropping message %q as suspected spam because address field was filled in as %q (time was %q)", emailMessage, honeyPotAddress, honeyPotTimeString)
	} else if honeyPotTime, err := strconv.Atoi(honeyPotTimeString); err != nil || honeyPotTime < 15000 {
		log.Printf("Dropping message %q as suspected spam because time was invalid: %q", emailMessage, honeyPotTimeString)
	} else {
		err = sendEmail(emailAddress, emailAddress, emailSubject, emailMessage)
		if err != nil {
			return errorResponse(500, "failed to send email", err)
		}
	}

	// Even if we dropped the message as spam, still return a success response so we don't leak that information
	return events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers: map[string]string{
			"location": "https://www.seafordcruk.org/",
		},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
