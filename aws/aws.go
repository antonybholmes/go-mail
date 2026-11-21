package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"strings"

	email "github.com/antonybholmes/go-mail"
	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-sys/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sestypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/panjf2000/ants"
)

type (
	SQSSender struct {
		Client   *sqs.Client
		QueueURL string
	}

	SQSEmailConsumer struct {
		QueueURL     string
		Threads      int
		Client       *sqs.Client
		EmailHandler email.MailQueue
		pool         *ants.Pool
		backoff      *sys.Backoff
	}

	SESMailer struct {
		svc  *sesv2.Client
		from *string
	}

	SQSEmailQueue struct {
		client *SQSSender
	}
)

const (
	WaitTimeSecs = 10
)

var (
	AwsTypeString = aws.String("String")
	AwsTypeNumber = aws.String("Number")
	AwsTypeBinary = aws.String("Binary")
	AwsEmailValue = aws.String("email")
)

func NewSQSSender(client *sqs.Client, queueURL string) *SQSSender {
	return &SQSSender{
		Client:   client,
		QueueURL: queueURL,
	}
}

func (s *SQSSender) Send(ctx context.Context, payload []byte) error {
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.QueueURL),
		MessageBody: aws.String(string(payload))}

	if _, err := s.Client.SendMessage(ctx, input); err != nil {
		return fmt.Errorf("send SQS message: %w", err)
	}
	return nil
}

func (s *SQSSender) SendFifo(ctx context.Context, id string, groupID string, payload []byte) error {
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.QueueURL),
		MessageBody: aws.String(string(payload))}

	if strings.HasSuffix(s.QueueURL, ".fifo") {
		input.MessageGroupId = aws.String(groupID)
		input.MessageDeduplicationId = aws.String(id)
	}

	if _, err := s.Client.SendMessage(ctx, input); err != nil {
		return fmt.Errorf("send SQS message: %w", err)
	}
	return nil
}

func NewSQSEmailConsumer(queueURL string,
	threads int,
	emailHandler email.MailQueue,
	pool *ants.Pool,
) *SQSEmailConsumer {

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal().Msgf("Unable to load AWS config: %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	return &SQSEmailConsumer{
		QueueURL:     queueURL,
		Threads:      threads,
		Client:       client,
		EmailHandler: emailHandler,
		pool:         pool,
		backoff:      sys.NewDefaultBackoff(),
	}
}

// Consume starts consuming messages from the SQS queue.
// This is a long-running function that continuously polls the queue for messages.
func (c *SQSEmailConsumer) Consume(ctx context.Context) {

	log.Debug().Msgf("start sqs %s", c.QueueURL)

	for {
		// pull as many messages as we have worker threads
		messageBatch, err := c.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(c.QueueURL),
			MaxNumberOfMessages: int32(c.Threads),
			WaitTimeSeconds:     WaitTimeSecs,
		})

		if err != nil {
			log.Error().Msgf("failed to receive messages: %v", err)
			c.backoff.Sleep()
			continue
		}

		// if err = c.pool.Submit(func() {
		// 	err := processMessages(client, ctx, messageBatch)
		// 	if err != nil {
		// 		log.Error().Msgf("error processing messages: %v", err)
		// 	}
		// }); err != nil {
		// 	log.Error().Msgf("failed to submit processMessages: %v", err)
		// }

		// process each message individually using the worker pool
		for _, message := range messageBatch.Messages {
			submitErr := c.pool.Submit(func() {
				if err := c.processMessage(ctx, message); err != nil {
					log.Error().Msgf("error processing SQS message: %v", err)
				}
			})

			if submitErr != nil {
				log.Error().Msgf("failed to submit SQS message: %v", submitErr)
			}
		}
	}
}

// func processSqsMessages(client *sqs.Client, ctx context.Context, resp *sqs.ReceiveMessageOutput) error {

// 	var errs []error
// 	for _, message := range resp.Messages {
// 		if err := processSqsMessage(client, ctx, message); err != nil {
// 			errs = append(errs, err)
// 			continue
// 		}
// 	}

// 	if len(errs) > 0 {
// 		return errors.Join(errs...)
// 	}

// 	return nil
// }

func (c *SQSEmailConsumer) processMessage(ctx context.Context,
	message sqstypes.Message) error {

	messageID := "<unknown>"
	if message.MessageId != nil {
		messageID = *message.MessageId
	}

	if message.Body == nil || message.ReceiptHandle == nil {
		return fmt.Errorf("message %s has no body or receipt handle", messageID)

	}

	var m email.MailItem
	if err := json.Unmarshal([]byte(*message.Body), &m); err != nil {
		return fmt.Errorf("decode message %s: %w", messageID, err)
	}

	if err := c.EmailHandler.SendMail(&m); err != nil {
		return fmt.Errorf("send message %s: %w", messageID, err)
	}

	if _, err := c.Client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.QueueURL),
		ReceiptHandle: message.ReceiptHandle,
	}); err != nil {
		return fmt.Errorf("delete message %s: %w", messageID, err)
	}

	return nil
}

func NewSESMailer(from *mail.Address) email.Mailer {
	// Create a new SES session
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal().Msgf("unable to load SDK config, %v", err)
	}

	// Create a new SES client
	svc := sesv2.NewFromConfig(cfg)

	return &SESMailer{
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

func NewSQSEmailQueue(queueUrl string) *SQSEmailQueue {
	ctx := context.Background()

	// Load AWS config (region and credentials)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal().Msgf("unable to load SDK config: %v", err)
	}

	sqs := sqs.NewFromConfig(cfg)

	client := NewSQSSender(sqs, queueUrl)

	return &SQSEmailQueue{client: client}
}

// SendMail sends an email message to the SQS queue.
// The message is serialized as JSON and includes message attributes.
func (queue *SQSEmailQueue) SendMail(ctx context.Context, mail *email.MailItem) error {

	body, err := json.Marshal(mail)

	if err != nil {
		return err
	}

	err = queue.client.Send(ctx, body)

	// Send to SQS
	// _, err = queue.client.SendMessage(ctx, &sqs.SendMessageInput{
	// 	QueueUrl:    aws.String(queue.client.QueueUrl()),
	// 	MessageBody: aws.String(string(body)),
	// 	// Optional: add message attributes
	// 	MessageAttributes: map[string]sqstypes.MessageAttributeValue{
	// 		"type": {
	// 			DataType:    AwsTypeString,
	// 			StringValue: AwsEmailValue,
	// 		},
	// 	},
	// })

	return err
}
