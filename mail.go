package mail

import (
	"context"
	"net/mail"
)

type (
	Mailer interface {
		SendHtmlMail(ctx context.Context,
			to *mail.Address,
			subject string,
			html string) error
	}

	// QueueProcessor defines the interface for consuming messages from a queue.
	// Assumes a loop or continuous process for consuming messages from the queue.
	QueueProcessor interface {
		Consume(ctx context.Context)
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

	EmailJobSubmit     EmailType = "job-submit"
	EmailJobQueued     EmailType = "job-queued"
	EmailJobStarted    EmailType = "job-started"
	EmailJobError      EmailType = "job-error"
	EmailJobCancelled  EmailType = "job-cancelled"
	EmailJobInProgress EmailType = "job-in-progress"
	EmailJobFailed     EmailType = "job-failed"
	EmailJobComplete   EmailType = "job-complete"
)
