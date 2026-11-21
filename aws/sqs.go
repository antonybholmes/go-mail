package aws

import (
	"context"
	"encoding/json"
	"fmt"

	awsjobs "github.com/antonybholmes/go-jobs/aws"
	email "github.com/antonybholmes/go-mail"
	"github.com/antonybholmes/go-sys/log"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type (
	SQSEmailConsumer struct {
		Consumer     *awsjobs.SQSConsumer
		EmailHandler email.MailQueue
	}

	SQSEmailQueue struct {
		Sender *awsjobs.SQSSender
	}
)

func NewSQSEmailConsumer(consumer *awsjobs.SQSConsumer,
	emailHandler email.MailQueue,
) *SQSEmailConsumer {

	return &SQSEmailConsumer{
		Consumer:     consumer,
		EmailHandler: emailHandler,
	}
}

// Consume starts consuming messages from the SQS queue.
// This is a long-running function that continuously polls the queue for messages.
func (c *SQSEmailConsumer) Consume(ctx context.Context) {

	log.Debug().Msgf("start sqs email consumer for %s", c.Consumer.QueueURL)

	c.Consumer.Consume(ctx, func(ctx context.Context, message sqstypes.Message) error {
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

		if err := c.EmailHandler.SendMail(ctx, &m); err != nil {
			return fmt.Errorf("send message %s: %w", messageID, err)
		}

		if err := c.Consumer.DeleteMessage(ctx, message); err != nil {
			return fmt.Errorf("delete message %s: %w", messageID, err)
		}

		return nil
	})

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

func NewSQSEmailQueue(client *sqs.Client, queueUrl string) *SQSEmailQueue {
	// ctx := context.Background()

	// // Load AWS config (region and credentials)
	// cfg, err := config.LoadDefaultConfig(ctx)
	// if err != nil {
	// 	log.Fatal().Msgf("unable to load SDK config: %v", err)
	// }

	//sqs := sqs.NewFromConfig(cfg)

	sender := awsjobs.NewSQSSender(client, queueUrl)

	return &SQSEmailQueue{Sender: sender}
}

// SendMail sends an email message to the SQS queue.
// The message is serialized as JSON and includes message attributes.
func (queue *SQSEmailQueue) SendMail(ctx context.Context, mail *email.MailItem) error {

	body, err := json.Marshal(mail)

	if err != nil {
		return err
	}

	err = queue.Sender.Send(ctx, body)

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
