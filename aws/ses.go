package aws

import (
	"context"
	"net/mail"

	email "github.com/antonybholmes/go-mail"
	"github.com/antonybholmes/go-sys/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sestypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type (
	SESMailer struct {
		svc  *sesv2.Client
		from *string
	}
)

func NewSESMailer(awsConfig aws.Config, from *mail.Address) (email.Mailer, error) {

	// Use the provided AWS config to create a new SES client
	svc := sesv2.NewFromConfig(awsConfig)

	return &SESMailer{
		svc:  svc,
		from: aws.String(from.String()),
	}, nil
}

func (mailer *SESMailer) SendHtmlMail(ctx context.Context,
	to *mail.Address,
	subject string,
	html string) error {
	input := &sesv2.SendEmailInput{
		Content: &sestypes.EmailContent{
			Simple: &sestypes.Message{
				Body: &sestypes.Body{
					Html: &sestypes.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(html),
					},
					// Text: &types.Content{
					// 	Charset: aws.String("UTF-8"),
					// 	Data:    aws.String("Hello, world!\nThis is a test email sent with AWS SES V2."),
					// },
				},
				Subject: &sestypes.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(subject),
				},
			},
		},
		Destination: &sestypes.Destination{
			ToAddresses: []string{to.Address}, // Replace with the recipient's email address
		},
		FromEmailAddress: mailer.from, // Replace with your verified sender email address
	}

	// Send the email
	resp, err := mailer.svc.SendEmail(ctx, input)
	if err != nil {
		return err
	}

	// Print the message ID on success
	log.Debug().Msgf("email sent: %s", *resp.MessageId)

	return nil
}
