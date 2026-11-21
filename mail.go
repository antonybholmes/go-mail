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

	QueueProcessor interface {
		Consume(ctx context.Context)
	}
)
