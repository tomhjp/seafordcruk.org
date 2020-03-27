package main

import (
	"context"
	"encoding/base64"

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
	const emailMessage = "Test email message"

	err := sendEmail(emailAddress, emailAddress, emailSubject, emailMessage)
	if err != nil {
		return errorResponse(500, "failed to send email", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "Hello from my first AWS Lambda",
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
