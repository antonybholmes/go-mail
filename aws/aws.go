package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type (
	SqsSender struct {
		Client   *sqs.Client
		QueueURL string
	}
)

func NewSqsSender(client *sqs.Client, queueURL string) *SqsSender {
	return &SqsSender{
		Client:   client,
		QueueURL: queueURL,
	}
}

func (s *SqsSender) Send(ctx context.Context, payload []byte) error {
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.QueueURL),
		MessageBody: aws.String(string(payload))}

	if _, err := s.Client.SendMessage(ctx, input); err != nil {
		return fmt.Errorf("send SQS message: %w", err)
	}
	return nil
}

func (s *SqsSender) SendFifo(ctx context.Context, id string, groupID string, payload []byte) error {
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
