package mail

import (
	"context"
	"net/mail"

	"github.com/antonybholmes/go-jobs"
)

type (
	Mailer interface {
		SendHtmlMail(ctx context.Context,
			to *mail.Address,
			subject string,
			html string) error
	}

	// QueueConsumer defines the interface for consuming messages from a queue.
	// Assumes a loop or continuous process for consuming messages from the queue.
	QueueConsumer interface {
		Consume(ctx context.Context)
	}

	Payload struct {
		DataType string `json:"dataType"`
		Data     string `json:"data"`
	}

	MailItem struct {
		Payload *Payload      `json:"payload,omitempty"`
		Job     *jobs.BaseJob `json:"job,omitempty"`
		// Name of the user, if applicable.
		UserID    string    `json:"userId,omitempty"`
		Name      string    `json:"name,omitempty"`
		To        string    `json:"to"`
		From      string    `json:"from"`
		EmailType EmailType `json:"type"`
		TTL       string    `json:"ttl,omitempty"`
		LinkUrl   string    `json:"linkUrl,omitempty"`
		Mode      string    `json:"mode,omitempty"`
	}

	// MailQueue defines the interface for sending mail items to a queue
	MailQueue interface {
		SendMail(ctx context.Context, mail *MailItem) error
	}

	EmailType string
)

const (
	EmailTypeVerify          EmailType = "verify"
	EmailTypeVerified        EmailType = "verified"
	EmailTypePasswordless    EmailType = "passwordless"
	EmailTypePasswordReset   EmailType = "password-reset"
	EmailTypePasswordUpdated EmailType = "password-updated"
	EmailTypeEmailReset      EmailType = "email-reset"
	EmailTypeEmailUpdated    EmailType = "email-updated"
	EmailTypeAccountCreated  EmailType = "account-created"
	EmailTypeAccountUpdated  EmailType = "account-updated"
	EmailTypeOTP             EmailType = "otp"
	EmailTypeTest            EmailType = "test"

	EmailTypeJobSubmitted  EmailType = "job-submitted"
	EmailTypeJobQueued     EmailType = "job-queued"
	EmailTypeJobStarted    EmailType = "job-started"
	EmailTypeJobError      EmailType = "job-error"
	EmailTypeJobCancelled  EmailType = "job-cancelled"
	EmailTypeJobInProgress EmailType = "job-in-progress"
	EmailTypeJobFailed     EmailType = "job-failed"
	EmailTypeJobComplete   EmailType = "job-complete"

	WaitTimeSecs = 10
)

func NewJobSubmittedMailItem(job *jobs.BaseJob, userID, to, from string) *MailItem {
	return &MailItem{
		UserID:    userID,
		Job:       job,
		To:        to,
		From:      from,
		EmailType: EmailTypeJobSubmitted,
	}
}
