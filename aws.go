package mail

import (
	"context"
	"encoding/json"
	"net/mail"

	"github.com/antonybholmes/go-sys/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sestypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type (
	SesMailer struct {
		svc  *sesv2.Client
		from *string
	}

	SQSEmailQueue struct {
		client   *sqs.Client
		queueUrl *string
		ctx      context.Context
	}
)

var (
	AwsTypeString = aws.String("String")
	AwsTypeNumber = aws.String("Number")
	AwsTypeBinary = aws.String("Binary")
	AwsEmailValue = aws.String("email")
)

func NewSesMailer(from *mail.Address) Mailer {
	// Create a new SES session
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal().Msgf("unable to load SDK config, %v", err)
	}

	// Create a new SES client
	svc := sesv2.NewFromConfig(cfg)

	return &SesMailer{
		svc:  svc,
		from: aws.String(from.String()),
	}
}

// func (mailer *SMTPMailer) SetUser(user string, password string) *SMTPMailer {
// 	mailer.user = user
// 	mailer.password = password
// 	return mailer
// }

// func (mailer *SMTPMailer) SetHost(host string, port uint) *SMTPMailer {
// 	mailer.host = host
// 	mailer.port = port
// 	mailer.addr = fmt.Sprintf("%s:%d", mailer.host, mailer.port)
// 	return mailer
// }

// func (mailer *SMTPMailer) SetFrom(from *mail.Address) *SMTPMailer {
// 	mailer.from = from
// 	return mailer
// }

func (mailer *SesMailer) SendHtmlMail(ctx context.Context,
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

func NewSqsEmailQueue(queueUrl string) *SQSEmailQueue {
	ctx := context.Background()

	// Load AWS config (region and credentials)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal().Msgf("unable to load SDK config: %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	return &SQSEmailQueue{client: client, queueUrl: aws.String(queueUrl), ctx: ctx}
}

// SendMail sends an email message to the SQS queue.
// The message is serialized as JSON and includes message attributes.
func (queue *SQSEmailQueue) SendMail(mail *MailItem) error {

	body, err := json.Marshal(mail)

	if err != nil {
		return err
	}

	// Send to SQS
	_, err = queue.client.SendMessage(queue.ctx, &sqs.SendMessageInput{
		QueueUrl:    queue.queueUrl,
		MessageBody: aws.String(string(body)),
		// Optional: add message attributes
		MessageAttributes: map[string]sqstypes.MessageAttributeValue{
			"type": {
				DataType:    AwsTypeString,
				StringValue: AwsEmailValue,
			},
		},
	})

	return err
}
