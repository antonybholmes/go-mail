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

	// MailQueue defines the interface for sending mail items to a queue
	MailQueue interface {
		SendMail(ctx context.Context, mail *jobs.MailItem) error
	}
)

const (
	EmailTypeVerify          jobs.EmailType = "verify"
	EmailTypeVerified        jobs.EmailType = "verified"
	EmailTypePasswordless    jobs.EmailType = "passwordless"
	EmailTypePasswordReset   jobs.EmailType = "password-reset"
	EmailTypePasswordUpdated jobs.EmailType = "password-updated"
	EmailTypeEmailReset      jobs.EmailType = "email-reset"
	EmailTypeEmailUpdated    jobs.EmailType = "email-updated"
	EmailTypeAccountCreated  jobs.EmailType = "account-created"
	EmailTypeAccountUpdated  jobs.EmailType = "account-updated"
	EmailTypeOTP             jobs.EmailType = "otp"
	EmailTypeTest            jobs.EmailType = "test"

	WaitTimeSecs = 10
)
